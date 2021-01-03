package http

import (
	"bytes"
	"fmt"
	"strings"
	"text/template"

	"github.com/georgemac/truce"
)

type Bindings map[string]*Binding

func BindingsFrom(api truce.API) Bindings {
	b := Bindings{}

	config := &truce.HTTP{Versions: []string{"1.0", "1.1", "2.0"}}
	if api.Transports.HTTP != nil {
		config = api.Transports.HTTP
	}

	tmpl, err := template.New("prefix").Parse(config.Prefix)
	if err != nil {
		panic(err)
	}

	buf := &bytes.Buffer{}
	if err := tmpl.Execute(buf, api); err != nil {
		panic(err)
	}

	config.Prefix = buf.String()

	for _, f := range api.Functions {
		if binding := NewBinding(config, f); binding != nil {
			b[f.Name] = binding
		}
	}

	return b
}

type Binding struct {
	Function  truce.Function
	Method    string
	Path      Path
	BodyVar   string
	BodyType  string
	HasReturn bool
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

func NewBinding(config *truce.HTTP, function truce.Function) *Binding {
	if function.Transports.HTTP == nil {
		return nil
	}

	transport := *function.Transports.HTTP

	b := &Binding{Function: function}

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
	}

	b.Method = transport.Method

	for _, arg := range transport.Arguments {
		if len(arg.Value) == 0 {
			continue
		}

		if arg.Value[0] != '$' {
			continue
		}

		target := arg.Value
		n := strings.Index(target, ".")
		if n > 0 {
			target = target[:n]
		}

		switch target {
		case "$body":
			if a, ok := args[arg.Name]; ok {
				b.BodyVar = a.variable
				b.BodyType = a.typ
			}
		case "$path":
			// given a path variable is targetted
			// create reverse lookup entry in path map
			if len(arg.Value) > n {
				pathMappings[arg.Value[n+1:]] = args[arg.Name].variable
			}
		}
	}

	for _, part := range strings.Split(config.Prefix, "/") {
		if part == "" {
			continue
		}

		b.Path = append(b.Path, Element{Type: "static", Value: part})
	}

	b.Path = append(b.Path, parsePath(pathMappings, transport.Path)...)

	return b
}
