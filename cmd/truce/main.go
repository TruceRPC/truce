package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/TruceRPC/truce"
	"github.com/TruceRPC/truce/internal/outputs"
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
	exitOnError(err)

	var truce truce.Truce
	exitOnError(truce.UnmarshalCUE(targetRaw))

	switch os.Args[1] {
	case "validate", "val":
		break
	case "generate", "gen":
		for _, versions := range truce.Truce {
			for _, definition := range versions {
				exitOnError(outputs.Write(definition))
			}
		}
	default:
		fmt.Printf("unexpected sub-command: %q\n", flag.Arg(0))
		os.Exit(2)
	}
}

func exitOnError(err error) {
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func pad(n int) (v string) {
	for i := 0; i < n; i++ {
		v += " "
	}
	return
}
