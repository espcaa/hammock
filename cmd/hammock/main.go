package main

import (
	"fmt"
	"log"
	"os"
	"runtime"

	gioapp "gioui.org/app"
	"gioui.org/op"

	"github.com/espcaa/hammock/internal/app"
	"github.com/espcaa/hammock/internal/login"
	"github.com/espcaa/hammock/internal/slack"
)

func main() {

	if len(os.Args) > 1 && os.Args[1] == "__login" {
		runtime.LockOSThread() // must own main thread
		loginURL := "https://app.slack.com/ssb/signin?" + slack.GenerateSSBParams()

		if len(os.Args) > 2 {
			loginURL = os.Args[2]
		}
		fmt.Print(login.RunLoginWebview(loginURL))
		return
	}

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
