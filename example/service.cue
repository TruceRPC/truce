package truce

import "strings"

#resources: {
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

outputs:
	"example": "1": {
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

specifications: {
	"example": "1": {
		transports: {
			http: {
				versions: ["1.0", "1.1", "2.0"]
				prefix: "/api/v1"
				errors: {
					"401": type: "NotAuthorized"
					"404": type: "NotFound"
				}
			}
		}
		functions: {
			for resourceName, resource in #resources {
				"Get\(resourceName)": {
					arguments: [
						{name: "id", type: "string"},
					]
					return: {
						name: "\(strings.ToLower(resourceName))"
						type: resourceName
					}
					transports: http: {
						path:                 "/\(strings.ToLower(resourceName))s/{id}"
						method:               "GET"
						arguments: id: {from: "path", var: "id"}
					}
				}
				"Get\(resourceName)s": {
					arguments: [
						{name: "limit", type: "int"},
					]
					return: {
						name: "\(strings.ToLower(resourceName))s"
						type: "[]\(resourceName)"
					}
					transports: http: {
						path:                    "/\(strings.ToLower(resourceName))s"
						method:                  "GET"
						arguments: limit: {from: "query", var: "limit"}
					}
				}
				"Put\(resourceName)": {
					arguments: [
						{name: "\(strings.ToLower(resourceName))", type: "Put\(resourceName)Request"},
					]
					return: {
						name: "\(strings.ToLower(resourceName))"
						type: resourceName
					}
					transports: http: {
						path:   "/\(strings.ToLower(resourceName))s"
						method: "PUT"
						arguments: "\(strings.ToLower(resourceName))": from: "body"
					}
				}
				"Patch\(resourceName)": {
					arguments: [
						{name: "id", type:                               "string"},
						{name: "\(strings.ToLower(resourceName))", type: "Patch\(resourceName)Request"},
					]
					return: {
						name: "\(strings.ToLower(resourceName))"
						type: resourceName
					}
					transports: http: {
						path:   "/\(strings.ToLower(resourceName))s/{id}"
						method: "PATCH"
						arguments: {
							id: {from: "path", var: "id"}
							"\(strings.ToLower(resourceName))": from: "body"
						}
					}
				}
			}
		}
		types: {
			NotAuthorized: {
				type: "error"
				fields: message: type: "string"
			}
			NotFound: {
				type: "error"
				fields: message: type: "string"
			}
			for resourceName, resource in #resources {
				"\(resourceName)": {
					fields: {
						id: type: "string"
						for k, v in resource.fields {
							"\(k)": {
								for k1, v1 in v {"\(k1)": v1}
							}
						}
					}
				}
				"Put\(resourceName)Request": {
					fields: {
						for k, v in resource.fields {
							"\(k)": {
								for k1, v1 in v {"\(k1)": v1}
							}
						}
					}
				}
				"Patch\(resourceName)Request": {
					fields: {
						for k, v in resource.fields {
							"\(k)": {
								for k1, v1 in v {"\(k1)": v1}
							}
						}
					}
				}
			}
		}
	}
}
