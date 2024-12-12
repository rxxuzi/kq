// cmd/kq/main.go
package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/rxxuzi/kq/internal/app"
	"github.com/rxxuzi/kq/internal/util"
)

func main() {
	// フラグの定義
	sizeFlag := flag.String("size", "900x600", "Window size in WIDTHxHEIGHT format")
	flag.Parse()

	// 必須の引数を確認
	if flag.NArg() < 1 {
		fmt.Println("Usage: kq [--size WIDTHxHEIGHT] <file.html>")
		os.Exit(1)
	}

	// 引数からHTMLファイルのパスを取得
	htmlPath := flag.Arg(0)
	if _, err := os.Stat(htmlPath); os.IsNotExist(err) {
		log.Fatalf("File does not exist: %s", htmlPath)
	}

	// 絶対パスを取得
	absPath, err := filepath.Abs(htmlPath)
	if err != nil {
		log.Fatalf("Failed to obtain absolute path: %v", err)
	}
	width, height, err := parseSize(*sizeFlag)
	if err != nil {
		log.Fatalf("Invalid size format: %v", err)
	}
	url := util.FileURL(absPath)
	app.RunWebViewWithSize(url, width, height)
}

func parseSize(sizeStr string) (int, int, error) {
	parts := strings.Split(sizeStr, "x")
	if len(parts) != 2 {
		return 0, 0, fmt.Errorf("size must be in WIDTHxHEIGHT format")
	}
	width, err := strconv.Atoi(parts[0])
	if err != nil {
		return 0, 0, fmt.Errorf("invalid width: %v", err)
	}
	height, err := strconv.Atoi(parts[1])
	if err != nil {
		return 0, 0, fmt.Errorf("invalid height: %v", err)
	}
	return width, height, nil
}
