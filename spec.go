package truce

import "fmt"

type Specification struct {
	Versions map[string]API `yaml:"version"`
}

type API struct {
	Transports []Transport `yaml:"transports"`
	Functions  []Function  `yaml:"functions"`
	Types      []Type      `yaml:"types"`
}

// Transport to be defined at a later date.
type Transport struct {
	Type   string      `yaml:"type"`
	Config interface{} `yaml:"-"`
}

func (t *Transport) UnmarshalYAML(unmarshal func(interface{}) error) error {
	transport := struct {
		Type string `yaml:"type"`
	}{}

	if err := unmarshal(&transport); err != nil {
		return err
	}

	t.Type = transport.Type

	switch transport.Type {
	case "http":
		var http HTTP
		if err := unmarshal(&http); err != nil {
			return err
		}

		t.Config = http
	default:
		return fmt.Errorf("%q transport not supported", transport.Type)
	}

	return nil
}

type Function struct {
	Name       string              `yaml:"name"`
	Arguments  []Field             `yaml:"arguments"`
	Return     Field               `yaml:"return"`
	Transports []FunctionTransport `yaml:"transports"`
}

// FunctionTransport to be defined at a later date.
type FunctionTransport struct {
	Type   string      `yaml:"type"`
	Config interface{} `yaml:"-"`
}

func (t *FunctionTransport) UnmarshalYAML(unmarshal func(interface{}) error) error {
	transport := struct {
		Type string `yaml:"type"`
	}{}

	if err := unmarshal(&transport); err != nil {
		return err
	}

	t.Type = transport.Type

	switch transport.Type {
	case "http":
		var http HTTPFunction
		if err := unmarshal(&http); err != nil {
			return err
		}

		t.Config = http
	default:
		return fmt.Errorf("%q transport not supported for function", transport.Type)
	}

	return nil
}

type ArgumentValue struct {
	Name  string `yaml:"name"`
	Value string `yaml:"value"`
}

type TypeName string

const (
	StringTypeName = TypeName("string")
	FloatTypeName  = TypeName("float")
	IntTypeName    = TypeName("int")
)

type Type struct {
	Name   string  `yaml:"name"`
	Fields []Field `yaml:"fields"`
}

type Field struct {
	Name string   `yaml:"name"`
	Type TypeName `yaml:"type"`
}
