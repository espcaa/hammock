package app

import (
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
	r.applyOptionsFor(cur)
	return cur.Layout(gtx)
}

func (r *Router) applyOptionsFor(s Screen) {
	if r.win == nil || s == nil || s == r.lastSized {
		return
	}
	if sizer, ok := s.(WindowSizer); ok {
		r.win.Option(sizer.WindowOptions()...)
	}
	r.lastSized = s
}
