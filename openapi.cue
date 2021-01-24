package truce

import "strings"

import "regexp"

_#schemaObj: [_=string]: {
	_type: string

	let mapMatcher = "^map[[]([^]])*[]](.*)$"
	_isMap:   _type =~ mapMatcher
	_isArray: bool
	_isPrim:  bool
	_base:    string
	_format:  string

	let trimmed = strings.TrimPrefix("\(_type)", "[]")
	if trimmed == "int" {
		_base:   "integer"
		_format: "int64"
	}
	if trimmed == "float64" {
		_base:   "number"
		_format: "float"
	}
	if trimmed == "bool" {
		_base:   "boolean"
		_format: ""
	}
	if trimmed != "int" &&
		trimmed != "bool" &&
		trimmed != "float64" &&
		trimmed != "byte" {
		_base:   trimmed
		_format: ""
	}
	if _type == "[]byte" {
		// []byte is represented on the wire as base64([]byte("..."))
		_isArray: false
		_isPrim:  true
		_base:    "string"
		_format:  "base64"
	}
	if _type != "[]byte" {
		_isArray: _type =~ "^[[][]].*$"
		_isPrim:  _type =~ "^([]][]])?[a-z].*$"
	}

	if _isMap {
		_mapParts:  regexp.FindAllSubmatch(mapMatcher, _type, 3)
		_isValPrim: _mapParts[2] =~ "^([]][]])?[a-z].*$"

		type: "object"
		additionalProperties: {
			if _isValPrim {
				type: _mapParts[2]
			}
			if !_isValPrim {
				"$ref": "#/components/schemas/\(_mapParts[2])"
			}
		}
	}
	if !_isMap {
		if _isArray {
			type: "array"
			items: {
				if _isPrim {
					"type": _base
				}
				if !_isPrim {
					"$ref": "#/components/schemas/\(_base)"
				}
			}
		}
		if !_isArray {
			if _isPrim {
				"type": _base
				if _format != "" {
					"format": _format
				}
			}
			if !_isPrim {
				"$ref": "#/components/schemas/\(_base)"
			}
		}
	}
}

openapi3: {
	for apiName, apiVersions in specifications {
		for apiVersion, apiDef in apiVersions {
			"\(apiName)": "\(apiVersion)": {
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
										let schema = _#schemaObj & {schema: {_type: fieldDef.type}}
										"\(fieldName)": schema.schema
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
								}
							}
						}
					}
				}
			}
		}
	}
}
