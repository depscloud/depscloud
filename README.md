# depscloud-cli

A simple command line interface that makes the API a little more digestible.

## Status

Early alpha. Hacked together in a night.

## Cheat Sheet

**List modules within the source repository**

```bash
$ depscloud-cli get modules -u https://github.com/deps-cloud/api.git
{"manages":{"language":"go","system":"vgo","version":"latest"},"module":{"language":"go","organization":"github.com","module":"deps-cloud/api"}}
{"manages":{"language":"node","system":"npm","version":"0.1.0"},"module":{"language":"node","organization":"deps-cloud","module":"api"}}
```

**List source repositories for the given module**

```bash
$ depscloud-cli get sources -l go -o github.com -m deps-cloud/api
{"source":{"url":"https://github.com/deps-cloud/api.git"},"manages":{"language":"go","system":"vgo","version":"latest"}}
```

**List modules that depend on a given module**

```bash
$ depscloud-cli get dependents -l go -o github.com -m deps-cloud/api
{"depends":{"language":"go","version_constraint":"v0.1.0","scopes":["direct"]},"module":{"language":"go","organization":"github.com","module":"deps-cloud/gateway"}}
{"depends":{"language":"go","version_constraint":"v0.1.0","scopes":["direct"]},"module":{"language":"go","organization":"github.com","module":"deps-cloud/tracker"}}
{"depends":{"language":"go","version_constraint":"v0.1.0","scopes":["direct"]},"module":{"language":"go","organization":"github.com","module":"deps-cloud/indexer"}}
```

**List modules that a given module depends on**

```bash
$ depscloud-cli get dependencies -l go -o github.com -m deps-cloud/api
{"depends":{"language":"go","version_constraint":"v1.3.0","scopes":["direct"]},"module":{"language":"go","organization":"github.com","module":"gogo/protobuf"}}
{"depends":{"language":"go","version_constraint":"v0.3.2","scopes":["indirect"]},"module":{"language":"go","organization":"golang.org","module":"x/text"}}
{"depends":{"language":"go","version_constraint":"v0.0.0-20190628185345-da137c7871d7","scopes":["indirect"]},"module":{"language":"go","organization":"golang.org","module":"x/net"}}
{"depends":{"language":"go","version_constraint":"v1.3.2","scopes":["direct"]},"module":{"language":"go","organization":"github.com","module":"golang/protobuf"}}
{"depends":{"language":"go","version_constraint":"v0.0.0-20190916214212-f660b8655731","scopes":["direct"]},"module":{"language":"go","organization":"google.golang.org","module":"genproto"}}
{"depends":{"language":"go","version_constraint":"v0.0.0-20190626221950-04f50cda93cb","scopes":["indirect"]},"module":{"language":"go","organization":"golang.org","module":"x/sys"}}
{"depends":{"language":"go","version_constraint":"v1.11.2","scopes":["direct"]},"module":{"language":"go","organization":"github.com","module":"grpc-ecosystem/grpc-gateway"}}
{"depends":{"language":"go","version_constraint":"v1.23.1","scopes":["direct"]},"module":{"language":"go","organization":"google.golang.org","module":"grpc"}}
```
