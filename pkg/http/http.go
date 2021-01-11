package http

import (
	"bytes"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"text/template"

	"github.com/georgemac/truce"
)

type Bindings struct {
	Errors    []Error
	Functions map[string]*Function
}

func BindingsFrom(api truce.API) (Bindings, error) {
	b := Bindings{
		Functions: map[string]*Function{},
	}

	config := &truce.HTTP{Versions: []string{"1.0", "1.1", "2.0"}}
	if api.Transports.HTTP != nil {
		config = api.Transports.HTTP
	}

	tmpl, err := template.New("prefix").Parse(config.Prefix)
	if err != nil {
		return b, err
	}

	buf := &bytes.Buffer{}
	if err := tmpl.Execute(buf, api); err != nil {
		return b, err
	}

	config.Prefix = buf.String()

	for _, f := range api.Functions {
		fn, err := NewFunction(config, f)
		if err != nil {
			return b, err
		}

		if fn != nil {
			b.Functions[f.Name] = fn
		}
	}

	for _, err := range config.Errors {
		t, ok := api.Types[err.Type]
		if !ok {
			return b, errors.New("cannot locate type definition for transport error")
		}

		if t.Type != "error" {
			return b, errors.New("transport error type definition is not error")
		}

		code, err := strconv.ParseInt(err.StatusCode, 10, 64)
		if err != nil {
			return b, fmt.Errorf("parsing status code: %w", err)
		}

		b.Errors = append(b.Errors, Error{
			Definition: t,
			StatusCode: int(code),
		})
	}

	return b, nil
}

type Error struct {
	Definition truce.Type
	StatusCode int
}

type Function struct {
	Definition  truce.Function
	Method      string
	Path        Path
	Query       map[string]string
	BodyVar     string
	BodyType    string
	HasReturn   bool
	ReturnType  string
	ReturnIsPtr bool
}

type Path []Element

func parsePath(vars map[string]string, v string) (p Path) {
	for i, val := range strings.Split(v, "/") {
		if i == 0 && val == "" {
			continue
		}

		var (
			typ  = "static"
			varn string
		)
		if len(val) > 0 && val[0] == '{' && val[len(val)-1] == '}' {
			typ = "variable"
			val = val[1 : len(val)-1]
			varn = vars[val]
		}

		p = append(p, Element{
			Type:  typ,
			Value: val,
			Var:   varn,
		})
	}

	return
}

func (p Path) String() (v string) {
	for _, e := range p {
		v += "/" + e.String()
	}
	return
}

func (p Path) FmtString() (v string) {
	for _, e := range p {
		v += "/" + e.FmtString()
	}
	return
}

func (p Path) ArgString() (v string) {
	var i int
	for _, e := range p {
		if e.Type != "variable" {
			continue
		}
		if i > 0 {
			v += ", "
		}
		v += e.Var
		i++
	}
	return
}

type Element struct {
	Type  string
	Value string
	Var   string
}

func (e Element) String() string {
	switch e.Type {
	case "static":
		return e.Value
	case "variable":
		return "{" + e.Value + "}"
	default:
		panic("element type not supported")
	}
}

func (e Element) FmtString() string {
	switch e.Type {
	case "static":
		return e.Value
	case "variable":
		return "%v"
	default:
		panic("element type not supported")
	}
}

func NewFunction(config *truce.HTTP, function truce.Function) (*Function, error) {
	if function.Transports.HTTP == nil {
		return nil, nil
	}

	transport := *function.Transports.HTTP

	b := &Function{Definition: function, Query: map[string]string{}}

	type argument struct {
		variable string
		typ      string
	}

	var (
		pathMappings = map[string]string{}
		args         = map[string]argument{}
	)

	for i, field := range function.Arguments {
		args[field.Name] = argument{
			typ:      string(field.Type),
			variable: fmt.Sprintf("v%d", i),
		}
	}

	if function.Return.Name != "" {
		b.HasReturn = true
		b.ReturnType = string(function.Return.Type)

		if len(b.ReturnType) < 1 {
			return nil, errors.New("return type cannot be empty")
		}

		if b.ReturnType[0] == '*' {
			b.ReturnType = b.ReturnType[1:]
			b.ReturnIsPtr = true
		}
	}

	b.Method = transport.Method

	for _, arg := range transport.Arguments {
		a, ok := args[arg.Name]

		switch arg.From {
		case "body":
			if !ok {
				continue
			}

			b.BodyVar = a.variable
			b.BodyType = a.typ
		case "path":
			if !ok {
				continue
			}

			pathMappings[arg.Var] = args[arg.Name].variable
		case "query":
			if !ok {
				continue
			}

			b.Query[arg.Var] = args[arg.Name].variable
		case "static":
			// TODO(georgemac)
		}
	}

	for _, part := range strings.Split(config.Prefix, "/") {
		if part == "" {
			continue
		}

		b.Path = append(b.Path, Element{Type: "static", Value: part})
	}

	b.Path = append(b.Path, parsePath(pathMappings, transport.Path)...)

	return b, nil
}
