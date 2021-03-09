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

For most development on this project, you will need one repositories.

```bash
# setup a workspace for all depscloud
mkdir depscloud && cd $_

# clone necessary repositories
git clone git@github.com:depscloud/depscloud.git
```

## Building changes

Every component can be built using docker.
When building a container locally it's tagged using the `latest` tag.
This allows it to be deployed using our [docker] configuration.
A common workflow is to build the changes to your container and redeploy the docker stack.

```bash
# make [name]/docker
make tracker/docker

# make run/docker[/platform]
make run/docker
```

By default, we run with a SQLite configuration.
You can run other platforms like CockroachDB, MariaDB, PostgreSQL, and MySQL as well.

## Rebuilding the base images

These images are are used to to build the remainder of the application as well as a base
for the runtime docker images. The default is to build the base images without a registry
prefix:

```bash
make dockerfiles
```

To build the base images locally tagged as required for the downstream builds:

```bash
make USE_REGISTRY=1 dockerfiles
```

Remember that an explicit pull or removal of the locally built images maybe neccessary to
use the public versions

```
docker pull ocr.sh/depscloud/base
docker pull ocr.sh/depscloud/devbase
```

[docker]: https://deps.cloud/docs/deploy/docker/
