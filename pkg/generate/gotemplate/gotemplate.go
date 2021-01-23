package gotemplate

import (
	"embed"
	"text/template"
)

//go:embed tmpl/*.go.tmpl
var templates embed.FS

func ParseAll() (*template.Template, error) {
	return template.
		New("").
		Funcs(tmplFuncs).
		ParseFS(templates, "tmpl/*.go.tmpl")
}
