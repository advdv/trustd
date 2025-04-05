// Package gui allows our Go code to serve the files.
package gui

import (
	"embed"
	"io/fs"
)

// Files exports the bundled frontend code.
//
//go:embed dist/*
var embedded embed.FS

// Dist exports the dist files as a filesystem.
var Dist fs.FS

func init() {
	Dist, _ = fs.Sub(embedded, "dist")
}
