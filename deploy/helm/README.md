# Ruthenium [helm](https://helm.sh/docs/intro/using_helm/) chart

Helm chart to deploy [Ruthenium](https://github.com/my-cloud/ruthenium) cryptocurrency node




## Prerequisites

* Kubernetes: `>= 1.24.0-0`
* Helm: `>= 3.0`

## Getting Started

For a quick install with the default configuration:

```bash
$ helm install ruthenium --set secrets[0].data.privateKey=<MyPrivateKey>
```

## Source Code

* [Ruthenium](https://github.com/my-cloud/ruthenium)

## Values

| Key                                                    | Type | Default                      | Description                                                                                                                           |
|--------------------------------------------------------|------|------------------------------|---------------------------------------------------------------------------------------------------------------------------------------|
| app.containers[0].annotations                          | map | {}                           | any relevant annotation                                                                                                               |
| app.containers[0].args                                 | list | []                           | any application argument                                                                                                              |
| app.containers[0].autoReload                           | string | "true"                       | specifies if the container should be restarted on each update                                                                         |
| app.containers[0].command[0]                           | string | "/app/validatornode"         | application binary path                                                                                                               |
| app.containers[0].env                                  | map | {}                           | any environment variable `ENV_NAME: value`                                                                                            |
| app.containers[0].health.liveness.initialDelaySeconds  | int | 40                           | kubernetes liveness initialDelaySeconds                                                                                               |
| app.containers[0].health.liveness.periodSeconds        | int | 5                            | kubernetes liveness periodSeconds                                                                                                     |
| app.containers[0].health.readiness.initialDelaySeconds | int | 30                           | kubernetes readiness initialDelaySeconds                                                                                              |
| app.containers[0].health.readiness.periodSeconds       | int | 1                            | kubernetes readiness periodSeconds                                                                                                    |
| app.containers[0].health.type                          | string | "grpc"                       | kubernetes healthcheck type                                                                                                           |
| app.containers[0].image.name                           | string | "ghcr.io/my-cloud/ruthenium" | container image                                                                                                                       |
| app.containers[0].image.pullPolicy                     | string | "IfNotPresent"               | container pull policy                                                                                                                 |
| app.containers[0].image.tagOverride                    | string | "latest"                     | container tag                                                                                                                         |
| app.containers[0].name                                 | string | "node"                       | container name                                                                                                                        |
| app.containers[0].resources.limits.memory              | string | "512Mi"                      | kubernetes resources limits memory                                                                                                    |
| app.containers[0].resources.requests.memory            | string | "128Mi"                      | kubernetes resources requests memory                                                                                                  |
| app.containers[0].secret.ruthenium.PRIVATE_KEY         | string | "privateKey"                 | Gets the `privateKey` key of the `ruthenium` secret and populate a container environment variable `PRIVATE_KEY` with the secret value |
| app.containers[0].service[0].port                      | string | "8106"                       | kubernetes service port                                                                                                               |
| app.containers[0].service[0].protocol                  | string | "TCP"                        | kubernetes service protocol                                                                                                           |
| app.containers[0].service[0].targetPort                | string | "8106"                       | kubernetes service container listening port                                                                                           |
| app.containers[0].storage.data                         | string | "/tmp"                       | will mount the `data` volume in the specified path                                                                                    |
| app.containers[1].annotations                          | map | {}                           | any relevant annotation                                                                                                               |
| app.containers[1].autoReload                           | string | "true"                       | specifies if the container should be restarted on each update                                                                         |
| app.containers[1].command[0]                           | string | "/app/ui"                    | any application argument                                                                                                              |
| app.containers[1].env                                  | map | {}                           | any environment variable ENV_NAME: value                                                                                              |
| app.containers[1].env.VALIDATOR_IP                     | string | "127.0.0.1"                  | specifies the `VALIDATOR_IP` environment variable with the validator node IP address                                                  |
| app.containers[1].health.liveness.initialDelaySeconds  | int | 120                          | kubernetes liveness initialDelaySeconds                                                                                               |
| app.containers[1].health.liveness.path                 | string | "/health/liveness"           | kubernetes liveness path                                                                                                              |
| app.containers[1].health.liveness.periodSeconds        | int | 5                            | kubernetes readiness periodSeconds                                                                                                    |
| app.containers[1].health.readiness.initialDelaySeconds | int | 30                           | kubernetes readiness initialDelaySeconds                                                                                              |
| app.containers[1].health.readiness.path                | string | "/health/readiness"          | kubernetes readiness path                                                                                                             |
| app.containers[1].health.readiness.periodSeconds       | int | 1                            | kubernetes readiness periodSeconds                                                                                                    |
| app.containers[1].health.type                          | string | "httpGet"                    | kubernetes healthcheck type                                                                                                           |
| app.containers[1].image.name                           | string | "ghcr.io/my-cloud/ruthenium" | container image                                                                                                                       |
| app.containers[1].image.pullPolicy                     | string | "IfNotPresent"               | container pull policy                                                                                                                 |
| app.containers[1].image.tagOverride                    | string | "latest"                     | container tag override                                                                                                                |
| app.containers[1].name                                 | string | "ui"                         | container name                                                                                                                        |
| app.containers[1].resources.limits.memory              | string | "512Mi"                      | kubernetes resources limits memory                                                                                                    |
| app.containers[1].resources.requests.memory            | string | "128Mi"                      | kubernetes resources requests memory                                                                                                  |
| app.containers[1].secret.ruthenium.PRIVATE_KEY         | string | "privateKey"                 | Gets the privateKey key of the ruthenium secret and populate a container environment variable PRIVATE_KEY with the secret value       |
| app.containers[1].service[0].port                      | string | "80"                         | service port                                                                                                                          |
| app.containers[1].service[0].protocol                  | string | "TCP"                        | service protocol                                                                                                                      |
| app.containers[1].service[0].targetPort                | string | "8080"                       | kubernetes service container listening port                                                                                           |
| app.containers[1].storage.data                         | string | "/data"                      | will mount the `data` volume in the specified path                                                                                    |
| app.replicas                                           | int | 1                            | Number of replicas. it is useless to set over `1` due to Ruthenium fundamental mechanism                                              |
| app.storage.data.accessModes[0]                        | string | "ReadWriteMany"              | storage access mode                                                                                                                   |
| app.storage.data.size                                  | string | "128M"                       | storage size                                                                                                                          |
| app.storage.data.storageClass                          | string | "kube-data"                  | storage class                                                                                                                         |
| app.type                                               | string | "statefulset"                | kubernetes kind (only `statefulset` is available for now)                                                                             |
| app.volumes[0].emptyDir                                | map | {}                           | volume type                                                                                                                           |
| app.volumes[0].name                                    | string | "shared-data"                | volume name                                                                                                                           |
| global.autoReload                                      | string | "true"                       | auto reload set globally (useful in meta deployments)                                                                                 |
| global.image.pullPolicy                                | string | "Always"                     | image pull policies set globally (useful in meta deployments)                                                                         |
| global.registries                                      | list | []                           | list of registries (useful with private registries)                                                                                   |
| secrets[0].annotations                                 | map | {}                           | secret annotation                                                                                                                     |
| secrets[0].data.privateKey                             | string | null                         | secret private key                                                                                                                    |
| secrets[0].name                                        | string | "ruthenium"                  | secret name                                                                                                                           |
| url.domains[0]                                         | string | "ruthenium.example.com"      | domain name on which ruthenium would be available                                                                                     |

