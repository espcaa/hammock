package app

import (
	"gioui.org/layout"
	"gioui.org/op/paint"

	"github.com/espcaa/hammock/internal/ui"
)

type App struct {
	theme  *ui.Theme
	router *Router
}

func New() *App {
	theme := ui.NewTheme()
	router := &Router{}
	a := &App{theme: theme, router: router}
	router.stack = []Screen{NewHomeScreen(theme, router)}
	return a
}

func (a *App) Layout(gtx layout.Context) layout.Dimensions {
	paint.Fill(gtx.Ops, a.theme.Bg)
	return a.router.Layout(gtx)
}
