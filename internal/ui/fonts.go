package ui

import (
	_ "embed"

	"gioui.org/font"
	"gioui.org/font/opentype"
)

//go:embed fonts/NunitoSans-Regular.ttf
var nsRegular []byte

//go:embed fonts/NunitoSans-Bold.ttf
var nsBold []byte

//go:embed fonts/NunitoSans-Italic.ttf
var nsItalic []byte

//go:embed fonts/NunitoSans-BoldItalic.ttf
var nsBoldItalic []byte

func mustParse(b []byte) opentype.Face {
	f, err := opentype.Parse(b)
	if err != nil {
		panic("bad font: " + err.Error())
	}
	return f
}

func nunitoCollection() []font.FontFace {
	const ui = "Nunito Sans"
	const mono = "JetBrains Mono"
	const emoji = "Noto Color Emoji"
	return []font.FontFace{
		{Font: font.Font{Typeface: ui, Weight: font.Normal}, Face: mustParse(nsRegular)},
		{Font: font.Font{Typeface: ui, Weight: font.Bold}, Face: mustParse(nsBold)},
		{Font: font.Font{Typeface: ui, Weight: font.Normal, Style: font.Italic}, Face: mustParse(nsItalic)},
		{Font: font.Font{Typeface: ui, Weight: font.Bold, Style: font.Italic}, Face: mustParse(nsBoldItalic)},
	}
}
