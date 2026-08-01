package misc

import (
	"image"
	"image/color"
)

func Darken(c color.NRGBA, factor float32) color.NRGBA {
	return color.NRGBA{
		R: uint8(float32(c.R) * factor),
		G: uint8(float32(c.G) * factor),
		B: uint8(float32(c.B) * factor),
		A: c.A,
	}
}

func Lighten(c color.NRGBA, amount float32) color.NRGBA {
	lerp := func(v uint8) uint8 {
		return uint8(float32(v) + (255-float32(v))*amount)
	}
	return color.NRGBA{R: lerp(c.R), G: lerp(c.G), B: lerp(c.B), A: c.A}
}

func GradientImage(width, height int, startColor, endColor color.NRGBA) *image.RGBA {
	img := image.NewRGBA(image.Rect(0, 0, width, height))
	for y := range img.Bounds().Dy() {
		t := float32(y) / float32(height-1)
		r := uint8(float32(startColor.R)*(1-t) + float32(endColor.R)*t)
		g := uint8(float32(startColor.G)*(1-t) + float32(endColor.G)*t)
		b := uint8(float32(startColor.B)*(1-t) + float32(endColor.B)*t)
		a := uint8(float32(startColor.A)*(1-t) + float32(endColor.A)*t)
		for x := range img.Bounds().Dx() {
			img.Set(x, y, color.NRGBA{R: r, G: g, B: b, A: a})
		}
	}
	return img
}
