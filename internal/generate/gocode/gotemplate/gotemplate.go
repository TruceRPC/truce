package gotemplate

import (
	"embed"
	"text/template"
)

//go:embed tmpl/*.go.tmpl
var templates embed.FS

// ParseAll parses and returns all templates required for Go generation.
func ParseAll() (*template.Template, error) {
	return template.
		New("").
		Funcs(tmplFuncs).
		ParseFS(templates, "tmpl/*.go.tmpl")
}
