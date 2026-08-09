package app

import (
	"log"

	"gioui.org/app"
	gioapp "gioui.org/app"
	"gioui.org/layout"
	"gioui.org/op"
)

type Router struct {
	stack     []Screen
	lastSized Screen
	win       *gioapp.Window
}

type WindowSizer interface {
	WindowOptions() []app.Option
}

func NewRouter(initial Screen) *Router {
	return &Router{stack: []Screen{initial}}
}

func (r *Router) SetWindow(w *gioapp.Window) {
	r.win = w
}

func (r *Router) Push(gtx layout.Context, s Screen) {
	r.stack = append(r.stack, s)
	gtx.Execute(op.InvalidateCmd{})
}

func (r *Router) Pop(gtx layout.Context) {
	if len(r.stack) > 1 {
		r.stack = r.stack[:len(r.stack)-1]
		gtx.Execute(op.InvalidateCmd{})
	}
}

func (r *Router) Replace(gtx layout.Context, s Screen) {
	gtx.Execute(op.InvalidateCmd{})
	if len(r.stack) == 0 {
		r.stack = []Screen{s}
		return
	}
	r.stack[len(r.stack)-1] = s
}

func (r *Router) Current() Screen {
	if len(r.stack) == 0 {
		return nil
	}
	return r.stack[len(r.stack)-1]
}

func (r *Router) Layout(gtx layout.Context) layout.Dimensions {
	cur := r.Current()
	if cur == nil {
		return layout.Dimensions{}
	}
	if r.lastSized != cur {
		r.applyOptionsFor(cur)
		r.lastSized = cur
	}
	return cur.Layout(gtx)
}

func (r *Router) applyOptionsFor(s Screen) {
	log.Printf("Applying window options for screen: %T", s)
	if sizer, ok := s.(WindowSizer); ok {
		log.Printf("Screen %T implements WindowSizer, applying options", s)
		for _, opt := range sizer.WindowOptions() {
			log.Printf("Applying option: %v", opt)
		}
		r.win.Option(sizer.WindowOptions()...)
	}
	r.lastSized = s
}

func (r *Router) Invalidate() {
	if r.win != nil {
		r.win.Invalidate()
	}
}
