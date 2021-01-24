package main

import (
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"os"

	"github.com/georgemac/truce"
	"github.com/georgemac/truce/internal/generate/gocode"
	"github.com/georgemac/truce/internal/generate/openapi"
)

func main() {
	if len(os.Args) < 3 {
		usage := fmt.Sprintf("Usage: %s", os.Args[0])
		fmt.Printf("%s <command>\n", usage)
		fmt.Printf("%s val[idate] <specification>\n", pad(len(usage)))
		fmt.Printf("%s gen[erate] <specification>\n", pad(len(usage)))
		os.Exit(2)
	}

	targetRaw, err := ioutil.ReadFile(os.Args[2])
	if err != nil {
		panic(err)
	}

	val, err := truce.Compile(targetRaw)
	if err != nil {
		panic(err)
	}

	var spec truce.Specification

	if err = val.Decode(&spec); err != nil {
		panic(err)
	}

	switch os.Args[1] {
	case "validate", "val":
		break
	case "generate", "gen":
		for n, versions := range spec.Outputs {
			for v, output := range versions {
				versions, ok := spec.Specifications[n]
				if !ok {
					fmt.Printf("\"%s v%s\" specification not found\n", n, v)
					os.Exit(1)
				}

				api, ok := versions[v]
				if !ok {
					fmt.Printf("\"%s v%s\" specification not found\n", n, v)
					os.Exit(1)
				}

				g, err := gocode.New(api)
				if err != nil {
					panic(err)
				}

				if goOut := output.Go; goOut != nil {
					fmt.Printf("Generating Service \"%s v%s\" Go Types\n", n, v)

					if target := goOut.Types; target != nil {
						writeTo(target.Path, g.WriteTypesTo,
							gocode.WithPackageName(target.Pkg))
					}

					if target := goOut.Server; target != nil {
						writeTo(target.Path, g.WriteServerTo,
							gocode.WithPackageName(target.Pkg),
							gocode.WithServerName(target.Type))
					}

					if target := goOut.Client; target != nil {
						writeTo(target.Path, g.WriteClientTo,
							gocode.WithPackageName(target.Pkg),
							gocode.WithClientName(target.Type))
					}
				}

				if oaOut := output.OpenAPI; oaOut != nil {
					fmt.Printf("Generating Service \"%s v%s\" Open API Specification\n", n, v)

					writeTo(oaOut.Path, func(w io.Writer, _ ...gocode.Option) error {
						if err := openapi.WriteJSON(w, val, n, v); err != nil {
							return err
						}

						return nil
					})

				}
			}
		}
	default:
		fmt.Printf("unexpected sub-command: %q\n", flag.Arg(0))
		os.Exit(2)
	}
}

func writeTo(path string, w func(io.Writer, ...gocode.Option) error, opts ...gocode.Option) {
	var f io.WriteCloser = nopWriteCloser{os.Stdout}
	if path != "<stdout>" {
		var err error
		f, err = os.OpenFile(path, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0644)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	}

	defer f.Close()

	if err := w(f, opts...); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

type nopWriteCloser struct {
	io.Writer
}

func (n nopWriteCloser) Close() error { return nil }

func pad(n int) (v string) {
	for i := 0; i < n; i++ {
		v += " "
	}
	return
}
