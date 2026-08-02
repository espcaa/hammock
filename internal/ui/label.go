package ui

import (
	"image/color"

	"gioui.org/font"
	"gioui.org/layout"
	"gioui.org/unit"
	"gioui.org/widget/material"
)

type LabelStyle struct {
	Text   string
	Size   unit.Sp
	Color  color.NRGBA
	Weight font.Weight
	Italic bool
	th     *material.Theme
}

func (t *Theme) Label(text string, size unit.Sp, weight font.Weight, italic bool) LabelStyle {
	return LabelStyle{
		Text:   text,
		Size:   size,
		Color:  t.Text,
		Weight: weight,
		Italic: italic,
		th:     t.Base,
	}
}

func (l LabelStyle) Layout(gtx layout.Context) layout.Dimensions {
	lbl := material.Label(l.th, l.Size, l.Text)
	lbl.Color = l.Color
	lbl.Font.Weight = l.Weight
	if l.Italic {
		lbl.Font.Style = font.Italic
	}
	return lbl.Layout(gtx)
}
