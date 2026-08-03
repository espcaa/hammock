package login

import (
	"regexp"
	"strings"

	webview "github.com/espcaa/webview_go"
)

func RunLoginWebview(authURL string) string {
	w := webview.New(true)
	defer w.Destroy()

	w.SetTitle("Sign in")
	w.SetUserAgent("Mozilla/5.0 (Macintosh; Intel Mac OS X 15_6_0) AppleWebKit/537.36 (KHTML, like Gecko) Slack/4.48.102 Chrome/144.0.7559.236 Electron/40.8.2 Safari/537.36")
	w.SetSize(480, 640, webview.HintNone)

	capturedUrl := ""

	w.Navigate(authURL)

	w.OnNavigation(func(url string) bool {
		if strings.HasPrefix(url, "slack://") {
			capturedUrl = url
			w.Terminate()
			return false
		}
		return true
	})

	go func() {

	}()

	w.Run()

	// get the magic code from the captured URL
	regex := `magic-login/([^?]+)`
	code := ""
	if capturedUrl != "" {
		re := regexp.MustCompile(regex)
		if m := re.FindStringSubmatch(capturedUrl); m != nil {
			code = m[1]
		}
	}

	return code
}
