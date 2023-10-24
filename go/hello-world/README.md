# Hello World

In this example we'll compare two simple `hello-world` apps that read the
environment variable `HELLO_WORLD_WHO` and say hello to them. The first is a
"native" Kubernetes app using HTTP servers, Dockerfile, and various Kubernetes
resources. The second is a KubeFox component which will be deployed using
[fox](https://github.com/xigxog/fox), the KubeFox CLI.

To get started you'll need the following installed:

- [Go](https://go.dev/doc/install)
- [Git](https://github.com/git-guides/install-git)
- [Docker](https://docs.docker.com/engine/install/)
- [kubectl](https://kubernetes.io/docs/tasks/tools/)

Also you'll need access to a Kubernetes cluster. If you'd like to
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

First, build the app container images using Docker. If you are using kind
locally you can leave the container registry set to localhost, otherwise replace
it with the container registry you'd like to use.

```shell
export CONTAINER_REGISTRY="localhost"
docker buildx build ./backend --file Dockerfile --tag "$CONTAINER_REGISTRY/hello-world-backend:main"
docker buildx build ./frontend --file Dockerfile --tag "$CONTAINER_REGISTRY/hello-world-frontend:main"
```

Next you'll need to make the container images available to Kubernetes. If you
are using kind you can load the image directly without using a container
registry.

```shell
export KIND_CLUSTER="kind"
kind load docker-image --name "$KIND_CLUSTER" "$CONTAINER_REGISTRY/hello-world-backend:main"
kind load docker-image --name "$KIND_CLUSTER" "$CONTAINER_REGISTRY/hello-world-frontend:main"
```

Otherwise push the image to the container registry.

```shell
docker push "$CONTAINER_REGISTRY/hello-world-backend:main"
docker push "$CONTAINER_REGISTRY/hello-world-frontend:main"
```

Finally, create a Kubernetes Namespace and apply the ConfigMaps and Deployment
to run the app on Kubernetes.

```shell
kubectl create namespace hello-world
kubectl apply --namespace hello-world --filename hack/environments/world.yaml
kubectl apply --namespace hello-world --filename hack/deployments/

# Example output:
# namespace/hello-world created
# configmap/env created
# deployment.apps/hello-world-backend created
# service/hello-world-backend created
# deployment.apps/hello-world-frontend created
# service/hello-world-frontend created
```

If everything worked you should see the Pod running. You can check using
kubectl.

```shell
kubectl get pods --namespace hello-world

# Example output:
# NAME                                    READY   STATUS    RESTARTS   AGE
# hello-world-backend-865d6697d5-2vwnw    1/1     Running   0          10s
# hello-world-frontend-5579b569c9-fdsnw   1/1     Running   0          19s
```

Time to test the app. To keep things simple you'll port forward to the Pod to
access its HTTP server. Open up a new terminal and run the following to start
the port forward. This will open the port `8080` on your workstation which will
forward requests to the Pod.

```shell
kubectl port-forward --namespace hello-world service/hello-world 8080:http

# Example output:
# Forwarding from 127.0.0.1:8080 -> 3333
# Forwarding from [::1]:8080 -> 3333
```

Finally send a HTTP request to the app.

```shell
curl http://127.0.0.1:8080/hello

# Example output:
#👋 Hello World!
```

It works! But how do you run the app in a different environment so you can
change who to say hello to? You need to update the ConfigMap `env` that contains
the `HELLO_WORLD_WHO` variable. Of course if you change what is running now the
`world` environment will no longer exist. Instead you can create a new Namespace
and run the app there with the updated ConfigMap. Try it out.

```shell
kubectl create namespace hello-universe
kubectl apply --namespace hello-universe --filename hack/environments/universe.yaml
kubectl apply --namespace hello-universe --filename hack/deployments/

# Example output:
# namespace/hello-universe created
# configmap/env created
# deployment.apps/hello-world-backend created
# service/hello-world-backend created
# deployment.apps/hello-world-frontend created
# service/hello-world-frontend created
```

Now you can test the app in the new environment. Once again open up a new
terminal and run the following to start the port forward but use port `8081`
this time.

```shell
kubectl port-forward --namespace hello-universe service/hello-world 8081:http

# Example output:
# Forwarding from 127.0.0.1:8081 -> 3333
# Forwarding from [::1]:8081 -> 3333
```

Then send a HTTP request to app.

```shell
curl http://127.0.0.1:8081/hello

# Example output:
#👋 Hello Universe!
```

Great! It's using the new environment. Take a look at what is running on
Kubernetes now. You can use a label from the Deployments to show Pods from
multiple namespaces.

```shell
kubectl get pods --all-namespaces --selector=app.kubernetes.io/name=hello-world-native

# Example output:
# NAMESPACE        NAME                                            READY   STATUS    RESTARTS   AGE
# hello-universe   hello-world-backend-9f67b958d-lwm6t             1/1     Running   0          2m30s
# hello-universe   hello-world-frontend-887674586-2q298            1/1     Running   0          2m25s
# hello-world      hello-world-backend-865d6697d5-tpbfr            1/1     Running   0          3m49s
# hello-world      hello-world-frontend-5579b569c9-fdsnw           1/1     Running   0          5m10s
```

## KubeFox

If this is your first time using KubeFox you'll need to install and setup
[fox](https://github.com/xigxog/fox) CLI tool. The easiest way is to use `go
install`, make sure that `$GOPATH/bin` directory is on your path. Or you can
download it from the [release page](https://github.com/xigxog/fox/releases). The
`config setup` command will guide you through the setup process

```shell
# This enables some extra output.
export FOX_INFO=true
export PATH=$PATH:$GOPATH/bin
go install github.com/xigxog/fox@latest
fox config setup
```

You'll need to run the following commands from the `kubefox` directory.

```shell
cd ../kubefox/
```

First, apply the KubeFox Environment resources. Environments are similar to
ConfigMaps but are cluster scoped so they can be used by multiple Namespaces as
the same time.

```shell
kubectl apply --filename hack/environments/

# Example output:
# environment.kubefox.xigxog.io/universe created
# environment.kubefox.xigxog.io/world created
```

Now you can `publish` the app using `fox`. This will build the container images,
push them to the registry, and deploy the app to the KubeFox platform running on
your Kubernetes cluster.

```shell
fox publish deployment-a --wait 5m
```

Now you can test the KubeFox app. To keep things simple you'll again port
forward, but this time you'll connect to the KubeFox broker with some help from
`fox`. Open up a new terminal and run the following to start the port forward.
This will open the port `8082` on your workstation which will forward requests
to the KubeFox platform.

```shell
fox proxy 8082
```

When KubeFox deploys an app it starts the components but will not automatically
send requests to it until it is released. But you can still test deployments by
providing some context. KubeFox needs two pieces of information, the deployment
to use and the environment to inject. These can be passed as headers or query
parameters.

```shell
curl "http://127.0.0.1:8082/hello?kf-dep=deployment-a&kf-env=world"

# Example output:
#👋 Hello World!
```

Next try switching to the `universe` environment created earlier. With KubeFox
there is no need to create another deployment to switch environments, simply
change the query parameter!

```shell
curl "http://127.0.0.1:8082/hello?kf-dep=deployment-a&kf-env=universe"

# Example output:
#👋 Hello Universe!
```

Now let's release the app so we don't have to specify all those details in the
request. It is important to understand that Fox works against the active state
of the Git repo. To deploy or release a different version of your app simply
checkout the tag, branch, or commit you want and let Fox do the rest.

```shell
fox release dev --env world --wait 5m
```

Try the same request from above, but this time don't specify the context. Since
the app has been released the request will get matched by the component's route
and the environment will be automatically injected by KubeFox.

```shell
curl "http://localhost:8082/hello"
```

Take a look at what is running on Kubernetes to support the KubeFox app.

```shell
kubectl get pods --all-namespaces --selector=app.kubernetes.io/name=hello-world-kubefox

# Example output:
# NAMESPACE     NAME                                                    READY   STATUS    RESTARTS   AGE
# kubefox-dev   hello-world-kubefox-backend-1403140-56d896767d-2r88l    1/1     Running   0          49s
# kubefox-dev   hello-world-kubefox-frontend-1403140-5bb7bd679b-7wp27   1/1     Running   0          48s
```

Notice that even though we have made a deployment, a release, and have two
environments there are still only two Pods running! This is possible because
KubeFox injects context at request time instead of deploy time. Adding
environments has nearly no overhead!
