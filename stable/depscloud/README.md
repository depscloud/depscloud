# deps.cloud

deps.cloud is a set of services used to extract, track, and monitor dependencies between systems
It is comprised of a distinct set of services, each working together to accomplish a task. 

## Introduction

This chart bootstraps the deps.cloud ecosystem on a [Kubernetes](http://kubernetes.io) cluster using the [Helm](https://helm.sh) package manager.

## Prerequisites

- Kubernetes 1.15+
- Helm 3.0+

## Installing the Chart

To install the chart with the release name `my-release`:

```bash
$ helm install my-release depscloud-incubator/depscloud
```

The command deploys deps.cloud on the Kubernetes cluster in the default configuration.
The [Parameters](#parameters) section lists the parameters that can be configured during installation.

> **Tip**: List all releases using `helm list`

## Uninstalling the Chart

To uninstall/delete the `my-release` deployment:

```bash
$ helm delete my-release
```

The command removes all the Kubernetes components associated with the chart and deletes the release.

## Parameters

The following table lists the configurable parameters of the Extractor chart and their default values.

| Parameter           | Description                                         | Default                       |
|---------------------|-----------------------------------------------------|-------------------------------|
| `extractor.enabled` | Whether the `extractor` process should be deployed  | `true`                        |
| `extractor.*`       | Set configuration for the `extractor` chart         | See the `extractor` chart     |
| `gateway.enabled`   | Whether the `gateway` process should be deployed    | `true`                        |
| `gateway.*`         | Set configuration for the `gateway` chart           | See the `gateway` chart       |
| `indexer.enabled`   | Whether the `indexer` process should be deployed    | `true`                        |
| `indexer.*`         | Set configuration for the `indexer` chart           | See the `indexer` chart       |
| `tracker.enabled`   | Whether the `tracker` process should be deployed    | `true`                        |
| `tracker.*`         | Set configuration for the `tracker` chart           | See the `tracker` chart       |
