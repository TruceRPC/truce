package main

import (
	"github.com/davecgh/go-spew/spew"
	"github.com/georgemac/truce"
	"gopkg.in/yaml.v2"
)

func main() {
	var spec truce.Specification

	if err := yaml.Unmarshal([]byte(example), &spec); err != nil {
		panic(err)
	}

	spew.Dump(spec)
}

var example = `version:
  0:
    transports:
      - type: http
        versions: ["1.0", "1.1", "2"]
        mappings:
          - prefix: "/api/v{{$version}}"
            mappings:
              - path: "/buckets"
                call:
                  name: GetBuckets
              - path: "/buckets/:id"
                call:
                  name: GetBucket
                  arguments:
                    - name: id
                      value: "$path.id"
    functions:
      - name: GetBuckets
        returns:
          - name: buckets
            type: "[]Bucket"
      - name: GetBucket
        arguments:
          - name: id
            type: string
        returns:
          - name: bucket
            type: Bucket
    types:
      - name: Bucket
        fields:
          - name: id
            type: string
          - name: name
            type: string`
