# Getting Started with NodeJS

To install:

``` bash
npm install --save @depscloud/api
```

Usage:

```javascript
const grpc = require('@grpc/grpc-js');

const { DependencyExtractor } = require('@depscloud/api/v1alpha/extractor');
const {
    SourceService,
    ModuleService,
    DependencyService,
} = require('@depscloud/api/v1alpha/tracker');

const target = "api.deps.cloud:443";
const creds = grpc.credentials.createSsl();

const sourceService = new SourceService(target, creds);
const moduleService = new ModuleService(target, creds);
const dependencyService = new DependencyService(target, creds);
```
