package slack

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/coder/websocket"
)

type RTMEvent struct {
	Type     string          `json:"type"`
	SubType  string          `json:"subtype,omitempty"`
	Channel  string          `json:"channel,omitempty"`
	User     json.RawMessage `json:"user,omitempty"`
	Text     string          `json:"text,omitempty"`
	Ts       string          `json:"ts,omitempty"`
	ThreadTs string          `json:"thread_ts,omitempty"`
	Team     string          `json:"team,omitempty"`
	Raw      json.RawMessage `json:"-"`
}

func (e RTMEvent) UserID() string {
	if len(e.User) == 0 {
		return ""
	}

	var s string
	if json.Unmarshal(e.User, &s) == nil {
		return s
	}

	var obj struct {
		ID string `json:"id"`
	}
	if json.Unmarshal(e.User, &obj) == nil {
		return obj.ID
	}
	return ""
}

type EventHandler func(teamID string, event RTMEvent)

type RTMConnection struct {
	teamID  string
	client  *Client
	conn    *websocket.Conn
	ctx     context.Context
	cancel  context.CancelFunc
	handler EventHandler
}

func (c *Client) ConnectRTM(teamID string, handler EventHandler) (*RTMConnection, error) {
	ws := c.Session.Workspaces[teamID]
	wsURL := ConstructWebsocketURL(ws.Token, teamID, ws.EnterpriseID)

	dialCtx, dialCancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer dialCancel()

	conn, _, err := websocket.Dial(dialCtx, wsURL, &websocket.DialOptions{
		HTTPHeader: http.Header{
			"User-Agent": {userAgent},
			"Cookie":     {"d=" + c.Session.DCookie},
		},
	})
	if err != nil {
		return nil, fmt.Errorf("RTM dial failed for %s: %w", teamID, err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	rtm := &RTMConnection{
		teamID:  teamID,
		client:  c,
		conn:    conn,
		ctx:     ctx,
		cancel:  cancel,
		handler: handler,
	}

	go rtm.readLoop()
	go rtm.pingLoop()
	return rtm, nil
}

func (rtm *RTMConnection) pingLoop() {
	ticker := time.NewTicker(15 * time.Second)
	defer ticker.Stop()

	id := 0
	for {
		select {
		case <-rtm.ctx.Done():
			return
		case <-ticker.C:
			id++
			msg, _ := json.Marshal(map[string]any{"type": "ping", "id": id})
			if err := rtm.conn.Write(rtm.ctx, websocket.MessageText, msg); err != nil {
				log.Printf("RTM ping failed for %s: %v", rtm.teamID, err)
				return
			}
		}
	}
}

func (rtm *RTMConnection) readLoop() {
	for {
		_, data, err := rtm.conn.Read(rtm.ctx)
		if err != nil {
			log.Printf("RTM read error for %s: %v", rtm.teamID, err)
			rtm.reconnect()
			return
		}

		var event RTMEvent
		if err := json.Unmarshal(data, &event); err != nil {
			log.Printf("RTM parse error: %v", err)
			continue
		}
		event.Raw = data

		if event.Type == "error" {
			log.Printf("RTM [%s] error: %s", rtm.teamID, string(data))
		}

		rtm.handler(rtm.teamID, event)
	}
}

func (rtm *RTMConnection) reconnect() {
	rtm.conn.Close(websocket.StatusNormalClosure, "reconnecting")

	for {
		log.Printf("RTM reconnecting for %s...", rtm.teamID)
		ws := rtm.client.Session.Workspaces[rtm.teamID]
		wsURL := ConstructWebsocketURL(ws.Token, rtm.teamID, ws.EnterpriseID)

		dialCtx, dialCancel := context.WithTimeout(rtm.ctx, 30*time.Second)
		conn, _, err := websocket.Dial(dialCtx, wsURL, &websocket.DialOptions{
			HTTPHeader: http.Header{
				"User-Agent": {userAgent},
				"Cookie":     {"d=" + rtm.client.Session.DCookie},
			},
		})
		dialCancel()

		if err != nil {
			log.Printf("RTM reconnect failed for %s: %v", rtm.teamID, err)
			time.Sleep(5 * time.Second)
			continue
		}

		rtm.conn = conn
		go rtm.readLoop()
		return
	}
}

func (rtm *RTMConnection) SendTyping(channelID string, threadTS string, id int) error {
	payload := map[string]any{"type": "user_typing", "channel": channelID, "id": id}
	if threadTS != "" {
		payload["thread_ts"] = threadTS
	}
	msg, _ := json.Marshal(payload)
	return rtm.conn.Write(rtm.ctx, websocket.MessageText, msg)
}

func (rtm *RTMConnection) Close() {
	rtm.cancel()
	rtm.conn.Close(websocket.StatusNormalClosure, "closing")
}

func ConstructWebsocketURL(token, teamID, enterpriseID string) string {
	params := url.Values{}
	params.Set("token", token)
	params.Set("sync_desync", "1")
	params.Set("slack_client", "desktop")
	params.Set("batch_presence_aware", "1")
	params.Set("no_query_on_subscribe", "1")
	params.Set("flannel", "3")
	params.Set("lazy_channels", "1")
	params.Set("start_args", "?agent=client&org_wide_aware=true&agent_version=1775954273&eac_cache_ts=true&cache_ts=0&name_tagging=true&only_self_subteams=true&connect_only=true&ms_latest=true")

	if enterpriseID != "" {
		// gateway_server uses the enterprise ID with E→T prefix swap
		params.Set("gateway_server", "T"+strings.TrimPrefix(enterpriseID, "E")+"-1")
		params.Set("enterprise_id", enterpriseID)
	} else {
		params.Set("gateway_server", teamID+"-1")
	}

	return "wss://wss-primary.slack.com/?" + params.Encode()
}
