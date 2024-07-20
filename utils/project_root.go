package utils

import (
	"path"
	"path/filepath"
	"runtime"
	"strings"
)

func FindProjectRoot() string {
	_, b, _, _ := runtime.Caller(0)
	d := path.Join(path.Dir(b))
	return strings.ReplaceAll(filepath.Dir(d), "\\", "/")
}
