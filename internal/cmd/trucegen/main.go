package main

import (
	"bytes"
	"io"
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

	// TODO(georgemac): put back ioutil.WriteFile once https://github.com/cuelang/cue/pull/664 is resolved.
	fi, err := os.OpenFile("./truce.go", os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0644)
	if err != nil {
		panic(err)
	}

	defer fi.Close()

	buf := bytes.NewBuffer(v)
	line, err := buf.ReadBytes('\n')
	for ; err == nil; line, err = buf.ReadBytes('\n') {
		if bytes.Contains(line, []byte("LookupField")) {
			_, _ = fi.Write([]byte("\t//lint:ignore SA1019 until FieldByName is produced by cue gocode https://github.com/cuelang/cue/pull/664\n"))
		}

		_, _ = fi.Write(line)
	}

	if err != io.EOF {
		panic(err)
	}

	// write the final line
	_, _ = fi.Write(line)
}
