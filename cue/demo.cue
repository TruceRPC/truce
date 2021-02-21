package truce

#crud: "example": "1": {
	outputs: {
		openapi: {
			version: 3
			path:    "example/swagger.json"
		}
		go: {
			types: {path: "example/types.go", pkg: "example"}
			server: {
				path: "example/server.go"
				type: "Server"
				pkg:  "example"
			}
			client: {
				path: "example/client.go"
				type: "Client"
				pkg:  "example"
			}
		}
	}
	resources: {
		"User": {
			fields: {
				name: type:   "string"
				age: type:    "int"
				height: type: "float64"
				labels: type: "object"
			}
		}
		"Post": {
			fields: {
				title: type:   "string"
				body: type:    "[]byte"
				draft: type:   "bool"
				created: type: "*timestamp"
			}
		}
	}
}
