package outputs

import (
	"github.com/TruceRPC/truce"
	"github.com/TruceRPC/truce/internal/outputs/gooutput"
	"github.com/TruceRPC/truce/internal/outputs/openapi"
)

func Write(def truce.Definition) error {
	if goOutput := def.Outputs.Go; goOutput != nil {
		if err := gooutput.Write(def.Spec, goOutput); err != nil {
			return err
		}
	}

	if oaOutput := def.Outputs.OpenAPI; oaOutput != nil {
		if err := openapi.Write(oaOutput); err != nil {
			return err
		}
	}

	return nil
}
