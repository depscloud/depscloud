# api

All deps.cloud API definitions consolidated into a single repository.
This repository currently produces 3 libraries:

| Tech   | Source                                  | Package                                                        |
|:-------|:----------------------------------------|:---------------------------------------------------------------|
| npm    | [source](packages/depscloud-api-nodejs) | [@depscloud/api](https://www.npmjs.com/package/@depscloud/api) |
| pip    | [source](packages/depscloud-api-python) | [depscloud_api](https://pypi.org/project/depscloud_api/) (coming soon!)      |
| go mod |                                         | [github.com/depscloud/api](https://github.com/depscloud/api)   |

## Getting Started with Go

To install:

```bash
go get -u github.com/depscloud/api
```

Usage:

```go
package main

import (
    "crypto/tls"

    "github.com/depscloud/api/v1alpha/extractor"
    "github.com/depscloud/api/v1alpha/tracker"

    "google.golang.org/grpc"
    "google.golang.org/grpc/credentials"
)

func main() {
    target := "api.deps.cloud:443"
    creds := credentials.NewTLS(&tls.Config{})

    conn, _ := grpc.Dial(target, grpc.WithTransportCredentials(creds))
    defer conn.Close()

    sourceService := tracker.NewSourceServiceClient(conn)
    moduleService := tracker.NewModuleServiceClient(conn)
    dependencyService := tracker.NewDependencyServiceClient(conn)
}
```
