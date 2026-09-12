package misc

import (
	"image/color"
	"strconv"
)

func HexColor(hex string) color.NRGBA {
	values, _ := strconv.ParseUint(string(hex[1:]), 16, 32)
	return color.NRGBA{R: uint8(values >> 16), G: uint8((values >> 8) & 0xFF), B: uint8(values & 0xFF), A: 255}
}
