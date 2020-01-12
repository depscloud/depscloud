# Extractor

Extractor encapsulates the logic for extracting dependency information from repositories in the [deps.cloud](https://deps.cloud) ecosystem.
It provides a gRPC API for matching and extracting contents of dependency manifest files.

## Introduction

This chart bootstraps a Extractor deployment on a [Kubernetes](http://kubernetes.io) cluster using the [Helm](https://helm.sh) package manager.

## Prerequisites

- Kubernetes 1.15+
- Helm 3.0+

## Installing the Chart

To install the chart with the release name `my-release`:

```bash
$ helm install my-release depscloud-incubator/tracker
```

The command deploys Extractor on the Kubernetes cluster in the default configuration.
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

| Parameter                                   | Description                                         | Default                       |
|---------------------------------------------|-----------------------------------------------------|-------------------------------|
| `replicaCount`                              | The number of instances to run                      | `1`                           |
| `image.repository`                          | The address of the registry hosting the image       | `depscloud/indexer`           |
| `image.pullPolicy`                          | The pull policy for the image                       | `IfNotPresent`                |
| `image.tag`                                 | The version of the image                            | `.Chart.AppVersion`           |
| `imagePullSecrets`                          | Registry secret names                               | `[]`                          |
| `nameOverride`                              | String to partially override the full name template | `""`                          |
| `fullnameOverride`                          | String to completely override the full name         | `""`                          |
| `serviceAccount.create`                     | Whether or not a service account should be created  | `true`                        |
| `serviceAccount.name`                       | The name of the service account                     | `""`                          |
| `podSecurityContext`                        | Provide any pod security context attributes         | `{}`                          |
| `securityContext`                           | Provide any security context attributes             | `{}`                          |
| `service.type`                              | The type of service used to address the tracker     | `Headless`                    |
| `service.port`                              | The port that should be exposed through the service | `8090`                        |
| `resources`                                 | Any resource constraints to place on the container  | `{}`                          |
| `nodeSelector`                              | Target the deployment to a certain class of nodes   | `{}`                          |
| `tolerations`                               | Identify any taints the process can tolerate        | `[]`                          |
| `affinity`                                  | Set up an an affinity based on attributes           | `{}`                          |
| `tls.secretName`                            | The name of the secret container certificate data for mutual TLS | `""`             |
