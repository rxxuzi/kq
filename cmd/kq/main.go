// cmd/kq/main.go
package main

import (
	"fmt"
	"log"

	"github.com/rxxuzi/kq/internal/cmd"
	"github.com/spf13/cobra"
)

const VERSION = "0.2.0"

func main() {
	var rootCmd = &cobra.Command{
		Use:   "kq",
		Short: "kq is a lightweight HTML viewer built in Go",
		Run: func(cmd *cobra.Command, args []string) {
			cmd.Help()
		},
	}

	rootCmd.CompletionOptions.DisableDefaultCmd = true

	rootCmd.AddCommand(cmd.NewRunCmd())
	rootCmd.AddCommand(cmd.NewPackCmd())
	rootCmd.AddCommand(&cobra.Command{
		Use:   "version",
		Short: "Displays the version of kq",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Printf("kq version %s\n", VERSION)
		},
	})

	if err := rootCmd.Execute(); err != nil {
		log.Fatalf("Error: %v", err)
	}
}
