// Package cmd internal/cmd/run.go
package cmd

import (
	"log"
	"os"
	"path/filepath"

	"github.com/rxxuzi/kq/internal/app"
	"github.com/rxxuzi/kq/internal/util"
	"github.com/spf13/cobra"
)

func NewRunCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "run [file]",
		Short: "Starts the viewer with the specified HTML file",
		Long:  `Starts the kq browser, which displays the specified HTML file.`,
		Args:  cobra.ExactArgs(1),
		Run:   run,
	}
}

func run(cmd *cobra.Command, args []string) {
	htmlPath := args[0]
	if _, err := os.Stat(htmlPath); os.IsNotExist(err) {
		log.Fatalf("File does not exist: %s", htmlPath)
	}

	absPath, err := filepath.Abs(htmlPath)
	if err != nil {
		log.Fatalf("Failed to obtain absolute path: %v", err)
	}

	u := util.FileURL(absPath)
	app.RunWebView(u)
}
