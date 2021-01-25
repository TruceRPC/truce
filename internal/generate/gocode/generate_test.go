package gocode

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

	var spec truce.Specification
	err = val.Decode(&spec)
	assert.NilError(t, err)

	versions, ok := spec.Specifications["example"]
	assert.Assert(t, ok)
	api, ok := versions["1"]
	assert.Assert(t, ok)

	t.Run("types.go", func(t *testing.T) {
		g, err := New(api)
		assert.NilError(t, err)

		actualData := bytes.NewBuffer(nil)
		err = g.WriteTypesTo(actualData, WithPackageName("example"))
		assert.NilError(t, err)

		expectedData, err := fs.ReadFile(testdata, "types.go.golden")
		assert.NilError(t, err)
		assert.DeepEqual(t, actualData.Bytes(), expectedData)
	})

	t.Run("server.go", func(t *testing.T) {
		g, err := New(api)
		assert.NilError(t, err)

		actualData := bytes.NewBuffer(nil)
		err = g.WriteServerTo(actualData,
			WithPackageName("example"),
			WithServerName("ExampleServer"))
		assert.NilError(t, err)

		expectedData, err := fs.ReadFile(testdata, "server.go.golden")
		assert.NilError(t, err)
		assert.DeepEqual(t, actualData.Bytes(), expectedData)
	})

	t.Run("client.go", func(t *testing.T) {
		g, err := New(api)
		assert.NilError(t, err)

		actualData := bytes.NewBuffer(nil)
		err = g.WriteClientTo(actualData,
			WithPackageName("example"),
			WithClientName("ExampleClient"))
		assert.NilError(t, err)

		expectedData, err := fs.ReadFile(testdata, "client.go.golden")
		assert.NilError(t, err)
		assert.DeepEqual(t, actualData.Bytes(), expectedData)
	})
}
