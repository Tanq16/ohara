package server

import "embed"

//go:embed index.html
var indexHTML []byte

//go:embed all:static
var staticFiles embed.FS
