package app

import (
	"bytes"
	_ "embed"
	"image"
	_ "image/jpeg"
	"log"

	"gioui.org/app"
	"gioui.org/font"
	"gioui.org/layout"
	"gioui.org/op/paint"
	"gioui.org/unit"
	"github.com/espcaa/hammock/internal/ui"
)

//go:embed assets/onboarding.jpg
var onboarding []byte

type HomeScreen struct {
	theme       *ui.Theme
	router      *Router
	loginButton ui.Button
	onboardImg  ui.Image
}

func NewHomeScreen(th *ui.Theme, r *Router) *HomeScreen {
	img, _, err := image.Decode(bytes.NewReader(onboarding))
	if err != nil {
		log.Fatalf("decode onboarding.jpg: %v", err)
	}

	oi := th.Image(img, 300, 0)
	oi.Fill = true

	return &HomeScreen{
		theme:      th,
		router:     r,
		onboardImg: oi,
	}
}

func (a *HomeScreen) Layout(gtx layout.Context) layout.Dimensions {
	paint.Fill(gtx.Ops, a.theme.Bg)

	if a.loginButton.Click.Clicked(gtx) {
		a.router.Push(gtx, NewSomethingElse(a.theme, a.router))
	}

	return layout.Flex{Axis: layout.Vertical, Alignment: layout.Middle}.Layout(gtx,
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			return a.onboardImg.Layout(gtx)
		}),
		layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
			return layout.Center.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
				return layout.Flex{Axis: layout.Vertical, Alignment: layout.Middle}.Layout(gtx,
					layout.Rigid(func(gtx layout.Context) layout.Dimensions {
						return a.theme.Label("Welcome to Hammock!", unit.Sp(30), font.Bold, false).Layout(gtx)
					}),
					layout.Rigid(layout.Spacer{Height: unit.Dp(8)}.Layout),
					layout.Rigid(layout.Spacer{Height: unit.Dp(8)}.Layout),
					layout.Rigid(func(gtx layout.Context) layout.Dimensions {
						btn := a.theme.Button(&a.loginButton, "Login")
						btn.TextSize = unit.Sp(20)
						btn.TextWeight = font.Medium
						return btn.Layout(gtx)
					}),
				)
			})
		}),
	)
}

func (h *HomeScreen) WindowOptions() []app.Option {
	return []app.Option{
		app.Size(unit.Dp(300), unit.Dp(400)),
		app.MaxSize(unit.Dp(300), unit.Dp(400)),
		app.MinSize(unit.Dp(300), unit.Dp(400)),
	}
}
