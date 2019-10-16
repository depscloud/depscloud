# api

All deps-cloud API definitions consolidated into a single repostory.
This repository currenlty produces 2 libraries:

* npm: [@deps-cloud/api](https://www.npmjs.com/package/@deps-cloud/api)
* vgo: [github.com/deps-cloud/api](https://github.com/deps-cloud/api)

## Getting Started with NodeJS

To install:

``` bash
npm install --save @deps-cloud/api
```

Usage:

```javascript
const grpc = require('grpc');

const { DependencyExtractor } = require('@deps-cloud/api/v1alpha/extractor/extractor');
const {
    SourceService,
    ModuleService,
    DependencyService,
} = require('@deps-cloud/api/v1alpha/tracker/tracker');

const credentials = grpc.credentials.createInsecure();

const dependencyExtractor = new DependencyExtractor('address', credentials);
const sourceService = new SourceService('address', credentials);
const moduleService = new ModuleService('address', credentials);
const dependencyService = new DependencyService('address', credentials);
```

## Getting Started with Go

To install:

```bash
go get -u github.com/deps-cloud/api
```

Usage:

```go
package main

import (
    "github.com/deps-cloud/api/v1alpha/extractor"
    "github.com/deps-cloud/api/v1alpha/tracker"

    "google.golang.org/grpc"
    "google.golang.org/grpc/credentials"
)

func dial(target string) *grpc.ClientConn {
    cc, _ := grpc.Dial(target, grpc.WithBlock(), grpc.WithInsecure())
    return cc
}

func main() {
    dependencyExtractor := extractor.NewDependencyExtractorClient(dial("address"))
    sourceService := tracker.NewSourceServiceClient(dial("address"))
    moduleService := tracker.NewModuleServiceClient(dial("address"))
    dependencyService := tracker.NewDependencyServiceClient(dial("address"))
}
```

## Recompiling Protocol Buffer Files

I've provided a docker container that encapsulates the dependencies for regenerating the different language files.

```bash
docker run --rm -it \
    -v $(pwd)/swagger:/go/src/github.com/deps-cloud/api/swagger \
    -v $(pwd)/v1alpha:/go/src/github.com/deps-cloud/api/v1alpha \
    depscloud/api-builder
```

You can quickly invoke this command using the `make compile-docker` target.

**NOTE:** Changes to `*.js` and `*.d.ts` are done manually.
Currently evaluating coge generation options.
