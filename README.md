Truce - RPC 🤝 People
----------------------

Truce is an API and RPC specification language.
Its goal is to enable machine friendly RPC definitions to be expressed and exposed through human friendly API specifications (like REST). 

RPC frameworks like gRPC, Twirp and WebRPC are awesome. They allow for service and procedure definitons to be expressed as configuration, which is then automatically generated into server and client definitions. The developer is left to primarily focus on the business logic.
The resulting wire types are often rigid and machine focussed. gRPC even requires protobuf and HTTP 2. Making it not particularly browser friendly, though plenty of projects exist to get this working. Each of these projects (to differing degrees) is pushing the developer to use code generation. Meaning, whichever language you're hoping to integrate with better have a well maintained set of code generation tools for that particular flavour of RPC.

API specification languages like OpenAPI are expressive. They allow for APIs to be defined which leverage HTTP in all its weird and wonderful ways. This ecosystem also offers code generators in order to simplify the scaffolding of these services. However, the breadth of capabilities made available by OpenAPI leads to complex generators and generated code.

Truce is in search of a sweet spot.

Here are some of the desirable traits:

- Machine generated servers and clients.
- Space for multiple transport protocols and versions (e.g. HTTP 1.1 and 2).
- Ability to express human readable wire-formats (REST with JSON).
- Versioned definitions by default.
