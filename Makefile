.PHONY: install test builder-image push lint fetch-dependencies binary-image

DOCKER_IMAGE ?= stakater/whitelister

# Default value "dev"
DOCKER_TAG ?= dev
REPOSITORY = ${DOCKER_IMAGE}:${DOCKER_TAG}
BUILDER ?= whitelister-builder
BINARY ?= whitelister

install:  fetch-dependencies

fetch-dependencies:
	dep ensure -v

test:
	go test -v ./...

builder-image:
	@docker build --network host -t "${BUILDER}" -f build/package/Dockerfile.build .

binary-image: builder-image
	@docker run --network host --rm "${BUILDER}" | docker build --network host -t "${REPOSITORY}" -f Dockerfile.run -

lint:
	golangci-lint run --enable-all --skip-dirs vendor

push:
	docker push ${DOCKER_IMAGE}:${DOCKER_TAG}
