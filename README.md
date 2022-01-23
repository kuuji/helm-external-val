# helm-external-val

## Overview

`helm-external-val` is a helm plugin that allows storing helm values in external source.
Currently it supports getting values from kubernetes [ConfigMaps](https://kubernetes.io/docs/concepts/configuration/configmap/)

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

#### Downloader plugin

The url is formatted as follows 

```
<source>://<namespace>/<name>
```

for example the url below will fetch the ConfigMap named `helm-values` from the namespace `kuuji`.

```
cm://kuuji/helm-values
```

Note: `namespace` is optional, not providing it will default to `default`