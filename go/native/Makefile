IMG := "localhost/hello-world:v0.1.0"
PORT := "3333"
WHO := "World"
NAMESPACE := "hello-world"

.PHONY: image
image:
	docker buildx build . -t "$(IMG)"

.PHONY: run
run: image
	docker run --rm -it --env HELLO_WORLD_WHO="$(WHO)" hello-world -addr ":$(PORT)"

.PHONY: install
install: image
ifdef KIND
	kind load docker-image --name "$(KIND)" "$(IMG)"
endif
	kubectl create ns $(NAMESPACE) || true
	kubectl apply -f resources -n $(NAMESPACE)