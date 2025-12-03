.PHONY: build-amd64 build-arm64 build-all

IMAGE_NAME ?= coutito/es2
TAG ?= latest

build:
	docker buildx build --platform linux/amd64 -t $(IMAGE_NAME):$(TAG)-amd64 .

# Build for AWS
build-aws:
	docker buildx build --platform linux/arm64 -t $(IMAGE_NAME):$(TAG)-arm64 .

build-all:
	docker buildx build --platform linux/amd64,linux/arm64 -t $(IMAGE_NAME):$(TAG) .
