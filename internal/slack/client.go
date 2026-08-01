package slack

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"log"
	"mime/multipart"
	"net/url"
	"strconv"
	"strings"
	"time"

	http "github.com/bogdanfinn/fhttp"
	"github.com/coder/websocket"

	tls_client "github.com/bogdanfinn/tls-client"
	"github.com/bogdanfinn/tls-client/profiles"
)

const userAgent = "Mozilla/5.0 (Macintosh; Intel Mac OS X 15_6_0) AppleWebKit/537.36 (KHTML, like Gecko) Slack/4.48.102 Chrome/144.0.7559.236 Electron/40.8.2 Safari/537.36"

type Category struct {
	ID                   string `json:"channel_section_id"`
	Name                 string `json:"name"`
	Type                 string `json:"type"`
	Emoji                string `json:"emoji,omitempty"`
	NextChannelSectionID string `json:"next_channel_section_id,omitempty"`
	LastUpdated          int64  `json:"last_updated,omitempty"`
	IsRedacted           bool   `json:"is_redacted,omitempty"`
	ChannelIDsPage       struct {
		ChannelIDs []string `json:"channel_ids"`
		Count      int      `json:"count"`
		Cursor     string   `json:"cursor,omitempty"`
	} `json:"channel_ids_page"`
}

type Client struct {
	Session              *SlackSession
	HTTP                 tls_client.HttpClient
	WebsocketConnections map[string]*websocket.Conn
}

func NewClient(session *SlackSession) *Client {
	options := []tls_client.HttpClientOption{
		tls_client.WithClientProfile(profiles.Chrome_120),
		tls_client.WithTimeout(30000),
	}

	httpClient, err := tls_client.NewHttpClient(tls_client.NewNoopLogger(), options...)
	if err != nil {
		log.Fatalf("Failed to initialize tls-client: %v", err)
	}

	return &Client{
		Session: session,
		HTTP:    httpClient,
	}
}

func (c *Client) Do(teamID string, method string, params url.Values) (json.RawMessage, error) {
	return c.DoWithQuery(teamID, method, params, nil)
}

func (c *Client) DoWithQuery(teamID string, method string, params url.Values, query url.Values) (json.RawMessage, error) {
	ws, ok := c.Session.Workspaces[teamID]
	if !ok {
		return nil, errors.New("unknown workspace: " + teamID)
	}

	if params == nil {
		params = url.Values{}
	}
	params.Set("token", ws.Token)

	baseURL := strings.TrimRight(ws.TeamURL, "/")
	apiURL := baseURL + "/api/" + method
	if query != nil {
		apiURL += "?" + query.Encode()
	}

	var reqBody strings.Builder
	w := multipart.NewWriter(&reqBody)
	for key, vals := range params {
		for _, val := range vals {
			w.WriteField(key, val)
		}
	}
	w.Close()

	req, err := http.NewRequest("POST", apiURL, strings.NewReader(reqBody.String()))
	if err != nil {
		return nil, err
	}

	req.Header = http.Header{
		"content-type": {w.FormDataContentType()},
		"user-agent":   {userAgent},
		"origin":       {"https://app.slack.com"},
		"cookie":       {"d=" + c.Session.DCookie},
		"accept":       {"*/*"},
	}

	resp, err := c.HTTP.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode == 429 {
		retryAfter, _ := strconv.Atoi(resp.Header.Get("Retry-After"))
		if retryAfter == 0 {
			retryAfter = 5
		}
		time.Sleep(time.Duration(retryAfter) * time.Second)
		return c.Do(teamID, method, params)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var envelope struct {
		OK    bool   `json:"ok"`
		Error string `json:"error"`
	}
	if err := json.Unmarshal(body, &envelope); err != nil {
		return nil, err
	}

	if !envelope.OK {
		if envelope.Error == "ratelimited" {
			time.Sleep(5 * time.Second)
			return c.Do(teamID, method, params)
		}
		return nil, errors.New("slack api error: " + envelope.Error)
	}

	return body, nil
}

func (c *Client) DoEdge(teamID string, resource string, payload map[string]any) (json.RawMessage, error) {
	ws, ok := c.Session.Workspaces[teamID]
	if !ok {
		return nil, errors.New("unknown workspace: " + teamID)
	}

	team := teamID
	if ws.EnterpriseID != "" {
		team = ws.EnterpriseID
	}

	payload["token"] = ws.Token
	payload["enterprise_token"] = ws.Token
	payload["_x_app_name"] = "client"
	payload["fp"] = "60"

	jsonData, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}

	apiURL := "https://edgeapi.slack.com/cache/" + team + "/" + resource + "?_x_app_name=client&fp=60&_x_num_retries=0"

	req, err := http.NewRequest("POST", apiURL, bytes.NewReader(jsonData))
	if err != nil {
		return nil, err
	}

	req.Header = http.Header{
		"content-type": {"application/json"},
		"cookie":       {"d=" + c.Session.DCookie},
		"user-agent":   {userAgent},
		"origin":       {"https://app.slack.com"},
		"accept":       {"*/*"},
	}
	req.Header[http.HeaderOrderKey] = []string{
		"content-type",
		"cookie",
		"user-agent",
		"origin",
		"accept",
	}

	resp, err := c.HTTP.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode == 429 {
		retryAfter, _ := strconv.Atoi(resp.Header.Get("Retry-After"))
		if retryAfter == 0 {
			retryAfter = 5
		}
		time.Sleep(time.Duration(retryAfter) * time.Second)
		return c.DoEdge(teamID, resource, payload)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var envelope struct {
		OK    bool   `json:"ok"`
		Error string `json:"error"`
	}
	if err := json.Unmarshal(body, &envelope); err != nil {
		return nil, err
	}

	if !envelope.OK {
		if envelope.Error == "ratelimited" {
			time.Sleep(5 * time.Second)
			return c.DoEdge(teamID, resource, payload)
		}
		return nil, errors.New("slack edge api error: " + envelope.Error)
	}

	return body, nil
}
