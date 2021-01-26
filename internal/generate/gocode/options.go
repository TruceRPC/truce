package gocode

// Option is a functional option for the Generator type.
type Option func(*Generator)

// Options is a list of Generator Options.
type Options []Option

// Apply applies the Options to a Generator.
func (o Options) Apply(g *Generator) {
	for _, o := range o {
		o(g)
	}
}

// WithServerName specifies the Go type name for the generated server.
func WithServerName(name string) Option {
	return func(g *Generator) {
		g.data.ServerName = name
	}
}

// WithClientName specifies the Go type name for the generated client.
func WithClientName(name string) Option {
	return func(g *Generator) {
		g.data.ClientName = name
	}
}

// WithPackageName specifies the package name of the generated files.
func WithPackageName(name string) Option {
	return func(g *Generator) {
		g.data.Package = name
	}
}
