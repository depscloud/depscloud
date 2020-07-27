![GitHub](https://img.shields.io/github/license/depscloud/extractor.svg)
![branch](https://github.com/depscloud/extractor/workflows/branch/badge.svg?branch=main)
![Google Analytics](https://www.google-analytics.com/collect?v=1&cid=555&t=event&ec=repo&ea=open&dp=deploy&dt=deploy&tid=UA-143087272-2)

# deps.cloud

[deps.cloud](https://deps.cloud/) is system that helps track and manage library usage across an organization.
Unlike many alternatives, it was built with portability in mind making easy for anyone to get started.

For more information on how to get involved take a look at our [project board](https://github.com/orgs/depscloud/projects/1).

## Deployment

This repository contains deployment configuration for the deps.cloud ecosystem.

**Supported storage drivers:**

* Driver name: `sqlite`
* Connection string: `file::memory:?cache=shared`

* Driver name: `mysql`
* Connection string: `user:password@tcp(depscloud-mysql:3306)/depscloud`

* Driver name: `postgres`
* Connection string: `postgres://user:password@depscloud-postgresql:5432/depscloud`

### Docker

Docker is great for trying the project out or for active development.
I do not recommend using the docker configuration in this repository for production deployments.

* [MySQL](docker/mysql/)
* [PostgreSQL](docker/postgres/)
* [SQLite](docker/sqlite/)

### Kubernetes

* [Raw Manifests](https://deps.cloud/docs/deployment/k8s/) (not recommended for production)
* [ArgoCD](https://github.com/depscloud/deploy/tree/main/examples/argocd)
* [Helm](https://github.com/depscloud/deploy/tree/main/examples/helm)
* FluxCD's [HelmOperator](https://github.com/depscloud/deploy/tree/main/examples/helm-operator)

### Helm Charts

The canonical source for Helm charts is the [Helm Hub](https://hub.helm.sh/), an aggregator for distributed chart repos.

This GitHub project is the source for the `depscloud` [Helm chart repository](https://v3.helm.sh/docs/topics/chart_repository/).

For more information about installing and using Helm, see the [Helm Docs](https://helm.sh/docs/).
For a quick introduction to Charts, see the [Chart Guide](https://helm.sh/docs/topics/charts/).

#### How do I install these charts?

```
$ helm repo add depscloud https://depscloud.github.io/deploy/charts
"depscloud" has been added to your repositories
```
