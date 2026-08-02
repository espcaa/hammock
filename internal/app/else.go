package app

import (
	"gioui.org/app"
	"gioui.org/layout"
	"gioui.org/unit"
	"github.com/espcaa/hammock/internal/ui"
)

type SomethingElse struct {
	theme   *ui.Theme
	backBtn ui.Button
	router  *Router
}

func NewSomethingElse(th *ui.Theme, r *Router) *SomethingElse {
	return &SomethingElse{
		theme:  th,
		router: r,
	}
}

func (s *SomethingElse) Layout(gtx layout.Context) layout.Dimensions {
	if s.backBtn.Click.Clicked(gtx) {
		s.router.Pop(gtx)
	}

	return layout.Center.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
		return layout.Flex{Axis: layout.Vertical, Alignment: layout.Middle}.Layout(gtx,
			layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				return s.theme.Button(&s.backBtn, "Back").Layout(gtx)
			}),
		)
	})
}

func (h *SomethingElse) WindowOptions() []app.Option {
	return []app.Option{
		app.MinSize(unit.Dp(300), unit.Dp(400)),
		app.MaxSize(unit.Dp(300), unit.Dp(400)),
		app.Size(unit.Dp(300), unit.Dp(400)),
	}
}
