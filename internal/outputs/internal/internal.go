package internal

import (
	"io"
	"os"
)

func WriterFor(path string) (io.WriteCloser, error) {
	var f io.WriteCloser = nopWriteCloser{os.Stdout}
	if path != "<stdout>" {
		var err error
		f, err = os.OpenFile(path, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0644)
		if err != nil {
			return nil, err
		}
	}

	return f, nil
}

type nopWriteCloser struct {
	io.Writer
}

func (n nopWriteCloser) Close() error { return nil }
