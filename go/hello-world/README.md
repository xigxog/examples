# Hello World

In this example we'll compare two simple `hello-world` apps. The first is a
"native" Kubernetes app using an HTTP server, Dockerfile, and various Kubernetes
resources. The second is a KubeFox component which will be deployed using
[fox](https://github.com/xigxog/fox), the KubeFox CLI.

To get started you'll need [Go](https://go.dev/doc/install),
[Git](https://github.com/git-guides/install-git),
[Docker](https://docs.docker.com/engine/install/), and
[kubectl](https://kubernetes.io/docs/tasks/tools/) installed on your
workstation. Also you'll need access to a Kubernetes cluster. If you'd like to
run a Kubernetes cluster on your workstation for testing we recommend using
[kind (Kubernetes in Docker)](https://kind.sigs.k8s.io/docs/user/quick-start/).

## Native

We have included everything needed to run the app in Kubernetes. This includes:

- Dockerfile to build and package the app into an OCI container.
- Kubernetes Deployment to run the app on the Kubernetes cluster.
- Kubernetes ConfigMap to store environment variables used by the app.
- Kubernetes Service to be able to send requests to the Pod.

Take a look at the various files, there is a lot going on.

There are a few steps to deploy the app to Kubernetes. Run all the following
commands from the `native` directory.

First, build the app container image using Docker. If you are using kind locally
you can leave the container registry set to localhost, otherwise replace it with
the container registry you'd like to use.

```shell
export CONTAINER_REGISTRY="localhost"
docker buildx build . -t "$CONTAINER_REGISTRY/hello-world:v0.1.0"
```

Next you'll need to make the container image available to Kubernetes. If you are
using kind you can load the image directly without using a container registry.

```shell
export KIND_CLUSTER="kind"
kind load docker-image --name "$KIND_CLUSTER" "$CONTAINER_REGISTRY/hello-world:v0.1.0"
```

Otherwise push the image to the container registry.

```shell
docker push "$CONTAINER_REGISTRY/hello-world:v0.1.0"
```

Finally, create a Kubernetes Namespace and apply the ConigMaps and Deployment to
run the app on Kubernetes.

```shell
kubectl create namespace hello-world
kubectl apply --namespace hello-world --filename resources/
```

If everything worked you should see the Pod running. You can check using
kubectl.

```shell
kubectl get pods --namespace hello-world

# Example output:
# NAME                         READY   STATUS    RESTARTS      AGE
# hello-world-7fcdb5bd-79hxk   1/1     Running   0             21s
```

Now let's test the app. To keep things simple you'll port forward to the Pod to
access its HTTP server. Open up a new terminal and run the following to start
the port forward. This will open the port `8080` on your workstation which will
forward all traffic to the Pod.

```shell
kubectl port-forward --namespace hello-world service/hello-world 8080:http

# Example output:
# Forwarding from 127.0.0.1:8080 -> 3333
# Forwarding from [::1]:8080 -> 3333
```

Finally send a HTTP request to the app.

```shell
curl http://127.0.0.1:8080/examples/hello-world

# Example output:
#👋 Hello World!
```

Phew, it works!

> TODO: deploy to `universe` environment

## KubeFox

If this is your first time using KubeFox you'll need to install and setup
[fox](https://github.com/xigxog/fox) CLI tool. The easiest way is to use `go
install`, just make sure that `$GOPATH/bin` directory is on your path. Or you
can download it from the [release page](https://github.com/xigxog/fox/releases).
The `config setup` command will guide you through the setup process

```shell
export PATH=$PATH:$GOPATH/bin
# This enables some extra output.
export FOX_INFO=true
go install github.com/xigxog/fox@latest
fox config setup
```

You'll need to run the following commands from the `kubefox` directory.

```shell
fox publish my-deployment --wait 5m
```

Now let's test the KubeFox app. To keep things simple you'll again port forward,
but this time you'll connect to the KubeFox broker with some help from `fox`.
Open up a new terminal and run the following to start the port forward. This
will open the port `8081` on your workstation which will forward all traffic to
the broker.

```shell
fox proxy 8081
```

```shell
curl http://127.0.0.1:8081/examples/hello-world?kf-dep=my-deployment&kf-env=env-world

# Example output:
#👋 Hello World!
```

```shell
curl http://127.0.0.1:8081/examples/hello-world?kf-dep=my-deployment&kf-env=env-universe

# Example output:
#👋 Hello Universe!
```
