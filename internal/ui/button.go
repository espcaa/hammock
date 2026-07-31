package ui

import (
	"image"

	"gioui.org/layout"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
	"gioui.org/widget"
	"gioui.org/widget/material"
	"github.com/espcaa/hammock/internal/misc"
)

type Button struct {
	Click widget.Clickable
}

type ButtonStyle struct {
	th   *Theme
	btn  *Button
	Text string
}

func (t *Theme) Button(btn *Button, txt string) ButtonStyle {
	return ButtonStyle{th: t, btn: btn, Text: txt}
}

func (b ButtonStyle) Layout(gtx layout.Context) layout.Dimensions {
	return b.btn.Click.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
		return layout.Stack{}.Layout(gtx,
			layout.Expanded(func(gtx layout.Context) layout.Dimensions {
				fill := b.th.SurfaceRaised
				if b.btn.Click.Hovered() {
					fill = misc.Lighten(fill, 0.06)
				}
				if b.btn.Click.Pressed() {
					fill = misc.Darken(fill, 0.4)
				}
				rr := gtx.Dp(b.th.Radius)
				defer clip.RRect{
					Rect: image.Rectangle{Max: gtx.Constraints.Min},
					SE:   rr, SW: rr, NE: rr, NW: rr,
				}.Push(gtx.Ops).Pop()
				paint.Fill(gtx.Ops, fill)

				// border
				paint.FillShape(gtx.Ops, b.th.Border,
					clip.Stroke{
						Path: clip.RRect{
							Rect: image.Rectangle{Max: gtx.Constraints.Min},
							SE:   rr, SW: rr, NE: rr, NW: rr,
						}.Path(gtx.Ops),
						Width: float32(gtx.Dp(1)),
					}.Op(),
				)
				return layout.Dimensions{Size: gtx.Constraints.Min}
			}),

			layout.Stacked(func(gtx layout.Context) layout.Dimensions {
				return layout.UniformInset(b.th.Gutter).Layout(gtx, func(gtx layout.Context) layout.Dimensions {
					lbl := material.Body1(b.th.Base, b.Text)
					lbl.Color = b.th.Text
					return lbl.Layout(gtx)
				})
			}),
		)
	})
}
