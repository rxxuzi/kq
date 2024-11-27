// main.go
package main

import (
	"fmt"
	"github.com/spf13/cobra"
	webview "github.com/webview/webview_go"
	"log"
	"net/url"
	"os"
	"path/filepath"
	"runtime"
)

const VERSION = "0.1.0"

func main() {
	var rootCmd = &cobra.Command{
		Use:   "kq",
		Short: "kq is a lightweight HTML viewer built with Go",
		Long:  `kq is a lightweight and fast HTML viewer built using Go and WebView.`,
		Run: func(cmd *cobra.Command, args []string) {
			cmd.Help()
		},
	}

	rootCmd.CompletionOptions.DisableDefaultCmd = true

	var runCmd = &cobra.Command{
		Use:   "run [file]",
		Short: "Run the HTML viewer with the specified HTML file",
		Long:  `Run launches the kq browser and displays the specified HTML file.`,
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			htmlPath := args[0]
			if _, err := os.Stat(htmlPath); os.IsNotExist(err) {
				log.Fatalf("File does not exist: %s", htmlPath)
			}

			absPath, err := filepath.Abs(htmlPath)
			if err != nil {
				log.Fatalf("Failed to get absolute path: %v", err)
			}

			u := fileURL(absPath)
			fmt.Printf("Opening: %s\n", u) // デバッグ用の出力

			debug := false // debug = off
			w := webview.New(debug)
			defer w.Destroy()
			w.SetTitle("kq Browser")
			w.SetSize(900, 600, webview.HintNone)
			w.Navigate(u)
			w.Run()
		},
	}

	var versionCmd = &cobra.Command{
		Use:   "version",
		Short: "Show the version of kq",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Printf("kq version %s\n", VERSION)
		},
	}

	rootCmd.AddCommand(runCmd)
	rootCmd.AddCommand(versionCmd)

	if err := rootCmd.Execute(); err != nil {
		log.Fatalf("Error: %v", err)
	}
}

func fileURL(path string) string {
	path = filepath.ToSlash(path)

	if runtime.GOOS == "windows" {
		// Windowsの場合、先頭にスラッシュを追加
		path = "/" + path
	}

	u := url.URL{
		Scheme: "file",
		Path:   path,
	}

	return u.String()
}
