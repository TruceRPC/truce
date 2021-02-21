package openapi

import (
	"bytes"
	"io/fs"
	"os"
	"testing"

	"github.com/TruceRPC/truce"
	"gotest.tools/v3/assert"
)

var testdata = os.DirFS("testdata")

func TestGenerator(t *testing.T) {
	data, err := fs.ReadFile(testdata, "service.cue")
	assert.NilError(t, err)

	var truce truce.Truce

	err = truce.UnmarshalCUE(data)
	assert.NilError(t, err)

	t.Run("swagger.json", func(t *testing.T) {
		actualData := bytes.NewBuffer(nil)
		err = truce.Truce["example"]["1"].Outputs.OpenAPI.MarshalJSON(actualData)
		assert.NilError(t, err)

		expectedData, err := fs.ReadFile(testdata, "swagger.json.golden")
		assert.NilError(t, err)
		assert.DeepEqual(t, actualData.Bytes(), expectedData)
	})
}
