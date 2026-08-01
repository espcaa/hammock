package app

import "gioui.org/layout"

type Screen interface {
	Layout(gtx layout.Context) layout.Dimensions
}
