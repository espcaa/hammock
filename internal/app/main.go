package app

import (
	"encoding/json"
	"image"
	"log"
	"sync"

	"gioui.org/app"
	"gioui.org/font"
	"gioui.org/layout"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
	"gioui.org/unit"
	"github.com/espcaa/hammock/internal/slack"
	"github.com/espcaa/hammock/internal/store"
	"github.com/espcaa/hammock/internal/store/db"
	"github.com/espcaa/hammock/internal/ui"
)

type MainScreen struct {
	theme         *ui.Theme
	router        *Router
	store         *store.Store
	cache         *ui.ImageCache
	client        *slack.Client
	currentTeamId string
	sidebar       *ui.Sidebar
	msgList       *ui.MessageList
	selChannelID  string

	mu          sync.Mutex
	pendingMsgs []db.Message
}

func NewMainScreen(th *ui.Theme, r *Router, s *store.Store, client *slack.Client) *MainScreen {
	m := &MainScreen{
		theme:  th,
		router: r,
		store:  s,
		client: client,
		cache:  ui.NewImageCache(),
	}
	m.sidebar = th.Sidebar(s, "", m.cache)
	m.msgList = ui.NewMessageList(th, s, m.currentTeamId)
	m.sidebar.OnSelect(func(ch db.Channel) {
		m.loadMessages(ch.ID)
	})
	return m
}

func (m *MainScreen) ensureTeam() {
	if m.currentTeamId == "" && m.store.SlackSession != nil && len(m.store.SlackSession.WorkspacesIds) > 0 {
		m.currentTeamId = m.store.SlackSession.WorkspacesIds[0]
	}
	if m.currentTeamId != "" {
		m.sidebar.SetTeam(m.currentTeamId)
		m.msgList.SetTeam(m.currentTeamId)
	}
}

func (m *MainScreen) Layout(gtx layout.Context) layout.Dimensions {
	m.ensureTeam()
	m.applyPending()

	width := gtx.Dp(ui.SidebarWidth)

	return layout.Flex{Axis: layout.Horizontal}.Layout(gtx,
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			gtx.Constraints.Min.X = width
			gtx.Constraints.Max.X = width
			return m.sidebar.Layout(gtx)
		}),
		layout.Flexed(1, m.content),
	)
}

func (m *MainScreen) loadMessages(channelID string) {
	if m.currentTeamId == "" {
		return
	}
	teamID := m.currentTeamId

	m.mu.Lock()
	m.selChannelID = channelID
	m.mu.Unlock()

	msgs, err := m.store.ListMessages(teamID, channelID)
	if err != nil {
		log.Printf("list messages %s: %v", channelID, err)
		msgs = nil
	}
	m.msgList.Set(msgs)

	go func() {
		messages, err := m.client.GetConversationHistory(teamID, channelID, "", 28)
		if err != nil {
			log.Printf("conversation.history %s: %v", channelID, err)
			return
		}
		if err := m.store.UpsertMessages(messages, teamID, channelID); err != nil {
			log.Printf("save messages %s: %v", channelID, err)
			return
		}
		fresh, err := m.store.ListMessages(teamID, channelID)
		if err != nil {
			log.Printf("relist messages %s: %v", channelID, err)
			return
		}
		m.mu.Lock()
		current := m.selChannelID == channelID
		if current {
			m.pendingMsgs = fresh
		}
		m.mu.Unlock()
		if current {
			m.router.Invalidate()
		}
		log.Printf("loaded %d messages for %s", len(fresh), channelID)
	}()

	go func() {
		resp, err := m.client.GetUsersFromChannel(teamID, channelID, 30)
		if err != nil {
			log.Printf("get users %s: %v", channelID, err)
			return
		}
		log.Printf("loaded %d users for %s", len(resp), channelID)

		err = m.store.UpsertUsers(resp, teamID)
		if err != nil {
			log.Printf("save users %s: %v", channelID, err)
			return
		}
		if m.cache.Pending(channelID) {
			m.cache.Get(channelID)
			m.router.Invalidate()
		}
	}()
}

func (m *MainScreen) applyPending() {
	m.mu.Lock()
	msgs := m.pendingMsgs
	m.pendingMsgs = nil
	m.mu.Unlock()
	if msgs != nil {
		m.msgList.Set(msgs)
	}
}

func (m *MainScreen) content(gtx layout.Context) layout.Dimensions {
	ch, ok := m.sidebar.Selected()
	if !ok {
		return layout.Center.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
			lbl := m.theme.Label("Select a channel", unit.Sp(16), font.Medium, false)
			lbl.Color = m.theme.Faint
			return lbl.Layout(gtx)
		})
	}

	return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			return layout.Inset{
				Left: unit.Dp(16), Right: unit.Dp(16), Top: unit.Dp(14), Bottom: unit.Dp(14),
			}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
				return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
					layout.Rigid(func(gtx layout.Context) layout.Dimensions {
						return m.theme.Label("#"+ch.Name, unit.Sp(20), font.Bold, false).Layout(gtx)
					}),
					layout.Rigid(func(gtx layout.Context) layout.Dimensions {
						topic := topicText(ch.Topic)
						if topic == "" {
							return layout.Dimensions{}
						}
						lbl := m.theme.Label(topic, unit.Sp(13), font.Normal, false)
						lbl.Color = m.theme.Muted
						return layout.Inset{Top: unit.Dp(2)}.Layout(gtx, lbl.Layout)
					}),
				)
			})
		}),
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			w := gtx.Constraints.Max.X
			defer clip.Rect{Max: image.Pt(w, 1)}.Push(gtx.Ops).Pop()
			paint.Fill(gtx.Ops, m.theme.Border)
			return layout.Dimensions{Size: image.Pt(w, 1)}
		}),
		layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
			return m.msgList.Layout(gtx)
		}),
	)
}

func (m *MainScreen) WindowOptions() []app.Option {
	return []app.Option{
		app.MinSize(unit.Dp(400), unit.Dp(300)),
		app.MaxSize(unit.Dp(8000), unit.Dp(6000)),
		app.Maximized.Option(),
	}
}

func topicText(raw json.RawMessage) string {
	if len(raw) == 0 {
		return ""
	}
	var topic struct {
		Value string `json:"value"`
	}
	if err := json.Unmarshal(raw, &topic); err != nil {
		return ""
	}
	return topic.Value
}
