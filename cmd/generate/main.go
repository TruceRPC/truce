package main

import (
	"os"

	"github.com/georgemac/truce"
	"github.com/georgemac/truce/pkg/generate"
	"gopkg.in/yaml.v2"
)

func main() {
	var spec truce.Specification

	if err := yaml.Unmarshal([]byte(example), &spec); err != nil {
		panic(err)
	}

	cmd := "types"
	if len(os.Args) > 1 {
		cmd = os.Args[1]
	}

	switch cmd {
	case "types":
		if err := generate.GenerateTypes(os.Stdout, spec.Versions["0"]); err != nil {
			panic(err)
		}
	case "client":
		if err := generate.GenerateClient(os.Stdout, spec.Versions["0"]); err != nil {
			panic(err)
		}
	case "server":
		if err := generate.GenerateServer(os.Stdout, spec.Versions["0"]); err != nil {
			panic(err)
		}
	}
}

var example = `version:
  0:
    transports:
      - type: http
        versions: ["1.0", "1.1", "2"]
        prefix: "/api/v{{$version}}"
    functions:
      - name: GetResources
        returns:
          - name: resources
            type: "[]Resource"
        transports:
          - type: http
            method: "GET"
            path: "/resources"
      - name: GetResource
        arguments:
          - name: id
            type: string
        return:
          name: resource
          type: Resource
        transports:
          - type: http
            method: "GET"
            path: "/resources/{id}"
            arguments:
              - name: id
                value: "$path.id"
      - name: PutResource
        arguments:
          - name: group_id
            type: string
          - name: resource
            type: PutResourceRequest
        return:
          name: resource
          type: Resource
        transports:
          - type: http
            path: "/group/{group_id}/resources"
            method: "PUT"
            arguments:
              - name: group_id
                value: "$path.group_id"
              - name: resource
                value: "$body"
    types:
      - name: Resource
        fields:
          - name: id
            type: string
          - name: group_id
            type: string
          - name: name
            type: string
      - name: PutResourceRequest
        fields:
          - name: name
            type: string`
