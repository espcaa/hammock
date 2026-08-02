package app

import (
	"math"

	gioapp "gioui.org/app"
	"gioui.org/io/event"
	"gioui.org/io/key"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/op/paint"

	"github.com/espcaa/hammock/internal/config"
	"github.com/espcaa/hammock/internal/ui"
)

type App struct {
	theme  *ui.Theme
	router *Router
	scale  float32
	conf   *config.Config
}

func New(w *gioapp.Window) *App {
	conf, _ := config.LoadConfig()
	theme := ui.NewTheme()
	router := &Router{}
	router.SetWindow(w)
	a := &App{theme: theme, router: router, scale: conf.Scale, conf: conf}
	router.stack = []Screen{NewHomeScreen(theme, router)}
	return a
}

func (a *App) Layout(gtx layout.Context) layout.Dimensions {
	a.handleZoom(gtx)

	if a.scale > 0 {
		gtx.Metric.PxPerDp *= a.scale
		gtx.Metric.PxPerSp *= a.scale
	}

	paint.Fill(gtx.Ops, a.theme.Bg)
	return a.router.Layout(gtx)
}

// handle app level zooming

func (a *App) handleZoom(gtx layout.Context) {
	filters := []event.Filter{
		key.Filter{Name: "=", Required: key.ModShortcut},
		key.Filter{Name: "-", Required: key.ModShortcut},
		key.Filter{Name: "0", Required: key.ModShortcut},
	}

	for {
		ev, ok := gtx.Event(filters...)
		if !ok {
			break
		}
		ke, ok := ev.(key.Event)
		if !ok || ke.State != key.Press {
			continue
		}
		switch ke.Name {
		case "=":
			a.scale = snap(a.scale + 0.1)
		case "-":
			a.scale = snap(a.scale - 0.1)
		case "0":
			a.scale = 1.0
		}
		a.scale = clamp(a.scale, 0.5, 3.0)
		a.conf.Scale = a.scale
		a.conf.Save()
		gtx.Execute(op.InvalidateCmd{})
	}
}

func snap(v float32) float32 {
	return float32(math.Round(float64(v)*10) / 10)
}

func clamp(v, lo, hi float32) float32 {
	if v < lo {
		return lo
	}
	if v > hi {
		return hi
	}
	return v
}
