package main

import (
	"io/ioutil"
	"os"

	"cuelang.org/go/cue"
	"cuelang.org/go/cue/load"
	"cuelang.org/go/encoding/gocode"
)

func main() {
	cwd, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	inst := cue.Build(load.Instances([]string{"."}, &load.Config{
		Dir:        cwd,
		ModuleRoot: cwd,
		Module:     "github.com/TruceRPC/truce",
	}))[0]
	if err := inst.Err; err != nil {
		panic(err)
	}

	v, err := gocode.Generate(".", inst, &gocode.Config{
		RuntimeVar: "Runtime",
	})
	if err != nil {
		panic(err)
	}

	if err := ioutil.WriteFile("./truce.go", v, 0644); err != nil {
		panic(err)
	}
}
