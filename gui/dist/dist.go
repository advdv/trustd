// Package guidist allows our Go code to serve the files.
package guidist

import "embed"

// Files exports the bundled frontend code.
//
//go:embed *
var Files embed.FS
