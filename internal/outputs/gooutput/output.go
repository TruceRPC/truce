package gooutput

import (
	"io"

	"github.com/TruceRPC/truce"
	"github.com/TruceRPC/truce/internal/generate/gocode"
	"github.com/TruceRPC/truce/internal/outputs/internal"
)

// Write generates and writes out the Go definitions for the provided
// api into the destinations described in goOutput.
func Write(api truce.API, goOutput *truce.GoOutput) error {
	g, err := gocode.New(api)
	if err != nil {
		return err
	}

	if target := goOutput.Types; target != nil {
		if err := writeTo(target.Path, g.WriteTypesTo,
			gocode.WithPackageName(target.Pkg)); err != nil {
			return err
		}
	}

	if target := goOutput.Server; target != nil {
		if err := writeTo(target.Path, g.WriteServerTo,
			gocode.WithPackageName(target.Pkg),
			gocode.WithServerName(target.Type)); err != nil {
			return err
		}
	}

	if target := goOutput.Client; target != nil {
		if err := writeTo(target.Path, g.WriteClientTo,
			gocode.WithPackageName(target.Pkg),
			gocode.WithClientName(target.Type)); err != nil {
			return err
		}
	}

	return nil
}

func writeTo(path string, w func(io.Writer, ...gocode.Option) error, opts ...gocode.Option) error {
	wr, err := internal.WriterFor(path)
	if err != nil {
		return err
	}

	defer wr.Close()

	return w(wr, opts...)
}
