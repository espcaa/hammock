//go:build darwin

package misc

/*
#cgo CFLAGS: -x objective-c
#cgo LDFLAGS: -framework Cocoa
#import <Cocoa/Cocoa.h>

static void makeTitlebarTransparent(void *nsviewPtr) {
    NSView *view = (NSView *)nsviewPtr;
    NSWindow *window = [view window];
    if (window == nil) { NSLog(@"hammock: window is nil!"); return; }

    NSLog(@"hammock: window=%@ styleMask before=%lu", window, (unsigned long)window.styleMask);

    window.titlebarAppearsTransparent = YES;
    window.titleVisibility = NSWindowTitleHidden;
    window.styleMask |= NSWindowStyleMaskFullSizeContentView;
    [window setMovableByWindowBackground:YES];

    NSLog(@"hammock: styleMask after=%lu transparent=%d", (unsigned long)window.styleMask, window.titlebarAppearsTransparent);
}

*/
import "C"
import (
	"unsafe"

	"gioui.org/app"
)

func StyleTitlebar(view uintptr) {
	C.makeTitlebarTransparent(unsafe.Pointer(view))
}

func NSViewHandle(e app.ViewEvent) uintptr {
	if ak, ok := e.(app.AppKitViewEvent); ok {
		return uintptr(ak.View)
	}
	return 0
}
