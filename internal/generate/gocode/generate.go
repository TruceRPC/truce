// Package gocode provides Go code generation for Truce specifications.
package gocode

import (
	"bytes"
	"go/format"
	"io"
	"text/template"

	"github.com/TruceRPC/truce"
	"github.com/TruceRPC/truce/internal/generate/gocode/gotemplate"
)

// Generator generates a go client, server and the necessary associated types
// for them to do their work.
type Generator struct {
	api      truce.API
	data     templatedata
	template *template.Template
}

type templatedata struct {
	Package    string
	Imports    string
	ClientName string
	ServerName string
	Types      map[string]truce.Type
	Functions  map[string]*gotemplate.Function
	Errors     []gotemplate.Error
}

// New creates a new Generator.
func New(api truce.API, opts ...Option) (Generator, error) {
	bindings, err := gotemplate.BindingsFrom(api)
	if err != nil {
		return Generator{}, err
	}

	t, err := gotemplate.ParseAll()
	if err != nil {
		return Generator{}, err
	}

	g := Generator{
		api:      api,
		template: t,
		data: templatedata{
			Package:    "types",
			ClientName: "Client",
			ServerName: "Server",
			Types:      api.Types,
			Functions:  bindings.Functions,
			Errors:     bindings.Errors,
		},
	}

	Options(opts).Apply(&g)

	return g, nil
}

// WriteTypesTo generates a Go types (for server and client) and writes the file
// content to an io.Writer. The generated code is gofmt'd.
func (g Generator) WriteTypesTo(w io.Writer, opts ...Option) error {
	Options(opts).Apply(&g)
	t := g.template.Lookup("types.go.tmpl")
	return writeGo(w, t, g.data)
}

// WriteClientTo generates a Go client and writes the file content to an
// io.Writer. The generated code is gofmt'd.
func (g Generator) WriteClientTo(w io.Writer, opts ...Option) error {
	Options(opts).Apply(&g)
	t := g.template.Lookup("client.go.tmpl")
	return writeGo(w, t, g.data)
}

// WriteServerTo generates a Go server and writes the file content to an
// io.Writer. The generated code is gofmt'd.
func (g Generator) WriteServerTo(w io.Writer, opts ...Option) error {
	Options(opts).Apply(&g)
	t := g.template.Lookup("server.go.tmpl")
	return writeGo(w, t, g.data)
}

func writeGo(w io.Writer, tmpl *template.Template, d templatedata) error {
	buf := &bytes.Buffer{}
	if err := tmpl.Execute(buf, d); err != nil {
		return err
	}

	src, err := format.Source(buf.Bytes())
	if err != nil {
		return err
	}

	_, err = w.Write(src)
	return err
}
