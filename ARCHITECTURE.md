# Architecture

I've had this information around the site/repository in one form or another.
Looks like this has picked up from the Rust community.
The default setup for deps.cloud deploys a full read-write system.

[![](https://mermaid.ink/img/eyJjb2RlIjoiZ3JhcGggTFJcbiAgc3ViZ3JhcGggY2x1c3RlclxuICAgIHRyYWNrZXIgLS0-IG15c3FsXG4gICAgdHJhY2tlciAtLT4gcG9zdGdyZXNxbFxuXG4gICAgZ2F0ZXdheSAtLT4gdHJhY2tlclxuXG4gICAgaW5kZXhlciAtLT4gdHJhY2tlclxuICAgIGluZGV4ZXIgLS0-IGV4dHJhY3RvclxuICBlbmRcbiAgXG4gIGRlcHMgLS0gZGVwc2Nsb3VkLmNvbXBhbnkuY29tIC0tPiBnYXRld2F5XG4iLCJtZXJtYWlkIjp7fSwidXBkYXRlRWRpdG9yIjpmYWxzZX0)](https://mermaid-js.github.io/mermaid-live-editor/#/edit/eyJjb2RlIjoiZ3JhcGggTFJcbiAgc3ViZ3JhcGggY2x1c3RlclxuICAgIHRyYWNrZXIgLS0-IG15c3FsXG4gICAgdHJhY2tlciAtLT4gcG9zdGdyZXNxbFxuXG4gICAgZ2F0ZXdheSAtLT4gdHJhY2tlclxuXG4gICAgaW5kZXhlciAtLT4gdHJhY2tlclxuICAgIGluZGV4ZXIgLS0-IGV4dHJhY3RvclxuICBlbmRcbiAgXG4gIGRlcHMgLS0gZGVwc2Nsb3VkLmNvbXBhbnkuY29tIC0tPiBnYXRld2F5XG4iLCJtZXJtYWlkIjp7fSwidXBkYXRlRWRpdG9yIjpmYWxzZX0)

Deployments can adopt a read-only mode by deploying a tracker with read-only permissions on the database.
The tracker handles read-only and read-write connections separately, so you have some degree of safety by only setting the read-only address. 
Any attempt to write data will result in an unsupported operation issue.
`api.deps.cloud` is run using the read-only configuration.
 
[![](https://mermaid.ink/img/eyJjb2RlIjoiZ3JhcGggTFJcbiAgc3ViZ3JhcGggY2x1c3RlclxuICAgIHRyYWNrZXItcmVhZG9ubHkgLS0gcmVhZC1vbmx5IC0tPiBteXNxbC9wb3N0Z3Jlc3FsXG4gICAgdHJhY2tlciAtLSByZWFkL3dyaXRlIC0tPiBteXNxbC9wb3N0Z3Jlc3FsXG5cbiAgICBnYXRld2F5IC0tPiB0cmFja2VyLXJlYWRvbmx5XG5cbiAgICBpbmRleGVyIC0tPiB0cmFja2VyXG4gICAgaW5kZXhlciAtLT4gZXh0cmFjdG9yXG4gIGVuZFxuICBcbiAgZGVwcyAtLSBhcGkuZGVwcy5jbG91ZCAtLT4gZ2F0ZXdheVxuIiwibWVybWFpZCI6e30sInVwZGF0ZUVkaXRvciI6ZmFsc2V9)](https://mermaid-js.github.io/mermaid-live-editor/#/edit/eyJjb2RlIjoiZ3JhcGggTFJcbiAgc3ViZ3JhcGggY2x1c3RlclxuICAgIHRyYWNrZXItcmVhZG9ubHkgLS0gcmVhZC1vbmx5IC0tPiBteXNxbC9wb3N0Z3Jlc3FsXG4gICAgdHJhY2tlciAtLSByZWFkL3dyaXRlIC0tPiBteXNxbC9wb3N0Z3Jlc3FsXG5cbiAgICBnYXRld2F5IC0tPiB0cmFja2VyLXJlYWRvbmx5XG5cbiAgICBpbmRleGVyIC0tPiB0cmFja2VyXG4gICAgaW5kZXhlciAtLT4gZXh0cmFjdG9yXG4gIGVuZFxuICBcbiAgZGVwcyAtLSBhcGkuZGVwcy5jbG91ZCAtLT4gZ2F0ZXdheVxuIiwibWVybWFpZCI6e30sInVwZGF0ZUVkaXRvciI6ZmFsc2V9)

## Components

`deps` - the command line interface (CLI) that allows users to easily query the graph for information.

`extractor` - a TypeScript gRPC service.
It's used to parse and extract dependencies from files. 

`gateway` - a Go HTTP / gRPC service.
It mediates communication with the backing services and provides RESTful interface to query the graph. 

`indexer` - a Go CronJob.
It pulls repository information and update changes in the graph.
To do this, it shallow-clones the repository and walks its file system for information. 

`tracker` - A Go gRPC service.
It encapsulates two servers.
The first is a GraphStore which provides key-value like semantics for managing graph data stored in common SQL solutions.
The second is the domain specific interface which focuses on modeling relationships between sources, modules, and their dependencies.  

## Navigating the repo

`deps/` - This directory contains all the source code specific to the `deps` component.

`docker/` - This directory contains docker-compose configuration for spinning up depscloud with different databases.
This was ported in from the [depscloud/deploy](https://github.com/depscloud/deploy) repository to help simplify development.

`dockerfiles/` - This directory contains various `Dockerfile` used to build different parts of the project.
The only Dockerfile that isn't under this directory is the one used for the `extractor` process.

`extractor/` - This directory contains all the source code specific to the `extractor` component.
This service is entirely stateless.
It's only responsibility is to parse string text for the various files and return a standard dependency format.
This makes it easy to scale out as needed. 

`gateway/` - This directory contains all the source code specific to the `gateway` component.
It serves a generic gRPC proxy that makes it easy to route requests to the appropriate backend services.
It also leverages `grpc-gateway` to map RESTful routes to gRPC services.
Because the upstream `tracker` process leverages long-lived streams, `gateway` must also support them.
This requires a little more attention on the load balancing side of the world.

`indexer/` - This directory contains all the source code specific to the `indexer` component.
While it currently runs as a CronJob, the next step is to turn this more into a controller.
It will be responsible for pulling your upstream VCS provider and determining when changes have occurred.
This will allow for a more real-time update of information.

`internal/` - This directory contains source code shared by many of the Go processes.

`monorepo/` - Mostly scripts I (@mjpitz) have used to merge the once numerous repositories into a single repo.
Using git-subtree made it easy to migrate git-history.
Not deleted yet, as I'm contemplating just how far I want to go down the monorepo route.

`tracker/` - This directory contains all the source code specific to the `tracker` component.
This service is stateless and relatively easy to scale out.
The gRPC services it offers leverage long-lived streams so load balancing requires a little more attention.
