# helm-external-val

## Overview

`helm-external-val` is a helm plugin that fetches helm values from external source.
Currently it supports getting values from kubernetes [ConfigMaps](https://kubernetes.io/docs/concepts/configuration/configmap/) and kubernetes [Secrets](https://kubernetes.io/docs/concepts/configuration/secret/)

## Getting started

### Installation

Local install

```
# Clone the project
# Build the project
go build
# Install the plugin
helm plugin install .
```

TODO

### Usage

This plugin has 2 mode of operation.
- Using the plugin cli
  - Create a values.yaml from various sources
- As a downloader plugin
  - Feed a specially formatted url to `helm install|upgrade -f`

The latter is recommended as it fits well with gitops workflows.


#### CLI plugin

```
helm external-val cm -h
Get the content of values from a cm and write it to a file

Usage:
  helm-external-val cm <name> [flags]

Flags:
  -h, --help                    help for cm
      --kube_namespace string   The namespace to get the cm from (default "default")
  -o, --out string              The file to output the values to (default "values-cm.yaml")
```

```
helm external-val secret -h
Get the content of values from a secret and write it to a file

Usage:
  helm-external-val secret <name> [flags]

Flags:
  -h, --help                    help for secret
      --kube_namespace string   The namespace to get the secret from (default "default")
  -o, --out string              The file to output the values to (default "values-secret.yaml")
```

#### Downloader plugin

Helm will invoke the downloader plugin with 4 parameters `certFile keyFile caFile full-URL`. In our case we're ignoring the first 3.

The url has to be formatted as follows 

```
<source>://<namespace>/<name>
```

- source (required) : the protocol to use (`cm` and `secret` are currently supported)
- namespace (optional) : the namespace in which to look for the resource (defaults to `default`)
- name (required) : the name of the resource to fetch

for example the url below will fetch the ConfigMap named `helm-values` from the namespace `kuuji`.

```
cm://kuuji/helm-values
```
