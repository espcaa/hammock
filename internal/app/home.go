package app

import (
	"image"
	"image/color"

	"gioui.org/layout"
	"gioui.org/op/paint"
	"gioui.org/unit"
	"github.com/espcaa/hammock/internal/misc"
	"github.com/espcaa/hammock/internal/ui"
)

type HomeScreen struct {
	theme      *ui.Theme
	router     *Router
	sendBtn    ui.Button
	testAvatar image.Image
}

func NewHomeScreen(th *ui.Theme, r *Router) *HomeScreen {
	return &HomeScreen{
		theme:  th,
		router: r,
		testAvatar: misc.GradientImage(100, 100,
			color.NRGBA{R: 0x83, G: 0xA5, B: 0x98, A: 0xFF},
			color.NRGBA{R: 0xD5, G: 0xC4, B: 0xA1, A: 0xFF}),
	}
}

func (a *HomeScreen) Layout(gtx layout.Context) layout.Dimensions {
	paint.Fill(gtx.Ops, a.theme.Bg)

	if a.sendBtn.Click.Clicked(gtx) {
		a.router.Push(gtx, NewSomethingElse(a.theme, a.router))
	}

	return layout.Center.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
		return layout.Flex{Axis: layout.Vertical, Alignment: layout.Middle}.Layout(gtx,
			layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				return a.theme.Avatar(a.testAvatar, unit.Dp(40)).Layout(gtx)
			}),
			layout.Rigid(layout.Spacer{Height: unit.Dp(12)}.Layout),
			layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				return a.theme.Button(&a.sendBtn, "Send").Layout(gtx)
			}),
		)
	})
}
