# Development Guide

This monorepo contains several independent processes.

* `deps` is a [command line tool] for interacting with the deps.cloud API.
* `extractor` accepts [manifest files] and extracts relationships from them.
* `gateway` provides both RESTful and gRPC interfaces. 
* `indexer` crawls repositories, calling the extractor and tracker appropriately. 
* `tracker` manages the dependency graph built on top of common databases.

[command line tool]: https://deps.cloud/docs/guides/cli/
[manifest files]: https://deps.cloud/docs/concepts/manifests/   

## Cloning projects

For most development on this project, you will need two repositories.
The first is this repository (depscloud) which contains all the source code.
The second is for the deployment configuration.

```bash
# setup a workspace for all depscloud
mkdir depscloud && cd $_

# clone necessary repositories
#   - the first is for the source code
#   - the second is for the deployment configuration
git clone git@github.com:depscloud/depscloud.git
git clone git@github.com:depscloud/deploy.git
```

## Building changes

Every component can be built using docker.
When building a container locally it's tagged using the `latest` tag.
This allows it to be deployed using our [docker] configuration.
A common workflow is to build the changes to your container and redeploy the docker stack.

```bash
# in depscloud/depscloud
# make [name]/docker
make tracker/docker

# in depscloud/deploy/docker/sqlite
docker-compose up
```

[docker]: https://deps.cloud/docs/deploy/docker/
