# For testing load on the service.
# go install github.com/rakyll/hey@latest
# hey -m GET -c 100 -n 10000 http://localhost:9020/
#
# Access metrics (9021)
# go install github.com/divan/expvarmon@latest
# expvarmon -ports=":9021" -endpoint="/debug/vars" -vars="build,requests,goroutines,errors,panics,mem:memstats.Alloc"
#
# To generate a private/public key PEM file.
# openssl genpkey -algorithm RSA -out private.pem -pkeyopt rsa_keygen_bits:2048
# openssl rsa -pubout -in private.pem -out public.pem
# or use make generatekey

run:
	go run app/services/sales-api/main.go

build:
	go build -ldflags "-X main.build=local"

lint:
	golangci-lint run ./...

generatekeys:
	go run app/tooling/admin/main.go -bits 4096

VERSION := 1.0
KIND_CLUSTER := starter-cluster  

all: sales-api

sales-api:
	docker build \
		-f infra/docker/Dockerfile.sales-api \
		-t sales-api-amd64:${VERSION} \
		--build-arg BUILD_REF=${VERSION} \
		.

kind-up:
	kind create cluster \
		--image kindest/node:v1.24.3 \
		--name $(KIND_CLUSTER) \
		--config infra/k8s/kind/kind-config.yaml
	kubectl config set-context --current --namespace=sales-api-system

kind-down:
	kind delete cluster --name $(KIND_CLUSTER)

kind-load:
	cd infra/k8s/kind/sales-api-pod; kustomize edit set image sales-api-image=sales-api-amd64:$(VERSION)
	kind load docker-image sales-api-amd64:$(VERSION) --name $(KIND_CLUSTER)

kind-apply:
	kustomize build infra/k8s/kind/sales-api-pod | kubectl apply -f -

kind-status:
	kubectl get nodes -o wide
	kubectl get svc -o wide
	kubectl get pods -o wide --watch

kind-status-all:
	kubectl get nodes -o wide
	kubectl get svc -o wide
	kubectl get pods -o wide --watch --all-namespaces

kind-logs:
	kubectl logs -l app=sales-api --all-containers=true -f --tail=100

kind-restart:
	kubectl rollout restart deployment sales-api-pod

kind-update: all kind-load kind-restart

kind-update-apply: all kind-load kind-apply

kind-describe:
	kubectl describe pod -l app=sales-api