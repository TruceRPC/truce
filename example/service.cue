import "strings"

import "list"

#resources: [
	{
		name: "User"
		fields: [
			{name: "name", type: "string"},
		]
	},
	{
		name: "Post"
		fields: [
			{name: "title", type: "string"},
			{name: "body", type:  "string"},
		]
	},
]

apis: [
	{
		version: "1"
		transports: {
			http: {
				versions: ["1.0", "1.1", "2.0"]
				prefix: "/api/v\(version)"
			}
		}
		functions: list.FlattenN([ for resource in #resources {
			[ {
				name: "Get\(resource.name)"
				arguments: [
					{
						name: "id"
						type: "string"
					},
				]
				return: {
					name: strings.ToLower(resource.name)
					type: resource.name
				}
				transports: http: {
					path:   strings.Join(["", strings.ToLower(resource.name) + "s", "{id}"], "/")
					method: "GET"
					arguments: [{
						name:  "id"
						value: "$path.id"
					}]
				}
			},
				{
					name: "Get\(resource.name)s"
					return: {
						name: "\(strings.ToLower(resource.name))s"
						type: "[]\(resource.name)"
					}
					transports: http: {
						path:   strings.Join(["", strings.ToLower(resource.name) + "s"], "/")
						method: "GET"
					}
				},
				{
					name: "Put\(resource.name)"
					arguments: [
						{
							name: strings.ToLower(resource.name)
							type: "Put\(resource.name)Request"
						},
					]
					return: {
						name: strings.ToLower(resource.name)
						type: resource.name
					}
					transports: http: {
						path:   strings.Join(["", strings.ToLower(resource.name) + "s"], "/")
						method: "PUT"
						arguments: [{
							name:  strings.ToLower(resource.name)
							value: "$body"
						}]
					}
				},
				{
					name: "Patch\(resource.name)"
					arguments: [
						{
							name: "id"
							type: "string"
						},
						{
							name: strings.ToLower(resource.name)
							type: "Patch\(resource.name)Request"
						},
					]
					return: {
						name: strings.ToLower(resource.name)
						type: resource.name
					}
					transports: http: {
						path:   strings.Join(["", strings.ToLower(resource.name) + "s", "{id}"], "/")
						method: "PATCH"
						arguments: [
							{
								name:  "id"
								value: "$path.id"
							},
							{
								name:  strings.ToLower(resource.name)
								value: "$body"
							}]
					}
				},
			]},
		], 1)
		types: list.FlattenN([ for resource in #resources {
			[
				{
					name:   resource.name
					fields: [
						{name: "id", type: "string"},
					] + resource.fields
				},
				{
					name:   "Put\(resource.name)Request"
					fields: resource.fields
				},
				{
					name:   "Patch\(resource.name)Request"
					fields: resource.fields
				},
			]},
		], 1)
	},
]
