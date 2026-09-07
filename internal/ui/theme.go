package ui

import (
	"image/color"

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

	base.Shaper = text.NewShaper(text.WithCollection(nunitoCollection()))
	base.Face = "Nunito Sans"

	return &Theme{
		Base:          base,
		Bg:            hex(0x0A0A0B), // near-black base
		Surface:       hex(0x121214), // panel
		SurfaceRaised: hex(0x1A1A1D), // raised card
		Overlay:       hex(0x232327), // popovers / hover
		Border:        hex(0x2A2A2F), // visible divider
		BorderMuted:   hex(0x1E1E22), // subtle divider
		Text:          hex(0xF4F4F5), // near-white
		Muted:         hex(0xA1A1AA), // secondary text
		Faint:         hex(0x71717A), // tertiary / disabled
		TextOnPrimary: hex(0xFFFFFF), // text on primary fill
		Primary:       hex(0x7C6CF6), // electric indigo
		Danger:        hex(0xF4515C), // red
		Success:       hex(0x36D399), // emerald
		Warning:       hex(0xFBBF24), // amber
		Info:          hex(0x60A5FA), // sky blue

		Radius: unit.Dp(8),
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
