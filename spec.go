package truce

import "cuelang.org/go/cue"

var Runtime = &cue.Runtime{}

type Specification struct {
	APIs []API `cue:"apis"`
}

func Unmarshal(data []byte, spec *Specification) error {
	target, err := Runtime.Compile("target", data)
	if err != nil {
		return nil
	}

	filled, err := cuegenInstance.Fill(target.Value())
	if err != nil {
		return err
	}

	val := filled.Value()
	if err := val.Err(); err != nil {
		return err
	}

	if err := val.Decode(spec); err != nil {
		return err
	}

	return nil
}

type API struct {
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
	Versions []string `cue:"versions"`
	Prefix   string   `cue:"prefix"`
}

type Function struct {
	Name       string            `cue:"name"`
	Arguments  []Field           `cue:"arguments"`
	Return     Field             `cue:"return"`
	Transports FunctionTransport `cue:"transports"`
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
	Value string `cue:"value"`
}

type TypeName string

const (
	StringTypeName = TypeName("string")
	FloatTypeName  = TypeName("float")
	IntTypeName    = TypeName("int")
)

type Type struct {
	Name   string           `cue:"name"`
	Fields map[string]Field `cue:"fields"`
}

type Field struct {
	Name string   `cue:"name"`
	Type TypeName `cue:"type"`
}
