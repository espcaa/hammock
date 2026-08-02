package main

import (
	"log"
	"os"

	gioapp "gioui.org/app"
	"gioui.org/op"

	"github.com/espcaa/hammock/internal/app"
	"github.com/espcaa/hammock/internal/misc"
)

func main() {
	go func() {
		w := new(gioapp.Window)
		w.Option(gioapp.Title("Hammock"), gioapp.Decorated(false))
		if err := loop(w); err != nil {
			log.Fatal(err)
		}
		os.Exit(0)
	}()
	gioapp.Main()
}

func loop(w *gioapp.Window) error {
	a := app.New()
	var styled bool
	var ops op.Ops
	for {
		switch e := w.Event().(type) {
		case gioapp.ViewEvent:
			if !styled {
				if handle := misc.NSViewHandle(e); handle != 0 {
					log.Println("hammock: got view handle, styling")
					w.Run(func() { misc.StyleTitlebar(handle) })
					styled = true
				} else {
					log.Println("hammock: viewevent but handle==0")
				}
			}

		case gioapp.DestroyEvent:
			return e.Err
		case gioapp.FrameEvent:
			gtx := gioapp.NewContext(&ops, e)

			a.Layout(gtx)
			e.Frame(gtx.Ops)
		}
	}
}
