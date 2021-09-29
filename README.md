# deps.cloud

_⚠️ After much internal conflict, I've decided to move this project into maintenance mode. This comes after a long 3+ 
years of working on this project in open source with little to no involvement from others. After trying to present this 
at several conferences, it's clear either the community isn't ready for or does not require such a building block. Most
individuals that have come to the project open issues, but have not seemed interested in contributing anything beyond
a ticket. **IF interest picks up again, I'm always happy to take the project off the back burner.** For now, I'm just
too burnt out managing a project that doesn't seem wanted / needed / desired._

deps.cloud is a tool to help companies understand what libraries and projects their systems use.
It works by detecting dependencies defined in common [manifest files] (`pom.xml`, `package.json`, `go.mod`, etc).
Using this information, we’re able to answer questions about project dependencies.

* What versions of _k8s.io/client-go_ do we depend on?
* Which projects use _eslint_ as a non-dev dependency?
* What open source libraries do we use the most?

[manifest files]: https://deps.cloud/docs/concepts/terminology/#manifest-file

## To start using deps.cloud

See our documentation on [deps.cloud](https://deps.cloud/docs/).

## To start developing deps.cloud

Take a look at our [contributing guidelines] and [project board].

```bash
# setup a workspace for all depscloud
mkdir depscloud && cd $_

# clone repository
git clone git@github.com:depscloud/depscloud.git
```

[contributing guidelines]: https://deps.cloud/docs/contrib/
[project board]: https://github.com/orgs/depscloud/projects/1

# Support

Join our [mailing list] to get access to virtual events and ask any questions there.

We also have a [Slack] channel.

[mailing list]: https://groups.google.com/a/deps.cloud/forum/#!forum/community/join
[Slack]: https://depscloud.slack.com/join/shared_invite/zt-fd03dm8x-L5Vxh07smWr_vlK9Qg9q5A

## Checks

![](https://img.shields.io/badge/dynamic/json?color=blue&label=api.deps.cloud%20status&query=%24.state&url=https%3A%2F%2Fapi.deps.cloud%2Fhealth)

**Branch**

[![branch workflow](https://github.com/depscloud/depscloud/workflows/branch/badge.svg?branch=main)](https://github.com/depscloud/depscloud/actions?query=workflow%3Abranch+branch%3Amain)
[![coverage](https://img.shields.io/codecov/c/gh/depscloud/depscloud/main)](https://codecov.io/gh/depscloud/depscloud)
[![dockerfiles workflow](https://github.com/depscloud/depscloud/workflows/dockerfiles/badge.svg?branch=main)](https://github.com/depscloud/depscloud/actions?query=workflow%3Adockerfiles+branch%3Amain)
[![goreleaser branch workflow](https://github.com/depscloud/depscloud/workflows/goreleaser-branch/badge.svg?branch=main)](https://github.com/depscloud/depscloud/actions?query=workflow%3Agoreleaser-branch+branch%3Amain)

**Release**

[![nightly workflow](https://github.com/depscloud/depscloud/workflows/nightly/badge.svg?branch=main)](https://github.com/depscloud/depscloud/actions?query=workflow%3Anightly+branch%3Amain)
[![extractor-tag workflow](https://github.com/depscloud/depscloud/workflows/extractor-tag/badge.svg)](https://github.com/depscloud/depscloud/actions?query=workflow%3Aextractor-tag)
[![goreleaser-tag workflow](https://github.com/depscloud/depscloud/workflows/goreleaser-tag/badge.svg)](https://github.com/depscloud/depscloud/actions?query=workflow%3Agoreleaser-tag)

**Image**

[![extractor docker hub](https://img.shields.io/docker/v/depscloud/extractor?color=blue&label=extractor%20version&sort=semver)](https://hub.docker.com/r/depscloud/extractor/tags)
[![extractor image](https://img.shields.io/docker/image-size/depscloud/extractor?label=extractor%20image&sort=semver)](https://hub.docker.com/r/depscloud/extractor/tags)
[![extractor pulls](https://img.shields.io/docker/pulls/depscloud/extractor?label=extractor%20pulls)](https://hub.docker.com/r/depscloud/extractor/tags)

[![gateway docker hub](https://img.shields.io/docker/v/depscloud/gateway?color=blue&label=gateway%20version&sort=semver)](https://hub.docker.com/r/depscloud/gateway/tags)
[![gateway image](https://img.shields.io/docker/image-size/depscloud/gateway?label=gateway%20image&sort=semver)](https://hub.docker.com/r/depscloud/gateway/tags)
[![gateway pulls](https://img.shields.io/docker/pulls/depscloud/gateway?label=gateway%20pulls)](https://hub.docker.com/r/depscloud/gateway/tags)

[![indexer docker hub](https://img.shields.io/docker/v/depscloud/indexer?color=blue&label=indexer%20version&sort=semver)](https://hub.docker.com/r/depscloud/indexer/tags)
[![indexer image](https://img.shields.io/docker/image-size/depscloud/indexer?label=indexer%20image&sort=semver)](https://hub.docker.com/r/depscloud/indexer/tags)
[![indexer pulls](https://img.shields.io/docker/pulls/depscloud/indexer?label=indexer%20pulls)](https://hub.docker.com/r/depscloud/indexer/tags)

[![tracker docker hub](https://img.shields.io/docker/v/depscloud/tracker?color=blue&label=tracker%20version&sort=semver)](https://hub.docker.com/r/depscloud/tracker/tags)
[![tracker image](https://img.shields.io/docker/image-size/depscloud/tracker?label=tracker%20image&sort=semver)](https://hub.docker.com/r/depscloud/tracker/tags)
[![tracker pulls](https://img.shields.io/docker/pulls/depscloud/tracker?label=tracker%20pulls)](https://hub.docker.com/r/depscloud/tracker/tags)

**License**

[![fossa](https://app.fossa.com/api/projects/git%2Bgithub.com%2Fdepscloud%2Fdepscloud.svg?type=large)](https://app.fossa.com/projects/git%2Bgithub.com%2Fdepscloud%2Fdepscloud?ref=badge_large)
![analytics](https://www.google-analytics.com/collect?v=1&cid=555&t=pageview&ec=repo&ea=open&dp=depscloud&dt=depscloud&tid=UA-143087272-2)
