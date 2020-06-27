# api

All deps.cloud API definitions consolidated into a single repository.
This repository currently produces 2 libraries:

* npm: [@depscloud/api](https://www.npmjs.com/package/@depscloud/api)
* vgo: [github.com/depscloud/api](https://github.com/depscloud/api)

## Getting Started with NodeJS

To install:

``` bash
npm install --save @depscloud/api
```

Usage:

```javascript
const grpc = require('grpc');

const { DependencyExtractor } = require('@depscloud/api/v1alpha/extractor');
const {
    SourceService,
    ModuleService,
    DependencyService,
} = require('@depscloud/api/v1alpha/tracker');

const credentials = grpc.credentials.createInsecure();

const dependencyExtractor = new DependencyExtractor('gateway:80', credentials);
const sourceService = new SourceService('gateway:80', credentials);
const moduleService = new ModuleService('gateway:80', credentials);
const dependencyService = new DependencyService('gateway:80', credentials);
```

## Getting Started with Go

To install:

```bash
go get -u github.com/depscloud/api
```

Usage:

```go
package main

import (
    "github.com/depscloud/api/v1alpha/extractor"
    "github.com/depscloud/api/v1alpha/tracker"

    "google.golang.org/grpc"
    "google.golang.org/grpc/credentials"
)

func dial(target string) *grpc.ClientConn {
    cc, _ := grpc.Dial(target, grpc.WithBlock(), grpc.WithInsecure())
    return cc
}

func main() {
    dependencyExtractor := extractor.NewDependencyExtractorClient(dial("gateway:80"))
    sourceService := tracker.NewSourceServiceClient(dial("gateway:80"))
    moduleService := tracker.NewModuleServiceClient(dial("gateway:80"))
    dependencyService := tracker.NewDependencyServiceClient(dial("gateway:80"))
}
```

## Recompiling Protocol Buffer Files

I've provided a docker container that encapsulates the dependencies for regenerating the different language files.

```bash
docker run --rm \
    -v $(PWD):/go/src/github.com/depscloud/api \
    depscloud/api-builder \
    bash scripts/compile.sh
```

You can quickly invoke this command using the `make compile-docker` target.

**NOTE:** Changes to `*.js` and `*.d.ts` are done manually.
Evaluating code generation options.
