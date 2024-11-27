// Package cmd internal/cmd/pack.go
package cmd

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/rxxuzi/kq/internal/types"
	"github.com/spf13/cobra"
)

var (
	sizeFlag    string
	entryFlag   string
	outputFlag  string
	previewFlag bool
	minifyFlag  bool
)

// NewPackCmd creates a new pack command
func NewPackCmd() *cobra.Command {
	packCmd := &cobra.Command{
		Use:   "pack <path> [<path2> ...]",
		Short: "Generate a JSON configuration file from specified paths",
		Long: `Generate a JSON configuration file from specified files or directories.
If multiple paths are provided, all files within those paths will be included as sources.`,
		Args: cobra.MinimumNArgs(1),
		Run:  pack,
	}

	// Adding flags
	packCmd.Flags().StringVarP(&sizeFlag, "size", "s", "800x600", "Window size in WIDTHxHEIGHT format")
	packCmd.Flags().StringVarP(&entryFlag, "entry", "e", "", "Entry point file (default is index.html)")
	packCmd.Flags().StringVarP(&outputFlag, "output", "o", "a.kq", "Output file name (default is a.kq)")
	packCmd.Flags().BoolVarP(&previewFlag, "preview", "p", false, "Preview JSON in the command line without saving to a file")
	packCmd.Flags().BoolVarP(&minifyFlag, "minify", "m", false, "Minify the JSON output")

	return packCmd
}

func pack(cmd *cobra.Command, args []string) {
	paths := args
	var absPaths []string

	// Convert all provided paths to absolute paths
	for _, path := range paths {
		absPath, err := filepath.Abs(path)
		if err != nil {
			log.Fatalf("Failed to get absolute path for '%s': %v", path, err)
		}
		absPaths = append(absPaths, absPath)
	}

	// Set application name to the base name of the first path
	name := filepath.Base(absPaths[0])
	if name == "" {
		name = "app"
	}

	// Handle --entry flag
	entry := entryFlag
	if entry == "" {
		entry = "index.html"
	}

	// Parse window size
	width, height, err := parseSize(sizeFlag)
	if err != nil {
		log.Fatalf("Invalid size format: %v", err)
	}

	// Collect source files from all provided paths
	sources, err := collectSources(absPaths)
	if err != nil {
		log.Fatalf("Failed to collect source files: %v", err)
	}

	// Resolve and verify the entry file
	entryAbsPath, err := resolveEntryPath(entry, absPaths, sources)
	if err != nil {
		log.Fatalf("Entry file not found: %v", err)
	}

	// Create the configuration structure
	config := types.Config{
		Name:    name,
		Version: "0.1.0",
		Source:  sources,
		Allow: []string{
			"*.html", "*.css", "*.js", "*.png", "*.pdf", "*.jpg", "*.svg",
		},
		Entry: entryAbsPath,
		Options: types.Options{
			Window: types.WindowOptions{
				Height:    height,
				Width:     width,
				Frameless: false,
				Resizable: true,
			},
			Security: types.SecurityOptions{
				NoScript:    false,
				LocalOnly:   false,
				AllowOrigin: []string{},
			},
			Debug: types.DebugOptions{
				Devtools: false,
				Console:  false,
			},
			Env: types.EnvOptions{
				SingleInstance: false,
			},
			UI: types.UIOptions{
				Theme: "light",
				Title: "kq Browser",
				Icon:  nil,
			},
		},
	}

	// Serialize the configuration to JSON
	var jsonData []byte
	if minifyFlag {
		jsonData, err = json.Marshal(config)
	} else {
		jsonData, err = json.MarshalIndent(config, "", "  ")
	}
	if err != nil {
		log.Fatalf("Failed to serialize JSON: %v", err)
	}

	if previewFlag {
		// If preview flag is set, print JSON to stdout
		fmt.Println(string(jsonData))
	} else {
		// Otherwise, write JSON to the specified output file
		err = os.WriteFile(outputFlag, jsonData, 0644)
		if err != nil {
			log.Fatalf("Failed to write to file '%s': %v", outputFlag, err)
		}
		fmt.Printf("Configuration file written to '%s'\n", outputFlag)
	}
}

// parseSize parses the size flag in WIDTHxHEIGHT format
func parseSize(sizeStr string) (width int, height int, err error) {
	parts := strings.Split(sizeStr, "x")
	if len(parts) != 2 {
		return 0, 0, fmt.Errorf("size must be in WIDTHxHEIGHT format")
	}
	width, err = strconv.Atoi(parts[0])
	if err != nil {
		return 0, 0, fmt.Errorf("invalid width: %v", err)
	}
	height, err = strconv.Atoi(parts[1])
	if err != nil {
		return 0, 0, fmt.Errorf("invalid height: %v", err)
	}
	return width, height, nil
}

// collectSources collects all source files from the provided absolute paths
func collectSources(paths []string) ([]string, error) {
	var sources []string

	for _, path := range paths {
		info, err := os.Stat(path)
		if err != nil {
			return nil, fmt.Errorf("failed to access '%s': %v", path, err)
		}

		if info.IsDir() {
			err = filepath.Walk(path, func(p string, info os.FileInfo, err error) error {
				if err != nil {
					return err
				}
				if !info.IsDir() {
					absFilePath, err := filepath.Abs(p)
					if err != nil {
						return err
					}
					sources = append(sources, absFilePath)
				}
				return nil
			})
			if err != nil {
				return nil, fmt.Errorf("failed to walk directory '%s': %v", path, err)
			}
		} else {
			absFilePath, err := filepath.Abs(path)
			if err != nil {
				return nil, fmt.Errorf("failed to get absolute path for '%s': %v", path, err)
			}
			sources = append(sources, absFilePath)
		}
	}

	return sources, nil
}

// resolveEntryPath verifies the entry file exists within the collected sources and returns its absolute path
func resolveEntryPath(entry string, paths []string, sources []string) (string, error) {
	var entryAbsPath string

	// If entry contains a path separator, treat it as a path
	if strings.Contains(entry, string(os.PathSeparator)) {
		absEntry, err := filepath.Abs(entry)
		if err != nil {
			return "", fmt.Errorf("failed to get absolute path for entry '%s': %v", entry, err)
		}

		// Check if the absolute entry path is in sources
		for _, src := range sources {
			if src == absEntry {
				entryAbsPath = src
				break
			}
		}

		if entryAbsPath == "" {
			return "", fmt.Errorf("entry file '%s' does not exist in the specified paths", absEntry)
		}
	} else {
		// Otherwise, treat it as a filename and search in sources
		for _, src := range sources {
			if strings.EqualFold(filepath.Base(src), entry) { // Case-insensitive comparison for Windows
				entryAbsPath = src
				break
			}
		}

		if entryAbsPath == "" {
			return "", fmt.Errorf("entry file '%s' does not exist in the specified paths", entry)
		}
	}

	return entryAbsPath, nil
}
