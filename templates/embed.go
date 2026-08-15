package templates

import "embed"

//go:embed *.tmpl
var FS embed.FS

//go:embed static
var StaticFS embed.FS
