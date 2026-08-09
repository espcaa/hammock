package ui

import (
	"image"
	"math"
	"time"

	"gioui.org/f32"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
	"gioui.org/unit"
)

type Spinner struct {
	th     *Theme
	Period time.Duration
	Size   unit.Dp
}

func (t *Theme) Spinner(size unit.Dp) Spinner {
	return Spinner{
		th:   t,
		Size: unit.Dp(size),
	}
}

func (s Spinner) Layout(gtx layout.Context) layout.Dimensions {
	if s.Period == 0 {
		s.Period = 1200 * time.Millisecond // ai said it was a good default, looks good
	}

	size := gtx.Dp(s.Size)
	if size == 0 {
		size = gtx.Dp(24)
	}

	t := gtx.Now.UnixNano() % int64(s.Period)
	frac := float32(t) / float32(s.Period)
	angle := frac * 2 * math.Pi

	center := f32.Pt(float32(size)/2, float32(size)/2)
	defer op.Affine(f32.Affine2D{}.Rotate(center, angle)).Push(gtx.Ops).Pop()

	radius := float32(size) / 2
	var p clip.Path
	p.Begin(gtx.Ops)
	p.MoveTo(f32.Pt(center.X+radius*0.8, center.Y))
	p.Arc(
		f32.Pt(-radius*0.8, 0),
		f32.Pt(-radius*0.8, 0),
		1.5*math.Pi,
	)
	paint.FillShape(gtx.Ops, s.th.Text, clip.Stroke{
		Path:  p.End(),
		Width: float32(size) / 10,
	}.Op())

	gtx.Execute(op.InvalidateCmd{})

	return layout.Dimensions{Size: image.Pt(size, size)}
}
