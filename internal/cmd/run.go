// internal/app/app.go
package app

import (
	webview "github.com/webview/webview_go"
)

// RunWebViewWithSize は指定されたURLとウィンドウサイズでWebViewを起動します。
func RunWebViewWithSize(url string, width, height int) {
	debug := false // デバッグモードの設定
	w := webview.New(debug)
	defer w.Destroy()
	w.SetTitle("kq Browser")
	w.SetSize(width, height, webview.HintNone)
	w.Navigate(url)
	w.Run()
}
