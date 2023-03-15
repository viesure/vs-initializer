# vs-initializer

k8s init container image for providing secrets from GCP Secret Manager inside config files. This is the go version.

## Purposes

vs-initializer is an initcontainer image for
[GKE](https://console.cloud.google.com/kubernetes)
that:
- parses all files inside a given directory and pulls secrets from
[GCP Secret Manager](https://console.cloud.google.com/security/secret-manager)
and puts the secret's value into files to enable easy access for the application container. This is achieved by mounting the
same volume to the initcontainer and the application container.
- injects environment variables into the actual app container - either from static values coming from a ConfigMap or
secret values obtained from
[GCP Secret Manager](https://console.cloud.google.com/security/secret-manager)

## Image location & version
The resulting image can be pulled from the
[public registry on GCP project vs-tools](https://console.cloud.google.com/gcr/images/vs-tools/global/toolimages/vs-initializer?project=vs-tools):

```gcr.io/vs-tools/toolimages/vs-initializer:latest```

## Configuration

The initcontainer image is configured using environment variables:

| environment variable | default value | description                                                                    |
|----------------------|---------------|--------------------------------------------------------------------------------|
| APP_LOG_LEVEL        | INFO          | Log level of app, one of:<br/>DEBUG, INFO, WARN, ERROR                         |
| TEMPLATE_DIR         | /data.tmpl    | Directory where the file(s) are read from                                      |
| OUTPUT_DIR           | /data         | Directory where the final files should be saved (for use inside app container) |
| ENV_SECRET           | app-env       | Name of Secret that holds the environment variables for app container          |
| APP_NAMESPACE        | *             | Namespace where to store environment ConfigMap for app container               |
* If APP_NAMESPACE is not defined, the namespace will be used from the k8s deployment (via the file `/var/run/secrets/kubernetes.io/serviceaccount/namespace`)

Please find an example for a deployment in [/examples/deployment.yml](examples/deployment.yml#L89-100)

### Secret Manager URLs

The URLs for accessing certain secrets in [GCP Secret Manager](https://console.cloud.google.com/security/secret-manager)
are constructed as follows:

`sm://<gcp-project>/<secret-name>[/<version>]?<labelname>=<labelvalue>&<labelname>=<labelvalue>&...`

| parameter       | required | description                                                                                          | 
|-----------------|----------|------------------------------------------------------------------------------------------------------|
| `<gcp-project>` | required | The project name in GCP that holds the secrets                                                       |
| `<secret-name>` | required | The value of the label "secret-name"                                                                 |
| `<version>`     | optional | The version of the secret value.<br/>If not specified, the latest version of the secret is selected. |
| `<labelname>`   | optional | The name of the label                                                                                |
| `<labelvalue>`  | optional | The value of the label                                                                               |

NOTE:
Each secret stored in [GCP Secret Manager](https://console.cloud.google.com/security/secret-manager)
must have exactly the same amount of labels that are used to query it using the URL plus one label called
"secret-name"

### environment file ".env"

Inside the template ConfigMap (or the template directory), a key/file `.env` can be defined that will be converted into
a separate Secret that can be used to configure environment variables for the actual app container.
Inside this `.env` key also [Secret Manager URLs](#secret-manager-urls) can be used, the init conmtainer will replace it with the actual
secret value from Secret Manager:

```yaml
...
apiVersion: v1
kind: Secret
...
data:
  ...
  .env: |
    ENVVAR1=example-value-1
    ENVVAR2=example-value-2
    DBSECRET=<Secret Manager URL>
...
```

An _empty_ Secret (without data) must be included, when the app container needs to be configured using
environment variables. Given the default name for this Secret ("app-env"), this Secret
is then configured as the/one source of the environment variables for the app container:

```yaml
...
apiVersion: apps/v1
kind: Deployment
...
spec:
  template:
    spec:
      initContainers:
        - name: app-init
        ...
      containers:
        - name: app
          image: app-image:1.0
          envFrom:
            - secretRef:
                name: app-env
```

When the app container is started, the values from the `.env` key in the Secret
will be available as environment variables

Find a complete example inside [/examples/deployment.yml](examples/deployment.yml)


## GKE Setup

### Workload Identity
The GKE Cluster running the deployment has to have [Workload Identity](https://cloud.google.com/kubernetes-engine/docs/concepts/workload-identity) enabled to be able
to access GCP resources and services - like the Secret Manager in this case.

### Service Account
A Kubernetes Service Account must be created for running the deployment, that has to
be mapped to an IAM Service Account (in GCP) that has the following IAM roles in the project
where the secrets are stored:

- Secret Manager Secret Accessor
- Secret Manager Viewer

A proper configuration of the Kubernetes Service Account with the necessary permissions
(Role, RoleMapping) is shown inside [/examples/deployment.yml](examples/deployment.yml)
