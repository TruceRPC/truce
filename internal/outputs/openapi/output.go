package openapi

import (
	"github.com/TruceRPC/truce"
	"github.com/TruceRPC/truce/internal/outputs/internal"
)

func Write(oaOut *truce.OpenAPIOutput) error {
	wr, err := internal.WriterFor(oaOut.Path)
	if err != nil {
		return err
	}

	defer wr.Close()

	return oaOut.MarshalJSON(wr)
}
