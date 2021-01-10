package truce

outputs: [n=string]: [v=string]: {
		name:       string
		version:    =~"^[0-9]+$"
		http?:      #HTTPOutput
} & {name: n, version: v}

#HTTPOutput: {
	types?:  #Target
	server?: #TypedTarget
	client?: #TypedTarget
}

#TypedTarget: {
	#Target
	type: string | *"Server"
}

#Target: {
	path: string | *"<stdout>"
	pkg:  string | *"types"
}

specifications: [n=string]: [v=string]: #API & {name: n, version: v}

#API: {
	name:    string
	version: =~"^[0-9]+$" // API Major version
	transports?: http?: #HTTP
	functions: [n=_]:   #Function & {"name": n}
	types: [n=_]:       #Type & {"name":     n}
}

#HTTP: {
	versions: [...string]
	prefix: string
}

#Function: {
	name: =~"^[A-Z][a-zA-Z]*$" // function names are capitalised
	arguments: [...#Field]
	return: #Field
	transports?: http?: #HTTPFunction
}

#Type: {
	name: =~"^[A-Z][a-zA-Z]*$" // type names are capitalised
	fields: [n=_]: #Field & {"name": n}
}

#primitives: ["bool", "int", "float64", "byte", "string"]
#slices: [ for x in #primitives {"[]\(x)"}]
#all: #primitives + #slices

#Field: {
	name: string
	type: or(#all) | =~"map[[][*]?[A-Za-z]+[]][*]?[A-Za-z]+" | =~"^[*]?[A-Z][a-zA-Z]*$" | =~"^[[][]][*]?[A-Z][a-zA-Z]*?" // can be primitive, Custom, *Custom, []primitive, []Custom, []*Custom.
}

#HTTPFunction: {
	path:   string
	method: "GET" | "POST" | "PUT" | "PATCH" | "OPTIONS" | "DELETE" | "HEAD" | "CONNECT" | "TRACE"
	arguments: [n=_]: #Argument & {"name": n}
}

#Argument: {
	name: string
	from: "body"
} | {
	name: string
	from: "path"
	var:  string
} | {
	name: string
	from: "query"
	var:  string
} | {
	name:  string
	value: string
}
