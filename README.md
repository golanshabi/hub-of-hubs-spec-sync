[comment]: # ( Copyright Contributors to the Open Cluster Management project )

# Hub-of-Hubs Spec Syncer
Red Hat Advanced Cluster Management Hub-of-Hubs Spec Syncer

(Syncer - the name taken from the [KCP project](https://github.com/kcp-dev/kcp/blob/main/contrib/demo/README.md#syncer))

## How it works

## Build to run locally

```
make build
```

## Run Locally

Disable the currently running controller in the cluster (if previously deployed):

```
kubectl scale deployment hub-of-hubs-spec-syncer --kubeconfig $TOP_HUB_CONFIG -n open-cluster-management --replicas 0
```

Set the following environment variables:

* DATABASE_URL
* WATCH_NAMESPACE

Set the `DATABASE_URL` according to the PostgreSQL URL format: `postgres://YourUserName:YourURLEscapedPassword@YourHostname:5432/YourDatabaseName?sslmode=verify-full`.

:exclamation: Remember to URL-escape the password, you can do it in bash:

```
python -c "import sys, urllib as ul; print ul.quote_plus(sys.argv[1])" 'YourPassword'
```

`WATCH_NAMESPACE` can be defined empty so the controller will watch all the namespaces.

```
./build/_output/bin/hub-of-hubs-spec-syncer --kubeconfig $TOP_HUB_CONFIG
```

## Build image

Define the `REGISTRY` environment variable.

```
make build-images
```

## Deploy to a cluster

1.  Create a secret with your database url:

    ```
    kubectl create secret generic hub-of-hubs-database-secret --kubeconfig $TOP_HUB_CONFIG -n open-cluster-management --from-literal=url=$DATABASE_URL
    ```

1.  Deploy the operator:

    ```
    IMAGE_TAG=latest envsubst < deploy/operator.yaml.template | kubectl apply --kubeconfig $TOP_HUB_CONFIG -n open-cluster-management -f -
    ```
