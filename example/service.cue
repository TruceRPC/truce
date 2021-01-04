import "strings"

#resources: {
	"User": {
		fields: name: type: "string"
	}
	"Post": {
		fields: {
			title: type: "string"
			body: type:  "string"
		}
	}
}

apis: [
	{
		version: "1"
		transports: {
			http: {
				versions: ["1.0", "1.1", "2.0"]
				prefix: "/api/v\(version)"
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
						path:   "\(strings.ToLower(resourceName))s/{id}"
						method: "GET"
						arguments: id: value: "$path.id"
					}
				}
				"Get\(resourceName)s": {
					return: {
						name: "\(strings.ToLower(resourceName))s"
						type: "[]\(resourceName)"
					}
					transports: http: {
						path:   "\(strings.ToLower(resourceName))s"
						method: "GET"
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
						arguments: "\(strings.ToLower(resourceName))": value: "$body"
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
							id: value:                                 "$path.id"
							"\(strings.ToLower(resourceName))": value: "$body"
						}
					}
				}
			}
		}
		types: {
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
	},
]
