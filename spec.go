package truce

import (
	"cuelang.org/go/cue"
)

var Runtime = &cue.Runtime{}

type Specification struct {
	Outputs        map[string]map[string]Output `cue:"outputs"`
	Specifications map[string]map[string]API    `cue:"specifications"`
}

func Compile(data []byte) (cue.Value, error) {
	target, err := Runtime.Compile("target", data)
	if err != nil {
		return cue.Value{}, err
	}

	filled, err := cuegenInstance.Fill(target.Value())
	if err != nil {
		return cue.Value{}, err
	}

	val := filled.Value()

	return val, val.Err()
}

type Output struct {
	Name    string    `cue:"name"`
	Version string    `cue:"version"`
	Go      *GoOutput `cue:"go"`
	OpenAPI *OpenAPI  `cue:"openapi"`
}

type GoOutput struct {
	Types  *Target      `cue:"types"`
	Server *TypedTarget `cue:"server"`
	Client *TypedTarget `cue:"client"`
}

type OpenAPI struct {
	Path string `cue:"path"`
}

type Target struct {
	Path string `cue:"path"`
	Pkg  string `cue:"pkg"`
}

type TypedTarget struct {
	Target
	Type string `cue:"type"`
}

type API struct {
	Name       string              `cue:"name"`
	Version    string              `cue:"version"`
	Transports Transport           `cue:"transports"`
	Functions  map[string]Function `cue:"functions"`
	Types      map[string]Type     `cue:"types"`
}

// Transport to be defined at a later date.
type Transport struct {
	HTTP *HTTP `cue:"http"`
}

type HTTP struct {
	Versions []string              `cue:"versions"`
	Prefix   string                `cue:"prefix"`
	Errors   map[string]*HTTPError `cue:"errors"`
}

type HTTPError struct {
	Type       string `cue:"type"`
	StatusCode string `cue:"statusCode"`
}

type Function struct {
	Name       string            `cue:"name"`
	Arguments  []Field           `cue:"arguments"`
	Return     OptionalField     `cue:"return"`
	Transports FunctionTransport `cue:"transports"`
}

type OptionalField struct {
	Field
	Present bool `cur:"present"`
}

// FunctionTransport to be defined at a later date.
type FunctionTransport struct {
	HTTP *HTTPFunction `cue:"http"`
}

type HTTPFunction struct {
	Path      string                   `cue:"path"`
	Method    string                   `cue:"method"`
	Arguments map[string]ArgumentValue `cue:"arguments"`
}

type ArgumentValue struct {
	Name  string `cue:"name"`
	From  string `cue:"from"`
	Var   string `cue:"var"`
	Value string `cue:"value"`
}

type Type struct {
	Name   string           `cue:"name"`
	Type   string           `cue:"type"`
	Fields map[string]Field `cue:"fields"`
}

type Field struct {
	Name string `cue:"name"`
	Type string `cue:"type"`
}
