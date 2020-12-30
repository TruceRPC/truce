package truce

type HTTP struct {
	Versions []string `yaml:"versions"`
	Prefix   string   `yaml:"prefix"`
}

type HTTPFunction struct {
	Path      string          `yaml:"path"`
	Method    string          `yaml:"method"`
	Arguments []ArgumentValue `yaml:"arguments"`
}
