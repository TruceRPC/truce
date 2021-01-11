Truce - An RPC Framework
------------------------

Truce is a combined API and RPC specification language. It is built in [Go](https://golang.org) and powered by [Cue](https://cuelang.org).

Its goal is to enable machine friendly RPC definitions to be expressed and exposed through human friendly API specifications (like REST). 

RPC frameworks like gRPC, Twirp and WebRPC are awesome. They allow for service and procedure definitons to be expressed as configuration, which is then automatically generated into server and client definitions. The developer is left to primarily focus on the business logic.
The resulting wire-types are often rigid and machine focussed. gRPC even requires protobuf and HTTP 2. Making it not particularly browser friendly, though plenty of projects exist to get this working. Each of these projects (to differing degrees) is pushing the developer to use code generation. Meaning, whichever language you're hoping to integrate with will require a well maintained set of code generation tools for that particular flavour of RPC.

API specification languages like OpenAPI are expressive. They allow for APIs to be defined which leverage HTTP in all its weird and wonderful ways. This ecosystem also offers code generators in order to simplify the scaffolding of these services. However, the breadth of capabilities made available by OpenAPI leads to lots of complexity in generators and generated code.

Truce is in search of a sweet spot. Here are some of the desired traits this project aims to incorporate:

- Be language agnostic (Starting with generators written in and for Go).
- Generate servers and clients.
- Make space for multiple target transport protocols and versions (e.g. HTTP 1.1 and 2, maybe even gRPC).
- Ability to express human readable wire-formats (REST with JSON).
- Versioned API definitions by default.

## Usage

> See [Building](#Building) to compile truce. This is currently not distributed by any means other than git.

Truce currently can generate Go struct definitions (types), clients (Go over http) and servers (Go over http) definitions.

```
bin/truce -src=<source cue definition> val[idate] # validate Truce source CUE definitions.
bin/truce -src=<source cue definition> gen[erate] # generate Truce types, client and server definitions based on source.
```

Try it out with the `example` directory:

```
bin/truce --src=example/service.cue gen
```

## Building

`make build` will output truce into a local `bin/` folder.

### Requirements

- Make
- Go
- Cue

## Example

As mentioned under [usage](#Usage) the [./example](./example) directory contains some examples of truce in action.

In particular try running go run [example/cmd/main.go](./example/cmd/main.go).

This runs a small web server which can be pinged on port 8080. This example server exercises the generated resources defined in the [service definition](./example/service.cue).

> The following demo requires `curl` and `jq`

```
# in one terminal
go run example/cmd/main.go

# in another terminal
# create first post
POST_ID=$(curl -X PUT --data-raw "{\"title\":\"Hello, World\",\"body\":\"Welcome to my blog\"}" localhost:8080/api/v1/posts | jq -r '.id')

# list posts
curl localhost:8080/api/v1/posts

# update post with "..."
curl -X PATCH --data-raw "{\"title\":\"Hello, World\",\"body\":\"Welcome to my blog...\"}" localhost:8080/api/v1/posts/$POST_ID

# list posts again
curl localhost:8080/api/v1/posts
```
