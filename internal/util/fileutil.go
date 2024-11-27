// Package util internal/util/fileutil.go
package util

import (
	"net/url"
	"path/filepath"
	"runtime"
)

func FileURL(path string) string {
	path = filepath.ToSlash(path)

	if runtime.GOOS == "windows" {
		path = "/" + path
	}

	u := url.URL{
		Scheme: "file",
		Path:   path,
	}

	return u.String()
}
