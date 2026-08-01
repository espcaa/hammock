package app

import (
	"gioui.org/layout"
	"gioui.org/op"
)

type Router struct {
	stack      []Screen
	invalidate bool
}

func NewRouter(initial Screen) *Router {
	return &Router{stack: []Screen{initial}}
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
	return r.Current().Layout(gtx)
}
