package app

import (
	"gioui.org/app"
	"gioui.org/font"
	"gioui.org/layout"
	"gioui.org/unit"
	"github.com/espcaa/hammock/internal/ui"
)

type MainScreen struct {
	theme    *ui.Theme
	router   *Router
	btnState *ui.Button
}

func NewMainScreen(th *ui.Theme, r *Router) *MainScreen {
	return &MainScreen{
		theme:    th,
		router:   r,
		btnState: &ui.Button{},
	}
}

func (m *MainScreen) Layout(gtx layout.Context) layout.Dimensions {
	return layout.Center.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
		return layout.Flex{Axis: layout.Vertical, Alignment: layout.Middle}.Layout(gtx,
			layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				return m.theme.Label("Welcome to Hammock!", unit.Sp(24), font.Normal, false).Layout(gtx)
			}),
			layout.Rigid(layout.Spacer{Height: unit.Dp(16)}.Layout),
			layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				btn := m.theme.Button(m.btnState, "Go to Something Else")
				btn.TextSize = unit.Sp(20)
				btn.TextWeight = font.Medium
				return btn.Layout(gtx)
			}),
		)
	})
}

func (m *MainScreen) WindowOptions() []app.Option {
	return []app.Option{
		app.MinSize(unit.Dp(400), unit.Dp(300)),
		app.MaxSize(unit.Dp(8000), unit.Dp(6000)),
		app.Maximized.Option(),
	}
}
