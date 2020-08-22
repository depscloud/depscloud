# deps.cloud

<img alt="logo" width="64" src="https://deps.cloud/favicons/android-chrome-512x512.png"/>

deps.cloud is a tool built to help companies understand how projects relate to one another.
It does this by detecting dependencies defined in common manifest files.
Using this information, we’re able to construct a dependency graph.
As a result we’re able to answer questions like:

* Which libraries get produced by a project?
* Which libraries do I depend on and what version?
* Which projects depend on library X and what version?
* Which projects can produce library X?
* Which projects do our systems use the most?

## To start using deps.cloud

See our documentation on [deps.cloud](https://deps.cloud/docs/).

## To start developing deps.cloud

Take a look at our [contributing guidelines] and [project board].

```bash
# setup a workspace for all depscloud
mkdir depscloud && cd $_

# clone necessary repositories
#   - the first is for the source code
#   - the second is for the deployment configuration
git clone git@github.com:depscloud/depscloud.git
git clone git@github.com:depscloud/deploy.git
```

[contributing guidelines]: https://deps.cloud/docs/contrib/
[project board]: https://github.com/orgs/depscloud/projects/1

# Support

Join our [mailing list] and ask any questions there.

We also have a [Slack] channel.

[mailing list]: https://groups.google.com/a/deps.cloud/forum/#!forum/community/join
[Slack]: https://depscloud.slack.com/join/shared_invite/zt-fd03dm8x-L5Vxh07smWr_vlK9Qg9q5A

## Branch Checks

![deps](https://github.com/depscloud/depscloud/workflows/deps-branch/badge.svg?branch=main)
![extractor](https://github.com/depscloud/depscloud/workflows/extractor-branch/badge.svg?branch=main)
![gateway](https://github.com/depscloud/depscloud/workflows/gateway-branch/badge.svg?branch=main)
![indexer](https://github.com/depscloud/depscloud/workflows/indexer-branch/badge.svg?branch=main)
![tracker](https://github.com/depscloud/depscloud/workflows/tracker-branch/badge.svg?branch=main)
![goreleaser](https://github.com/depscloud/depscloud/workflows/goreleaser-branch/badge.svg?branch=main)

## Release Checks

![extractor](https://github.com/depscloud/depscloud/workflows/extractor-tag/badge.svg)
![goreleaser](https://github.com/depscloud/depscloud/workflows/goreleaser-tag/badge.svg)

![extractor](https://img.shields.io/docker/v/depscloud/extractor?color=blue&label=extractor%20version&sort=semver)
![gateway](https://img.shields.io/docker/v/depscloud/gateway?color=blue&label=gateway%20version&sort=semver)
![indexer](https://img.shields.io/docker/v/depscloud/indexer?color=blue&label=indexer%20version&sort=semver)
![tracker](https://img.shields.io/docker/v/depscloud/tracker?color=blue&label=tracker%20version&sort=semver)

## License Checks

[![FOSSA Status](https://app.fossa.com/api/projects/git%2Bgithub.com%2Fdepscloud%2Fdepscloud.svg?type=large)](https://app.fossa.com/projects/git%2Bgithub.com%2Fdepscloud%2Fdepscloud?ref=badge_large)
![analytics](https://www.google-analytics.com/collect?v=1&cid=555&t=pageview&ec=repo&ea=open&dp=depscloud&dt=depscloud&tid=UA-143087272-2)
