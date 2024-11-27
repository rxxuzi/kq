// Package app internal/app/app.go
package app

import (
	webview "github.com/webview/webview_go"
)

func RunWebView(url string) {
	debug := false // デバッグモードの設定
	w := webview.New(debug)
	defer w.Destroy()
	w.SetTitle("kq Browser")
	w.SetSize(900, 600, webview.HintNone)
	w.Navigate(url)
	w.Run()
}
