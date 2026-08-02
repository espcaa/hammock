package login

import (
	webview "github.com/webview/webview_go"
)

func RunLoginWebview(authURL string) string {
	w := webview.New(true)
	defer w.Destroy()

	w.SetTitle("Sign in to Hammock")
	w.SetSize(480, 640, webview.HintNone)

	tokenCh := make(chan string, 1)

	w.Init(`
		window.reportToken = function(t) { window.__go_onToken(t); };
	`)

	w.Bind("__go_onToken", func(token string) {
		tokenCh <- token
		w.Terminate()
	})

	w.Navigate(authURL)

	go func() {

	}()

	w.Run()

	select {
	case tok := <-tokenCh:
		return tok
	default:
		return ""
	}
}
