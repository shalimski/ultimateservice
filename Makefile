run:
	go run main.go

build:
	go build -ldflags "-X main.build=local"

VERSION := 1.0
KIND_CLUSTER := starter-cluster  

all: service

service:
	docker build \
		-f infra/docker/Dockerfile \
		-t service-amd64:${VERSION} \
		--build-arg BUILD_REF=${VERSION} \
		.

kind-up:
	kind create cluster \
		--image kindest/node:v1.24.3 \
		--name $(KIND_CLUSTER) \
		--config infra/k8s/kind/kind-config.yaml
	kubectl config set-context --current --namespace=service-system

kind-down:
	kind delete cluster --name $(KIND_CLUSTER)

kind-load:
	kind load docker-image service-amd64:$(VERSION) --name $(KIND_CLUSTER)

kind-apply:
	kustomize build infra/k8s/kind/service-pod | kubectl apply -f -

kind-status:
	kubectl get nodes -o wide
	kubectl get svc -o wide
	kubectl get pods -o wide --watch

kind-status-all:
	kubectl get nodes -o wide
	kubectl get svc -o wide
	kubectl get pods -o wide --watch --all-namespaces

kind-logs:
	kubectl logs -l app=service --all-containers=true -f --tail=100

kind-restart:
	kubectl rollout restart deployment service-pod

kind-update: all kind-load kind-restart

kind-update-apply: all kind-load kind-apply

kind-describe:
	kubectl describe pod -l app=service