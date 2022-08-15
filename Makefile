IMAGE_NAME := "pfc_demo"
GO_VERSION := $(shell go version | cut -d' ' -f3 | sed "s/go//")
GO_VERSION_MAJOR := $(shell echo $(GO_VERSION) | cut -d' ' -f3 | sed "s/go//" | cut -d'.' -f1)
GO_VERSION_MINOR := $(shell echo $(GO_VERSION) | cut -d' ' -f3 | sed "s/go//" | cut -d'.' -f2,2)

# Ensure that binary is compiled for linux/amd64, even when developing on ARM architecture.
HOST_ARCH := $(shell uname -m)
ifeq ($(HOST_ARCH), arm64)
	DOCKER_PLATFORM := --platform linux/amd64
endif

.PHONY: image
image: ## Creates a service image by compiling the project inside a build-container (which calls `make build` below).
	DOCKER_BUILDKIT=1 docker build $(DOCKER_PLATFORM) -f .build/Dockerfile -t $(IMAGE_NAME) .

.PHONUY: build
build: version_check ## Compiles the project into a binary.
	CGO_ENABLED=0 go build -mod vendor -o go-app cmd/main.go

.PHONY: run
run: image ## Creates service image and runs it.
	docker run --rm $(IMAGE_NAME):latest

# XXX: Could run this in a separate container to avoid version check
.PHONY: vendor
vendor: version_check ## Updates and vendors all dependencies.
	go get -u -t ./...
	go mod tidy
	go mod vendor

.PHONY: test
test: version_check ## Runs tests.
	go test -v ./...

.PHONY: coverage
coverage: version_check ## Runs tests and generates coverage report.
	go test -v -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out

.PHONY: version_check
version_check: ## Checks that the go version is correct.
	echo "$(GO_VERSION_MAJOR)  $(GO_VERSION_MINOR)"
	@if [[ "$(GO_VERSION_MAJOR)" != 1 ]] || [[ "$(GO_VERSION_MINOR)" -lt "18" ]]; then \
		echo "WARNING: go version >=1.18 is required to build this project. Your version is $(GO_VERSION)."; \
		exit 1; \
	fi