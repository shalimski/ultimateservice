run:
	go run main.go

build:
	go build -ldflags "-X main.build=local"

VERSION := 1.0

all: service

service:
	docker build \
		-f infra/docker/Dockerfile \
		-t sevice-amd64:${VERSION} \
		--build-arg BUILD_REF=${VERSION} \
		.