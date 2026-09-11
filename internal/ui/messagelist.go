package ui

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"gioui.org/font"
	"gioui.org/layout"
	"gioui.org/unit"

	"github.com/espcaa/hammock/internal/store/db"
)

type MessageList struct {
	theme    *Theme
	list     layout.List
	messages []db.Message
}

func NewMessageList(th *Theme) *MessageList {
	return &MessageList{
		theme: th,
		list:  layout.List{Axis: layout.Vertical},
	}
}

func (m *MessageList) Set(messages []db.Message) {
	if len(messages) == 0 {
		m.messages = nil
		return
	}
	m.messages = make([]db.Message, len(messages))
	for i, msg := range messages {
		m.messages[len(messages)-1-i] = msg
	}
	m.list.ScrollToEnd = true
}

func (m *MessageList) Layout(gtx layout.Context) layout.Dimensions {
	if len(m.messages) == 0 {
		return layout.Center.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
			lbl := m.theme.Label("No messages yet", unit.Sp(14), font.Normal, false)
			lbl.Color = m.theme.Faint
			return lbl.Layout(gtx)
		})
	}
	return m.list.Layout(gtx, len(m.messages), m.msgRow)
}

func (m *MessageList) msgRow(gtx layout.Context, i int) layout.Dimensions {
	msg := m.messages[i]
	text := messageText(msg)
	ts := formatMessageTS(msg.Ts)

	return layout.Inset{
		Left: unit.Dp(16), Right: unit.Dp(16), Top: unit.Dp(6), Bottom: unit.Dp(6),
	}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
		return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
			layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				return layout.Flex{Axis: layout.Horizontal, Alignment: layout.Middle}.Layout(gtx,
					layout.Rigid(func(gtx layout.Context) layout.Dimensions {
						lbl := m.theme.Label(m.senderName(msg), unit.Sp(13), font.Bold, false)
						lbl.Color = m.theme.Text
						return lbl.Layout(gtx)
					}),
					layout.Rigid(layout.Spacer{Width: unit.Dp(8)}.Layout),
					layout.Rigid(func(gtx layout.Context) layout.Dimensions {
						lbl := m.theme.Label(ts, unit.Sp(11), font.Normal, false)
						lbl.Color = m.theme.Faint
						return lbl.Layout(gtx)
					}),
				)
			}),
			layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				return layout.Inset{Top: unit.Dp(2)}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
					lbl := m.theme.Label(text, unit.Sp(14), font.Normal, false)
					lbl.Color = m.theme.Text
					return lbl.Layout(gtx)
				})
			}),
		)
	})
}

func (m *MessageList) senderName(msg db.Message) string {
	if msg.User.Valid {
		return msg.User.String
	}
	return "System"
}

func messageText(msg db.Message) string {
	if msg.Text.Valid && strings.TrimSpace(msg.Text.String) != "" {
		return msg.Text.String
	}
	return "…"
}

func formatMessageTS(ts string) string {
	sec, nsec, err := splitTS(ts)
	if err != nil {
		return ""
	}
	return time.Unix(sec, nsec).Local().Format("3:04 PM")
}

func splitTS(ts string) (int64, int64, error) {
	parts := strings.SplitN(ts, ".", 2)
	if len(parts) == 0 || parts[0] == "" {
		return 0, 0, fmt.Errorf("bad ts %q", ts)
	}
	sec, err := strconv.ParseInt(parts[0], 10, 64)
	if err != nil {
		return 0, 0, err
	}
	if len(parts) == 2 {
		micro, err := strconv.ParseInt(parts[1], 10, 64)
		if err != nil {
			return 0, 0, err
		}
		return sec, micro * 1000, nil
	}
	return sec, 0, nil
}
