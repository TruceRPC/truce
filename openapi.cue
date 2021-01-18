package truce

import "strings"

_#schemaObj: [_=string]: {
	_type:    string
	_isPrim:  _type =~ "^([]][]])?[a-z].*$"
	_isArray: _type =~ "^[[][]].*$"
	_trimmed: strings.TrimPrefix("\(_type)", "[]")
	if _isArray {
		type: "array"
		items: {
			if _isPrim {
				"type": _trimmed
			}
			if !_isPrim {
				"$ref": "#/components/schemas/\(_trimmed)"
			}
		}
	}
	if !_isArray {
		if _isPrim {
			"type": _trimmed
		}
		if !_isPrim {
			"$ref": "#/components/schemas/\(_trimmed)"
		}
	}
}

openapi3: {
	for apiName, apiVersions in specifications {
		for apiVersion, apiDef in apiVersions {
			"\(apiName)": {
				openapi: "3.0.3"
				info: {
					title:   apiName
					version: apiVersion
				}
				components: {
					schemas: {
						for typeName, typeDef in apiDef.types {
							"\(typeName)": {
								type: "object"
								properties: {
									for fieldName, fieldDef in typeDef.fields {
										"\(fieldName)": {
											type: fieldDef.type
										}
									}
								}
							}
						}
					}
				}
				paths: {
					for fnName, fnDef in apiDef.functions {
						let httpFn = fnDef.transports.http
						"\(httpFn.path)": {
							"\(strings.ToLower(httpFn.method))": {
								operationId: fnName
								// request parameters
								_parameters: [
									for argDef in fnDef.arguments {
										let fnArg = fnDef.transports.http.arguments[argDef.name]
										let schemaObj = _#schemaObj & {schema: {_type: argDef.type}}
										if fnArg != null && fnArg.from != "" && fnArg.from != "body" {
											{
												name:        argDef.name
												in:          fnArg.from
												description: "\(argDef.name) from \(fnArg.from)"
												required:    true
												schema:      schemaObj.schema
											}
										}
									},
								]
								_nonEmptyParameters: [ for p in _parameters if len(p) > 0 {p}]
								if len(_nonEmptyParameters) > 0 {
									parameters: _nonEmptyParameters
								}

								// request body
								_body: {
									for argDef in fnDef.arguments {
										let fnArg = fnDef.transports.http.arguments[argDef.name]
										if fnArg != null && fnArg.from == "body" {
											_#schemaObj & {schema: {_type: argDef.type}}
										}
									}
								}
								if len(_body) > 0 {
									requestBody: {
										description: "\(fnDef.name) operation request body"
										content: {
											"application/json": _body
										}
									}
								}

								// responses
								responses: {
									"200": {
										description: "\(fnDef.name) operation 200 response"
										content: {
											"application/json": {
												_#schemaObj & {schema: {_type: fnDef.return.type}}
											}
										}
									}
									default: {
										description: "\(fnDef.name) operation default response"
										content: "application/json": {schema: type: "string"}
									}
								}
							}
						}
					}
				}
			}
		}
	}
}
