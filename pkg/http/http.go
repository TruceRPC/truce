package http

import (
	"fmt"
	"path"
	"strings"

	"github.com/georgemac/truce"
)

type Bindings map[string]*Binding

func BindingsFrom(api truce.API) Bindings {
	b := Bindings{}

	for _, f := range api.Functions {
		b[f.Name] = NewBinding(f)
	}

	return b
}

type Binding struct {
	Function  truce.Function
	Method    string
	Path      string
	PathFmt   string
	PathArgs  string
	BodyArg   string
	HasReturn bool
}

func NewBinding(function truce.Function) *Binding {
	b := &Binding{Function: function}

	var (
		pathMappings = map[string]string{}
		args         = map[string]string{}
	)

	for i, field := range function.Arguments {
		args[field.Name] = fmt.Sprintf("v%d", i)
	}

	if function.Return.Name != "" {
		b.HasReturn = true
	}

	transport, ok := b.Function.Transports[0].Config.(truce.HTTPFunction)
	if !ok {
		panic("unexpected type")
	}

	b.Path = transport.Path
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
			b.BodyArg = args[arg.Name]
		case "$path":
			// given a path variable is targetted
			// create reverse lookup entry in path map
			if len(arg.Value) > n {
				pathMappings[arg.Value[n+1:]] = args[arg.Name]
			}
		}
	}

	var (
		parts    = strings.Split(transport.Path, "/")
		pathArgs []string
	)
	for i, part := range parts {
		if len(part) > 0 && part[0] == '{' && part[len(part)-1] == '}' {
			// replace part with substitution
			parts[i] = "%v"

			// lookup argument reference in path mappings
			argName, ok := pathMappings[part[1:len(part)-1]]
			if !ok {
				panic("argument not defined")
			}

			pathArgs = append(pathArgs, argName)
		}
	}

	b.PathFmt = path.Join(append([]string{"/"}, parts...)...)
	b.PathArgs = strings.Join(pathArgs, ",")

	return b
}
