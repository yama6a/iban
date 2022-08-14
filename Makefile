IMAGE_NAME := "pfc_demo"
GO_VERSION := $(shell go version | cut -d' ' -f3 | sed "s/go//")

# Ensure that binary is compiled for linux/amd64, even when developing on ARM architecture.
HOST_ARCH := $(shell uname -m)
ifeq ($(HOST_ARCH), arm64)
	DOCKER_PLATFORM := --platform linux/amd64
endif

.PHONY: image
image: ## Creates a service image by compiling the project inside a build-container (which calls `make build` below).
	DOCKER_BUILDKIT=1 docker build $(DOCKER_PLATFORM) -f .build/Dockerfile -t $(IMAGE_NAME) .

.PHONUY: build
build: ## Compiles the project into a binary.
	# Version-validation needs refinement. 1.2* also allows for v1.2 (which is obviously not supported).
	# Try wrestling with Bash+regex to make this more robust one day:
	# 	e.g. [[ "$(GO_VERSION)" =~ [0-9] (1[8-9)|(2[0-9]) ]]  or something like that.
	@if [[ "$(GO_VERSION)" != 1.18* ]] && [[ "$(GO_VERSION)" != 1.19* ]] && [[ "$(GO_VERSION)" != 1.2* ]]; then \
		echo "WARNING: go version >=1.18 is required to build this project. Your version is $(GO_VERSION)."; \
		exit 1; \
	fi
	CGO_ENABLED=0 go build -mod vendor -o go-app cmd/main.go

.PHONY: run
run: image
	docker run --rm $(IMAGE_NAME):latest

