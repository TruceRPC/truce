package main

import (
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"strings"

	"github.com/georgemac/truce"
	"github.com/georgemac/truce/pkg/generate"
)

var (
	source = flag.String("src", "", "Filepath for source truce specification")
	dst    = flag.String("w", "<stdout>", "Destination to write generated output (can be in form 'path' or 'type:path' or blank for defaults)")
)

func main() {
	flag.Parse()

	targetRaw, err := ioutil.ReadFile(*source)
	if err != nil {
		panic(err)
	}

	var spec truce.Specification

	if err = truce.Unmarshal(targetRaw, &spec); err != nil {
		panic(err)
	}

	switch flag.Arg(0) {
	case "validate", "val":
		break
	case "generate", "gen":
		targets := map[string]io.Writer{
			"types":  os.Stdout,
			"client": os.Stdout,
			"server": os.Stdout,
		}

		if *dst != "<stdout>" {
			if *dst == "" {
				// apply defaults
				*dst = "types:types_truce_gen.go,client:client_truce_gen.go,server:server_truce_gen.go"
			}

			files := map[string]*os.File{}
			for _, target := range strings.Split(*dst, ",") {
				v := strings.SplitN(target, ":", 2)
				var (
					target = flag.Arg(1)
					path   = v[0]
				)

				if len(v) > 1 {
					target = v[0]
					path = v[1]
				}

				if path == "<stdout>" {
					targets[target] = os.Stdout
					continue
				}

				// cache opened files
				f, ok := files[path]
				if !ok {
					f, err = os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
					if err != nil {
						panic(err)
					}

					defer f.Close()
					files[path] = f
				}

				ts := []string{target}
				if target == "" {
					ts = []string{"types", "client", "server"}
				}
				for _, target := range ts {
					targets[target] = f
				}
			}
		}

		var generator func(truce.API) error
		curry := func(dst io.Writer, fn func(io.Writer, truce.API) error) func(truce.API) error {
			return func(a truce.API) error {
				return fn(dst, a)
			}
		}

		switch flag.Arg(1) {
		case "types":
			generator = curry(targets["types"], generate.GenerateTypes)
		case "client":
			generator = curry(targets["client"], generate.GenerateClient)
		case "server":
			generator = curry(targets["server"], generate.GenerateServer)
		case "", "all":
			generator = func(a truce.API) error {
				if err := generate.GenerateTypes(targets["types"], a); err != nil {
					panic(err)
				}
				if err := generate.GenerateClient(targets["client"], a); err != nil {
					panic(err)
				}
				return generate.GenerateServer(targets["server"], a)
			}
		default:
			fmt.Printf("unexpected generation selector: %q\n", flag.Arg(1))
			os.Exit(2)
		}

		for _, api := range spec.Specifications {
			if err := generator(api); err != nil {
				panic(err)
			}
		}
	default:
		fmt.Printf("unexpected sub-command: %q\n", flag.Arg(0))
		os.Exit(2)
	}
}
