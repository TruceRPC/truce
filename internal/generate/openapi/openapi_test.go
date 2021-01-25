package openapi

import (
	"bytes"
	"io/fs"
	"os"
	"testing"

	"cuelang.org/go/cue"
	"github.com/TruceRPC/truce"
	"gotest.tools/v3/assert"
)

var testdata = os.DirFS("testdata")

func TestGenerator(t *testing.T) {
	data, err := fs.ReadFile(testdata, "service.cue")
	assert.NilError(t, err)

	var val cue.Value
	val, err = truce.Compile(data)
	assert.NilError(t, err)

	t.Run("swagger.json", func(t *testing.T) {
		actualData := bytes.NewBuffer(nil)
		err = WriteJSON(actualData, val, "example", "1")
		assert.NilError(t, err)

		expectedData, err := fs.ReadFile(testdata, "swagger.json.golden")
		assert.NilError(t, err)
		assert.DeepEqual(t, actualData.Bytes(), expectedData)
	})
}
