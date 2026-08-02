package main

import (
	"log"
	"os"

	gioapp "gioui.org/app"
	"gioui.org/op"

	"github.com/espcaa/hammock/internal/app"
)

func main() {
	go func() {
		w := new(gioapp.Window)
		w.Option(gioapp.Title("Hammock"), gioapp.Decorated(true)) // false removes the entire thing naur, TODO: make cleaner macos title bar
		if err := loop(w); err != nil {
			log.Fatal(err)
		}
		os.Exit(0)
	}()
	gioapp.Main()
}

func loop(w *gioapp.Window) error {
	a := app.New(w)
	var ops op.Ops
	for {
		switch e := w.Event().(type) {
		case gioapp.DestroyEvent:
			return e.Err
		case gioapp.ConfigEvent:
			w.Invalidate()
			log.Printf("ConfigEvent: %v", e)
		case gioapp.FrameEvent:
			gtx := gioapp.NewContext(&ops, e)

			a.Layout(gtx)
			e.Frame(gtx.Ops)
		}
	}
}
