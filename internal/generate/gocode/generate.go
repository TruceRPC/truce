package gocode

import (
	"bytes"
	"go/format"
	"io"
	"text/template"

	"github.com/TruceRPC/truce"
	"github.com/TruceRPC/truce/internal/generate/gocode/gotemplate"
)

type Generator struct {
	api truce.API

	context  context
	template *template.Template
}

type context struct {
	Package    string
	Imports    string
	ClientName string
	ServerName string
	Types      map[string]truce.Type
	Functions  map[string]*gotemplate.Function
	Errors     []gotemplate.Error
}

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
		context: context{
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

func (g Generator) WriteTypesTo(w io.Writer, opts ...Option) error {
	Options(opts).Apply(&g)
	t := g.template.Lookup("type.go.tmpl")
	return writeGo(w, t, g.context)
}

func (g Generator) WriteClientTo(w io.Writer, opts ...Option) error {
	Options(opts).Apply(&g)
	t := g.template.Lookup("client.go.tmpl")
	return writeGo(w, t, g.context)
}

func (g Generator) WriteServerTo(w io.Writer, opts ...Option) error {
	Options(opts).Apply(&g)
	t := g.template.Lookup("server.go.tmpl")
	return writeGo(w, t, g.context)
}

func writeGo(w io.Writer, tmpl *template.Template, context interface{}) error {
	buf := &bytes.Buffer{}
	if err := tmpl.Execute(buf, context); err != nil {
		return err
	}

	src, err := format.Source(buf.Bytes())
	if err != nil {
		return err
	}

	_, err = w.Write(src)
	return err
}
