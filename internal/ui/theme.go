package ui

import (
	"image/color"

	"gioui.org/font/gofont"
	"gioui.org/text"
	"gioui.org/unit"
	"gioui.org/widget/material"
)

type Theme struct {
	Base          *material.Theme
	Bg            color.NRGBA
	Surface       color.NRGBA
	SurfaceRaised color.NRGBA
	Overlay       color.NRGBA
	Border        color.NRGBA
	BorderMuted   color.NRGBA
	Text          color.NRGBA
	Muted         color.NRGBA
	Faint         color.NRGBA
	TextOnPrimary color.NRGBA
	Primary       color.NRGBA
	Danger        color.NRGBA
	Success       color.NRGBA
	Warning       color.NRGBA
	Info          color.NRGBA

	Radius unit.Dp
	Gutter unit.Dp
}

func NewTheme() *Theme {
	base := material.NewTheme()
	base.Shaper = text.NewShaper(text.WithCollection(gofont.Collection()))
	return &Theme{
		Base:          base,
		Bg:            hex(0x1D2021), // bg0_hard
		Surface:       hex(0x282828), // bg0
		SurfaceRaised: hex(0x32302F), // bg0_soft
		Overlay:       hex(0x3C3836), // bg1
		Border:        hex(0x504945), // bg2
		BorderMuted:   hex(0x3C3836), // bg1
		Text:          hex(0xEBDBB2), // fg1
		Muted:         hex(0xA89984), // fg4
		Faint:         hex(0x928374), // gray
		TextOnPrimary: hex(0x1D2021), // bg0_hard
		Primary:       hex(0x83A598), // bright blue
		Danger:        hex(0xFB4934), // bright red
		Success:       hex(0xB8BB26), // bright green
		Warning:       hex(0xFABD2F), // bright yellow
		Info:          hex(0x8EC07C), // bright aqua

		Radius: unit.Dp(2),
		Gutter: unit.Dp(8),
	}
}

func hex(v uint32) color.NRGBA {
	return color.NRGBA{R: byte(v >> 16), G: byte(v >> 8), B: byte(v), A: 0xFF}
}

func toHex(c color.NRGBA) string {
	return "#" + hexByte(c.R) + hexByte(c.G) + hexByte(c.B)
}

func hexByte(b byte) string {
	const hex = "0123456789ABCDEF"
	return string([]byte{hex[b>>4], hex[b&0x0F]})
}
