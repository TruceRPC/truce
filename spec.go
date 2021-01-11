package truce

import "cuelang.org/go/cue"

var Runtime = &cue.Runtime{}

type Specification struct {
	Outputs        map[string]map[string]Output `cue:"outputs"`
	Specifications map[string]map[string]API    `cue:"specifications"`
}

func Unmarshal(data []byte, spec *Specification) error {
	target, err := Runtime.Compile("target", data)
	if err != nil {
		return err
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

type Output struct {
	Name    string      `cue:"name"`
	Version string      `cue:"version"`
	HTTP    *HTTPOutput `cue:"http"`
}

type HTTPOutput struct {
	Types  *Target      `cue:"types"`
	Server *TypedTarget `cue:"server"`
	Client *TypedTarget `cue:"client"`
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
	Versions []string           `cue:"versions"`
	Prefix   string             `cue:"prefix"`
	Errors   map[int]*HTTPError `cue:"errors"`
}

type HTTPError struct {
	Type       string `cue:"type"`
	StatusCode string `cue:"statusCode"`
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
	From  string `cue:"from"`
	Var   string `cue:"var"`
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
	Type   string           `cue:"type"`
	Fields map[string]Field `cue:"fields"`
}

type Field struct {
	Name string   `cue:"name"`
	Type TypeName `cue:"type"`
}
