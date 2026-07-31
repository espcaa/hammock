package ui

import (
	"image"
	"image/color"

	"gioui.org/layout"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
	"gioui.org/unit"
	"gioui.org/widget"
)

type Avatar struct {
	Img  image.Image
	Size unit.Dp
}

func (t *Theme) Avatar(img image.Image, size unit.Dp) Avatar {
	return Avatar{Img: img, Size: size}
}

func (a Avatar) Layout(gtx layout.Context) layout.Dimensions {
	sz := gtx.Dp(a.Size)
	gtx.Constraints = layout.Exact(image.Pt(sz, sz))
	rr := sz / 2
	defer clip.RRect{Rect: image.Rectangle{Max: image.Pt(sz, sz)}, NW: rr, NE: rr, SW: rr, SE: rr}.Push(gtx.Ops).Pop()
	if a.Img == nil {
		// if no image, fill with red
		paint.Fill(gtx.Ops, color.NRGBA{R: 0xFF, A: 0xFF})
		return layout.Dimensions{Size: image.Pt(sz, sz)}
	}
	return widget.Image{Src: paint.NewImageOp(a.Img), Fit: widget.Cover}.Layout(gtx)
}
