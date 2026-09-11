package ui

import (
	"image"
	"image/color"
	"log"
	"strconv"

	"gioui.org/font"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
	"gioui.org/unit"
	"gioui.org/widget"

	"github.com/espcaa/hammock/internal/misc"
	"github.com/espcaa/hammock/internal/store"
	"github.com/espcaa/hammock/internal/store/db"
)

const SidebarWidth = unit.Dp(300)

type Sidebar struct {
	theme  *Theme
	store  *store.Store
	teamID string

	list     layout.List
	clicks   []widget.Clickable
	channels []db.Channel
	loaded   bool
	cache    *ImageCache

	selectedID string
	onSelect   func(db.Channel)
}

func (t *Theme) Sidebar(st *store.Store, teamID string, cache *ImageCache) *Sidebar {
	return &Sidebar{
		theme:  t,
		store:  st,
		teamID: teamID,
		list:   layout.List{Axis: layout.Vertical},
		cache:  cache,
	}
}

func (s *Sidebar) SetTeam(teamID string) {
	if s.teamID == teamID {
		return
	}
	s.teamID = teamID
	s.loaded = false
	s.selectedID = ""
}

func (s *Sidebar) Reload() { s.loaded = false }

func (s *Sidebar) OnSelect(fn func(db.Channel)) {
	if fn == nil {
		fn = func(db.Channel) {}
	}
	s.onSelect = fn
}

func (s *Sidebar) Selected() (db.Channel, bool) {
	for _, ch := range s.channels {
		if ch.ID == s.selectedID {
			return ch, true
		}
	}
	return db.Channel{}, false
}

func (s *Sidebar) load() {
	s.loaded = true
	s.channels = nil
	s.clicks = nil
	if s.teamID == "" || s.store == nil {
		return
	}
	channels, err := s.store.ListChannels(s.teamID, []string{"channel-public"})
	if err != nil {
		log.Printf("sidebar: list channels for %s: %v", s.teamID, err)
		return
	}
	s.channels = channels
	s.clicks = make([]widget.Clickable, len(channels))
}

func (s *Sidebar) Layout(gtx layout.Context) layout.Dimensions {
	if !s.loaded {
		s.load()
	}

	return layout.Background{}.Layout(gtx,
		func(gtx layout.Context) layout.Dimensions {
			w := gtx.Constraints.Min.X
			h := gtx.Constraints.Min.Y
			defer clip.Rect{Max: image.Pt(w, h)}.Push(gtx.Ops).Pop()
			paint.Fill(gtx.Ops, s.theme.Surface)

			defer clip.Rect{Min: image.Pt(w-1, 0), Max: image.Pt(w, h)}.Push(gtx.Ops).Pop()
			paint.Fill(gtx.Ops, s.theme.Border)

			return layout.Dimensions{Size: image.Pt(w, h)}
		},
		func(gtx layout.Context) layout.Dimensions {
			return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
				layout.Rigid(s.workspaceHeader),
				layout.Rigid(func(gtx layout.Context) layout.Dimensions {
					lbl := s.theme.Label("Channels", unit.Sp(11), font.Bold, false)
					lbl.Color = s.theme.Faint
					return layout.Inset{Left: unit.Dp(10), Top: unit.Dp(4), Bottom: unit.Dp(4)}.Layout(gtx, lbl.Layout)
				}),
				layout.Flexed(1, s.channelList),
			)
		},
	)
}

func (s *Sidebar) workspaceHeader(gtx layout.Context) layout.Dimensions {
	return layout.Inset{
		Top: unit.Dp(12), Bottom: unit.Dp(8), Left: unit.Dp(10), Right: unit.Dp(10),
	}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
		return layout.Flex{Axis: layout.Horizontal, Alignment: layout.Middle}.Layout(gtx,
			layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				return s.theme.RemoteImage(s.cache, s.store.SlackSession.WorkspaceSessions[s.teamID].TeamIcon, unit.Dp(48)).Layout(gtx)
			}),
			layout.Rigid(layout.Spacer{Width: unit.Dp(10)}.Layout),
			layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				lbl := s.theme.Label(shortID(s.store.SlackSession.WorkspaceSessions[s.teamID].TeamName), unit.Sp(16), font.Bold, true)
				lbl.Color = s.theme.Faint
				return lbl.Layout(gtx)
			}),
		)
	})
}

func shortID(id string) string {
	if len(id) <= 11 || id == "" {
		return id
	}
	return id[:11]
}

func (s *Sidebar) channelList(gtx layout.Context) layout.Dimensions {
	if len(s.channels) == 0 {
		lbl := s.theme.Label("No channels", unit.Sp(13), font.Normal, false)
		lbl.Color = s.theme.Faint
		return layout.Inset{Left: unit.Dp(10)}.Layout(gtx, lbl.Layout)
	}
	return s.list.Layout(gtx, len(s.channels), s.channelRow)
}

func (s *Sidebar) channelRow(gtx layout.Context, i int) layout.Dimensions {
	ch := s.channels[i]
	c := &s.clicks[i]
	selected := ch.ID == s.selectedID

	if c.Clicked(gtx) && !selected {
		s.selectedID = ch.ID
		if s.onSelect != nil {
			s.onSelect(ch)
		}
		gtx.Execute(op.InvalidateCmd{})
	}

	return c.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
		return layout.Background{}.Layout(gtx,
			func(gtx layout.Context) layout.Dimensions {
				var bg color.NRGBA
				switch {
				case selected:
					bg = s.theme.Overlay
				case c.Hovered():
					bg = misc.Lighten(s.theme.Surface, 0.06)
				}
				if bg != (color.NRGBA{}) {
					rr := gtx.Dp(s.theme.Radius)
					defer clip.RRect{
						Rect: image.Rectangle{Max: gtx.Constraints.Min},
						NW:   rr, NE: rr, SW: rr, SE: rr,
					}.Push(gtx.Ops).Pop()
					paint.Fill(gtx.Ops, bg)
				}
				return layout.Dimensions{Size: gtx.Constraints.Min}
			},
			func(gtx layout.Context) layout.Dimensions {
				return layout.Inset{
					Top: unit.Dp(4), Bottom: unit.Dp(4), Left: unit.Dp(6), Right: unit.Dp(6),
				}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
					return layout.Flex{Axis: layout.Horizontal, Alignment: layout.Middle}.Layout(gtx,
						layout.Rigid(func(gtx layout.Context) layout.Dimensions {
							lbl := s.theme.Label("#", unit.Sp(14), font.Bold, false)
							lbl.Color = s.theme.Faint
							return lbl.Layout(gtx)
						}),
						layout.Rigid(layout.Spacer{Width: unit.Dp(6)}.Layout),
						layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
							lbl := s.theme.Label(ch.Name, unit.Sp(14), font.Medium, false)
							if selected {
								lbl.Color = s.theme.Text
							} else {
								lbl.Color = s.theme.Muted
							}
							return lbl.Layout(gtx)
						}),
						layout.Rigid(func(gtx layout.Context) layout.Dimensions {
							return s.unreadBadge(gtx, ch.Unread)
						}),
					)
				})
			},
		)
	})
}

func (s *Sidebar) unreadBadge(gtx layout.Context, n int64) layout.Dimensions {
	if n <= 0 {
		return layout.Dimensions{}
	}
	return layout.Background{}.Layout(gtx,
		func(gtx layout.Context) layout.Dimensions {
			rr := gtx.Dp(unit.Dp(9))
			rect := image.Rectangle{Max: gtx.Constraints.Min}
			defer clip.RRect{Rect: rect, NW: rr, NE: rr, SW: rr, SE: rr}.Push(gtx.Ops).Pop()
			paint.Fill(gtx.Ops, s.theme.Primary)
			return layout.Dimensions{Size: gtx.Constraints.Min}
		},
		func(gtx layout.Context) layout.Dimensions {
			return layout.Inset{
				Top: unit.Dp(1), Bottom: unit.Dp(1), Left: unit.Dp(7), Right: unit.Dp(7),
			}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
				lbl := s.theme.Label(strconv.FormatInt(n, 10), unit.Sp(11), font.Bold, false)
				lbl.Color = s.theme.TextOnPrimary
				return lbl.Layout(gtx)
			})
		},
	)
}
