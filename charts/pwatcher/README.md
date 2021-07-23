# pwatcher

Example Kubernetes Controller built with the `kubebuilder` to showcase different use-cases in a controller.

## Installing the Chart

To install the chart with the release name `my-release`:

```console
helm repo add pwatcher https://ahmedwaleedmalik.github.io/pwatcher
helm repo update
helm upgrade -i my-release pwatcher/pwatcher
```

The command deploys podinfo on the Kubernetes cluster in the default namespace.
The [configuration](#configuration) section lists the parameters that can be configured during installation.

## Uninstalling the Chart

To uninstall/delete the `my-release` deployment:

```console
helm delete my-release
```

## Configuration

The following tables lists the configurable parameters of the pwatcher chart and their default values:

| Parameter                 | Default                     | Description                    |
| ------------------------- | --------------------------- | ------------------------------ |
| `replicaCount`            | `1`                         | Desired number of pod replicas |
| `image.repository`        | `ahmedwaleedmalik/pwatcher` | Image repository               |
| `image.tag`               | `<chart.version>`           | Image tag                      |
| `image.pullPolicy`        | `IfNotPresent`              | Image pull policy              |
| `podFilter.enabled`       | `false`                     | Enable pod filter              |
| `podFilter.key`           | `timestamp`                 | Pod filter key                 |
| `namespaceFilter.enabled` | `false`                     | Enable namespace filter        |
| `namespaceFilter.key`     | `timestamp`                 | Namespace filter key           |
