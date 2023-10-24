#!/bin/bash

NAMESPACE=${1:-"hello-world"}
VERSION=${2:-"main"}
CONTAINER_REGISTRY=${CONTAINER_REGISTRY:-"localhost"}

docker buildx build ./backend --file Dockerfile --tag "$CONTAINER_REGISTRY/hello-world-backend:$VERSION"
docker buildx build ./frontend --file Dockerfile --tag "$CONTAINER_REGISTRY/hello-world-frontend:$VERSION"

if [ "$CONTAINER_REGISTRY" != "localhost" ]; then
    docker push "$CONTAINER_REGISTRY/hello-world-backend:$VERSION"
    docker push "$CONTAINER_REGISTRY/hello-world-frontend:$VERSION"
fi
if [ "$KIND_CLUSTER" != "" ]; then
    kind load docker-image --name "$KIND_CLUSTER" "$CONTAINER_REGISTRY/hello-world-backend:$VERSION"
    kind load docker-image --name "$KIND_CLUSTER" "$CONTAINER_REGISTRY/hello-world-frontend:$VERSION"
fi

kubectl create namespace $NAMESPACE 2>/dev/null || true
kubectl apply --namespace $NAMESPACE --filename hack/
