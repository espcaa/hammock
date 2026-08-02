//go:build !darwin

package misc

func StyleTitlebar(view uintptr) {}

func NSViewHandle(e any) uintptr { return 0 }
