package openapi

import (
	"bytes"
	"encoding/json"
	"io"

	"cuelang.org/go/cue"
)

// WriteJSON marshals the named and version specification in
// JSON to the provided writer.
func WriteJSON(w io.Writer, val cue.Value, name, version string) error {
	val = val.Lookup("openapi3", name, version)
	if err := val.Err(); err != nil {
		return err
	}

	data, err := val.MarshalJSON()
	if err != nil {
		return err
	}

	buf := &bytes.Buffer{}
	if err := json.Indent(buf, data, "", "    "); err != nil {
		return err
	}

	_, err = buf.WriteTo(w)
	return err
}
