// internal/app/app.go
package app

import (
	webview "github.com/webview/webview_go"
)

func RunWebViewWithSize(url string, width, height int) {
	debug := false
	w := webview.New(debug)
	defer w.Destroy()
	w.SetTitle("kq")
	w.SetSize(width, height, webview.HintNone)
	w.Navigate(url)
	w.Run()
}
