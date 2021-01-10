package generate

type Option func(*Generator)

type Options []Option

func (o Options) Apply(g *Generator) {
	for _, o := range o {
		o(g)
	}
}

func WithServerName(name string) Option {
	return func(g *Generator) {
		g.context.ServerName = name
	}
}

func WithClientName(name string) Option {
	return func(g *Generator) {
		g.context.ClientName = name
	}
}

func WithPackageName(name string) Option {
	return func(g *Generator) {
		g.context.Package = name
	}
}
