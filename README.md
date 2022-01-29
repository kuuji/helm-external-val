# helm-external-val

## Overview

`helm-external-val` is a helm plugin that fetches helm values from external source.
Currently it supports getting values from kubernetes [ConfigMaps](https://kubernetes.io/docs/concepts/configuration/configmap/) and kubernetes [Secrets](https://kubernetes.io/docs/concepts/configuration/secret/)

## Installation

### Local

```
helm plugin install https://github.com/kuuji/helm-external-val
```

or by specifying the version (git tag)

```
helm plugin install --version v0.0.4 https://github.com/kuuji/helm-external-val
```


### ArgoCD

### Via a custom image

The ArgoCD [recommended option](https://argo-cd.readthedocs.io/en/stable/user-guide/helm/#helm-plugins) is to build a repo-server image that includes the plugin.
See example below

```
FROM argoproj/argocd:v1.5.7

USER root
RUN apt-get update && \
    apt-get install -y \
        curl && \
    apt-get clean && \
    rm -rf /var/lib/apt/lists/* /tmp/* /var/tmp/*

USER argocd

RUN helm plugin install https://github.com/kuuji/helm-external-val

ENV HELM_PLUGINS="/home/argocd/.local/share/helm/plugins/"
```

### Via init container

Using the helm chart or by patching repo-server we can use an init-container to install helm-external-val in a shared `custom-tools` volume.
Then we set the helm plugins directory to `/custom-tools/helm-plugins` in the main repo-server container via the env `HELM_PLUGIN`
Credit to @jkroepke for this neat trick

```
repoServer:
  env:
    - name: HELM_PLUGINS
      value: /custom-tools/helm-plugins/
  volumes:
    - name: custom-tools
      emptyDir: {}
  volumeMounts:
    - mountPath: /custom-tools
      name: custom-tools
  initContainers:
  - name: download-tools
    args:
    - |
      mkdir -p /custom-tools/helm-plugins
      helm plugin install https://github.com/kuuji/helm-external-val
    command:
    - sh
    - -ec
    image: alpine/helm:latest
    imagePullPolicy: Always
    terminationMessagePath: /dev/termination-log
    terminationMessagePolicy: File
    volumeMounts:
    - mountPath: /custom-tools
      name: custom-tools
    env:
    - name: HELM_PLUGINS
      value: /custom-tools/helm-plugins
```

### Caveats

ArgoCD won't re-evaluate the external source unless you do a hard refresh on the application. This is because argocd caches the application manifest and won't evaluate the external values since the application manifest didn't change.

You can do this either in the ui (click the arrow under the refresh button) or via the cli like below.

```
argocd app get <application_name> --hard-refresh
```

## Usage

This plugin has 2 modes of operation.

- Using the plugin cli
  - Create a values.yaml locally from various sources

- As a downloader plugin
  - Feed a specially formatted url to `helm install|upgrade -f`

The latter is recommended as it fits well with gitops workflows.


### CLI plugin

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

### Downloader plugin

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
