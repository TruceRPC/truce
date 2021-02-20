package truce

import (
	_ "embed"

	"cuelang.org/go/cue"
	"cuelang.org/go/cue/build"
	_ "cuelang.org/go/pkg"
)

var (
	//go:embed cue/truce.cue
	truceSource []byte
	//go:embed cue/openapi.cue
	openAPISource []byte

	Runtime = &cue.Runtime{}

	Instance *cue.Instance
)

func must(err error) {
	if err != nil {
		panic(err)
	}
}

func init() {
	instance := build.NewContext().NewInstance("cue", nil)
	must(instance.AddFile("truce.cue", truceSource))
	must(instance.AddFile("openapi.cue", openAPISource))

	var err error
	Instance, err = Runtime.Build(instance)
	if err != nil {
		panic(err)
	}
}
