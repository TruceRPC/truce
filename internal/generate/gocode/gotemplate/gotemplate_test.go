package gotemplate

import (
	"testing"

	"gotest.tools/v3/assert"
)

func TestParseAll(t *testing.T) {
	tmpl, err := ParseAll()
	assert.NilError(t, err)

	known := []string{
		"types.go.tmpl",
		"server.go.tmpl",
		"client.go.tmpl",
	}
	for _, n := range known {
		assert.Assert(t, tmpl.Lookup(n) != nil)
	}

	assert.Assert(t, tmpl.Lookup("totally_unknown.go.tmpl") == nil)
}
