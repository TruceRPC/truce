package main

import (
	"io/ioutil"
	"os"

	"github.com/georgemac/truce"
	"github.com/georgemac/truce/pkg/generate"
	"gopkg.in/yaml.v2"
)

func main() {
	var spec truce.Specification

	example, err := ioutil.ReadFile("./example/service.yml")
	if err != nil {
		panic(err)
	}

	if err := yaml.Unmarshal(example, &spec); err != nil {
		panic(err)
	}

	cmd := "types"
	if len(os.Args) > 1 {
		cmd = os.Args[1]
	}

	switch cmd {
	case "types":
		if err := generate.GenerateTypes(os.Stdout, spec.APIs[0]); err != nil {
			panic(err)
		}
	case "client":
		if err := generate.GenerateClient(os.Stdout, spec.APIs[0]); err != nil {
			panic(err)
		}
	case "server":
		if err := generate.GenerateServer(os.Stdout, spec.APIs[0]); err != nil {
			panic(err)
		}
	}
}
