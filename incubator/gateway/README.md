# Gateway

Gateway is an HTTP proxy to the various gRPC services within the [deps.cloud](https://deps.cloud) ecosystem.
It translates RESTful API calls into gRPC calls for interacting with modules and their dependencies.

## Introduction

This chart bootstraps a Gateway deployment on a [Kubernetes](http://kubernetes.io) cluster using the [Helm](https://helm.sh) package manager.

## Prerequisites

- Kubernetes 1.15+
- Helm 3.0+

## Installing the Chart

To install the chart with the release name `my-release`:

```bash
$ helm install my-release depscloud-stable/tracker
```

The command deploys Gateway on the Kubernetes cluster in the default configuration.
The [Parameters](#parameters) section lists the parameters that can be configured during installation.

> **Tip**: List all releases using `helm list`

## Uninstalling the Chart

To uninstall/delete the `my-release` deployment:

```bash
$ helm delete my-release
```

The command removes all the Kubernetes components associated with the chart and deletes the release.

## Parameters

The following table lists the configurable parameters of the Gateway chart and their default values.

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
| `service.type`                              | The type of service used to address the tracker     | `ClusterIP`                   |
| `service.port`                              | The port that should be exposed through the service | `80`                          |
| `resources`                                 | Any resource constraints to place on the container  | `{}`                          |
| `nodeSelector`                              | Target the deployment to a certain class of nodes   | `{}`                          |
| `tolerations`                               | Identify any taints the process can tolerate        | `[]`                          |
| `affinity`                                  | Set up an an affinity based on attributes           | `{}`                          |
| `extractor.address`                         | The address of the extractor process                | `"{{ .Release.Name }}-extractor:80"` |
| `extractor.secretName`                      | The name of the secret containing certificates to the extractor | `""`              |
| `tracker.address`                           | The address of the tracker process                  | `"{{ .Release.Name }}-tracker:80"` |
| `tracker.secretName`                        | The name of the secret containing certificates to the tracker | `""`                |
| `tls.secretName`                            | The name of the secret container certificate data for mutual TLS | `""`             |
