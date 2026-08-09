package app

import (
	"gioui.org/app"
	gioapp "gioui.org/app"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/unit"
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
	if r.win == nil {
		r.lastSized = s
		return
	}
	sizer, ok := s.(WindowSizer)
	if !ok {
		r.lastSized = s
		return
	}
	opts := sizer.WindowOptions()

	// apply max/min size first & then the mode
	var sizeOpts, modeOpts []app.Option
	for _, o := range opts {
		if isModeOption(o) {
			modeOpts = append(modeOpts, o)
		} else {
			sizeOpts = append(sizeOpts, o)
		}
	}
	if len(sizeOpts) > 0 {
		r.win.Option(append(sizeOpts, app.Windowed.Option())...)
	}
	if len(modeOpts) > 0 {
		r.win.Option(modeOpts...)
	}
	r.lastSized = s
}

func isModeOption(o app.Option) bool {
	cnf := app.Config{Mode: app.Windowed}
	o(unit.Metric{PxPerDp: 1, PxPerSp: 1}, &cnf)
	return cnf.Mode != app.Windowed
}

func (r *Router) Invalidate() {
	if r.win != nil {
		r.win.Invalidate()
	}
}
