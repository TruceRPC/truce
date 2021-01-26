package gotemplate

import (
	"bytes"
	"errors"
	"fmt"
	"sort"
	"strconv"
	"strings"
	"text/template"

	"github.com/TruceRPC/truce"
)

// Bindings contains Go-specific details derived from a generalized Truce API
// specification.
type Bindings struct {
	Errors    []Error
	Functions map[string]*Function
}

// BindingsFrom derives a set of Bindings from the underlying API specification.
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

	// Sort the errors by status code so we have a deterministic order going
	// into any template phases.
	sort.Slice(b.Errors, func(i, j int) bool {
		return b.Errors[i].StatusCode < b.Errors[j].StatusCode
	})

	return b, nil
}

// Error is an error that can inform an HTTP response.
type Error struct {
	Definition truce.Type
	StatusCode int
}

// QueryParam represent a parameter mapped to a query variable
// It contains helper method to parse query parameters into go
// variables and vice-versa.
type QueryParam struct {
	Pos      int
	Type     string
	GoVar    string
	QueryVar string
}

// ToStringVar outputs the code necessary to coerce the Go variable
// onto the respective query parameter.
func (q QueryParam) ToStringVar() string {
	switch q.Type {
	case "string":
		return q.GoVar
	case "timestamp":
		return fmt.Sprintf("%s.Format(time.RFC3339)", q.GoVar)
	}

	return fmt.Sprintf("fmt.Sprintf(\"%%v\", %s)", q.GoVar)
}

// FromStringVar
func (q QueryParam) FromStringVar() (v string) {
	v = fmt.Sprintf("q%d := r.URL.Query().Get(\"%s\")\n", q.Pos, q.QueryVar)
	switch q.Type {
	case "int":
		v += fmt.Sprintf("%s, err := strconv.ParseInt(q%d, 10, 64); if err != nil { return }\n",
			q.GoVar, q.Pos)
	case "float64":
		v += fmt.Sprintf("%s, err := strconv.ParseFloat(q%d, 10, 64); if err != nil { return }\n",
			q.GoVar, q.Pos)
	case "bool":
		v += fmt.Sprintf("%s, err := strconv.ParseBool(q%d); if err != nil { return }\n",
			q.GoVar, q.Pos)
	case "timestamp":
		v += fmt.Sprintf("%s, err := time.Parse(time.RFC3339, q%d); if err != nil { return }\n",
			q.GoVar, q.Pos)
	default:
		v += fmt.Sprintf("%s := q%d", q.GoVar, q.Pos)
	}

	return
}

// Function contains information about a Go function and its associated routing
// information.
type Function struct {
	Definition  truce.Function
	Method      string
	Path        Path
	Query       map[string]QueryParam
	BodyVar     string
	BodyType    string
	HasReturn   bool
	ReturnType  string
	ReturnIsPtr bool
}

// Path is a collection of path segments (Elements) that represent a route's URL
// path.
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

// String implements fmt.Stringer
func (p Path) String() (v string) {
	for _, e := range p {
		v += "/" + e.String()
	}
	return
}

// FmtString returns the formatting strings of each Element as a joined path to
// be used as a sprintf formatting string.
func (p Path) FmtString() (v string) {
	for _, e := range p {
		v += "/" + e.FmtString()
	}
	return
}

// ArgString joins the Path's variable set into a list of arguments.
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

// Element is a Path segment. It can either be static or variable.
type Element struct {
	Type  string
	Value string
	Var   string
}

// String implements fmt.Stringer
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

// FmtString returns a sprintf formatting string for a Path segment
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

// NewFunction creates a Go-specific function definition from the abstract Truce definitions.
func NewFunction(config *truce.HTTP, function truce.Function) (*Function, error) {
	if function.Transports.HTTP == nil {
		return nil, nil
	}

	transport := *function.Transports.HTTP

	b := &Function{Definition: function, Query: map[string]QueryParam{}}

	type argument struct {
		variable string
		typ      string
	}

	var (
		pathMappings = map[string]string{}
		args         = map[string]argument{}
	)

	for _, field := range function.Arguments {
		args[field.Name] = argument{
			typ:      string(field.Type),
			variable: field.Name,
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

	var qpos int
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

			b.Query[arg.Var] = QueryParam{
				Pos:      qpos,
				QueryVar: arg.Var,
				GoVar:    a.variable,
				Type:     a.typ,
			}

			qpos++
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
