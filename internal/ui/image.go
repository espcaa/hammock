package ui

import (
	"image"
	"image/color"

	"gioui.org/layout"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
	"gioui.org/widget"
)

type Image struct {
	Op            paint.ImageOp
	Height, Width int
	Fill          bool
	valid         bool
}

func (t *Theme) Image(img image.Image, height, width int) Image {
	i := Image{Height: height, Width: width}
	if img != nil {
		i.Op = paint.NewImageOp(img)
		i.valid = true
	}
	return i
}

func (i Image) Layout(gtx layout.Context) layout.Dimensions {
	sz := image.Pt(i.Width, i.Height)

	if i.Fill {
		w := gtx.Constraints.Max.X
		src := i.Op.Size()
		if src.X > 0 {
			h := src.Y * w / src.X
			sz = image.Pt(w, h)
		}
	}

	if !i.valid {
		defer clip.Rect{Max: sz}.Push(gtx.Ops).Pop()
		paint.Fill(gtx.Ops, color.NRGBA{R: 0xFF, A: 0xFF})
		return layout.Dimensions{Size: sz}
	}

	gtx.Constraints = layout.Exact(sz)
	return widget.Image{
		Src: i.Op,
		Fit: widget.Fill,
	}.Layout(gtx)
}
