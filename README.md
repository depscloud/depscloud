# deps

A simple command line interface that makes the API a little more digestible.

## Support

![GitHub](https://img.shields.io/github/license/depscloud/cli.svg)
![branch](https://github.com/depscloud/cli/workflows/branch/badge.svg?branch=main)
![Google Analytics](https://www.google-analytics.com/collect?v=1&cid=555&t=event&ec=repo&ea=open&dp=depscloud%2Fcli&dt=depscloud%2Fcli&tid=UA-143087272-2)
[![FOSSA Status](https://app.fossa.com/api/projects/git%2Bgithub.com%2Fdepscloud%2Fcli.svg?type=shield)](https://app.fossa.com/projects/git%2Bgithub.com%2Fdepscloud%2Fcli?ref=badge_shield)

## Cheat Sheet

**List modules within the source repository**

```bash
$ deps get modules -u https://github.com/depscloud/api.git
{"manages":{"language":"go","system":"vgo","version":"latest"},"module":{"language":"go","organization":"github.com","module":"depscloud/api"}}
{"manages":{"language":"node","system":"npm","version":"0.1.0"},"module":{"language":"node","organization":"depscloud","module":"api"}}
```

**List source repositories for the given module**

```bash
$ deps get sources -l go -o github.com -m depscloud/api
{"source":{"url":"https://github.com/depscloud/api.git"},"manages":{"language":"go","system":"vgo","version":"latest"}}
```

**List modules that depend on a given module**

```bash
$ deps get dependents -l go -o github.com -m depscloud/api
{"depends":{"language":"go","version_constraint":"v0.1.0","scopes":["direct"]},"module":{"language":"go","organization":"github.com","module":"depscloud/gateway"}}
{"depends":{"language":"go","version_constraint":"v0.1.0","scopes":["direct"]},"module":{"language":"go","organization":"github.com","module":"depscloud/tracker"}}
{"depends":{"language":"go","version_constraint":"v0.1.0","scopes":["direct"]},"module":{"language":"go","organization":"github.com","module":"depscloud/indexer"}}
```

**List modules that a given module depends on**

```bash
$ deps get dependencies -l go -o github.com -m depscloud/api
{"depends":{"language":"go","version_constraint":"v1.3.0","scopes":["direct"]},"module":{"language":"go","organization":"github.com","module":"gogo/protobuf"}}
{"depends":{"language":"go","version_constraint":"v0.3.2","scopes":["indirect"]},"module":{"language":"go","organization":"golang.org","module":"x/text"}}
{"depends":{"language":"go","version_constraint":"v0.0.0-20190628185345-da137c7871d7","scopes":["indirect"]},"module":{"language":"go","organization":"golang.org","module":"x/net"}}
{"depends":{"language":"go","version_constraint":"v1.3.2","scopes":["direct"]},"module":{"language":"go","organization":"github.com","module":"golang/protobuf"}}
{"depends":{"language":"go","version_constraint":"v0.0.0-20190916214212-f660b8655731","scopes":["direct"]},"module":{"language":"go","organization":"google.golang.org","module":"genproto"}}
{"depends":{"language":"go","version_constraint":"v0.0.0-20190626221950-04f50cda93cb","scopes":["indirect"]},"module":{"language":"go","organization":"golang.org","module":"x/sys"}}
{"depends":{"language":"go","version_constraint":"v1.11.2","scopes":["direct"]},"module":{"language":"go","organization":"github.com","module":"grpc-ecosystem/grpc-gateway"}}
{"depends":{"language":"go","version_constraint":"v1.23.1","scopes":["direct"]},"module":{"language":"go","organization":"google.golang.org","module":"grpc"}}
```

## License
[![FOSSA Status](https://app.fossa.com/api/projects/git%2Bgithub.com%2Fdepscloud%2Fcli.svg?type=large)](https://app.fossa.com/projects/git%2Bgithub.com%2Fdepscloud%2Fcli?ref=badge_large)
