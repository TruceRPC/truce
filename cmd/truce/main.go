package main

import (
	"flag"
	"io"
	"io/ioutil"
	"os"

	"cuelang.org/go/cue"
	"cuelang.org/go/encoding/yaml"
	"github.com/georgemac/truce"
	"github.com/georgemac/truce/pkg/generate"
	gyaml "gopkg.in/yaml.v2"
)

var source = flag.String("src", "", "Filepath for source truce specification")

func main() {
	flag.Parse()

	truceRaw, err := ioutil.ReadFile("truce.cue")
	if err != nil {
		panic(err)
	}

	var r cue.Runtime

	core, err := r.Compile("truce", truceRaw)
	if err != nil {
		panic(err)
	}

	targetRaw, err := ioutil.ReadFile(*source)
	if err != nil {
		panic(err)
	}

	switch flag.Arg(0) {
	case "validate", "val":
		if err := yaml.Validate(targetRaw, core.Value()); err != nil {
			panic(err)
		}
	case "generate", "gen":
		if err := yaml.Validate(targetRaw, core.Value()); err != nil {
			panic(err)
		}

		var spec truce.Specification

		if err := gyaml.Unmarshal(targetRaw, &spec); err != nil {
			panic(err)
		}

		var generator func(io.Writer, truce.API) error

		switch flag.Arg(1) {
		case "types":
			generator = generate.GenerateTypes
		case "client":
			generator = generate.GenerateClient
		case "server":
			generator = generate.GenerateServer
		case "", "all":
			generator = func(w io.Writer, a truce.API) error {
				if err := generate.GenerateTypes(w, a); err != nil {
					panic(err)
				}
				if err := generate.GenerateClient(w, a); err != nil {
					panic(err)
				}
				return generate.GenerateServer(w, a)
			}
		}

		if err := generator(os.Stdout, spec.APIs[0]); err != nil {
			panic(err)
		}
	}
}
