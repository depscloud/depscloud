# deps.cloud

## Kubernetes Manifests

First up, let's create the common resources.
The MySQL deployment is optional, so feel free to leave that out.
You will need to configure some credentials later on.

```
$ kubectl create ns depscloud-system
$ kubectl apply -n depscloud-system -f https://deps-cloud.github.io/deploy/k8s/mysql.yaml
$ kubectl apply -n depscloud-system -f https://deps-cloud.github.io/deploy/k8s/depscloud-system.yaml
```

By default, the depscloud-system doesn't know anything about the MySQL being deployed.
This allows deployments to bring their own existing data store or leverage hosted solutions in the cloud.
To connect deps.cloud to the MySQL deployed above, simply supply the following configuration.

```bash
$ cat <<EOF | kubectl apply -f -
apiVersion: v1
kind: Secret
metadata:
  namespace: depscloud-system
  name: depscloud-tracker
stringData:
  STORAGE_DRIVER: mysql
  STORAGE_ADDRESS: user:password@tcp(mysql:3306)/depscloud
  STORAGE_READ_ONLY_ADDRESS: user:password@tcp(mysql-slave:3306)/depscloud
EOF
```

Be sure to change the username, password, target, and name appropriately.
Once this configuration is provided, the tracker pods should start up without an issue.

Once the tracker is configured and running, we can configure the indexer.
The indexer needs a config.yaml file to bootstrap the indexer with the repos it's intended to crawl.
The configuration below demonstrates how the indexer can be configured to index the deps-cloud projects.

```bash
$ cat <<EOF | kubectl apply -f -
apiVersion: v1
kind: Secret
metadata:
  namespace: depscloud-system
  name: depscloud-indexer
stringData:
  config.yaml: |
    accounts:
    - github:
        strategy: HTTP
        organizations:
        - deps-cloud
EOF
```

## Helm Charts

The canonical source for Helm charts is the [Helm Hub](https://hub.helm.sh/), an aggregator for distributed chart repos.

This GitHub project is the source for the `depscloud` [Helm chart repository](https://v3.helm.sh/docs/topics/chart_repository/).

For more information about installing and using Helm, see the [Helm Docs](https://helm.sh/docs/).
For a quick introduction to Charts, see the [Chart Guide](https://helm.sh/docs/topics/charts/).

### How do I install these charts?

```
$ helm repo add depscloud https://deps-cloud.github.io/deploy/charts
"depscloud" has been added to your repositories
```

### Resources

* [Chart Sources](https://github.com/deps-cloud/deploy)
* [deps.cloud](https://deps.cloud)
