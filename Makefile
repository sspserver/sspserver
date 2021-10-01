include .env
export

BUILD_GOOS ?= linux
BUILD_GOARCH ?= amd64
BUILD_CGO_ENABLED ?= 0

COMMIT_NUMBER ?= $(shell git log -1 --pretty=format:%h)

PROJECT_WORKSPACE := adnet-project

SHELL := /bin/bash -o pipefail
UNAME_OS := $(shell uname -s)
UNAME_ARCH := $(shell uname -m)

TMP_BASE := .tmp
TMP := $(TMP_BASE)/$(UNAME_OS)/$(UNAME_ARCH)
TMP_BIN = $(TMP)/bin
TMP_ETC := $(TMP)/etc
TMP_LIB := $(TMP)/lib
TMP_VERSIONS := $(TMP)/versions
TMP_FOSSA_GOPATH := $(TMP)/fossa/go

APP_TAGS := "nats allplatform"

unexport GOPATH
export GOPATH=$(abspath $(TMP))
export GO111MODULE := on
export GOBIN := $(abspath $(TMP_BIN))
export PATH := $(GOBIN):$(PATH)
export GOSUMDB := off
export GOFLAGS=-mod=mod
# Go 1.13 defaults to TLS 1.3 and requires an opt-out.  Opting out for now until certs can be regenerated before 1.14
# https://golang.org/doc/go1.12#tls_1_3
export GODEBUG := tls13=0

CONTAINER_IMAGE ?= rtb/sspserver

DOCKER_COMPOSE := docker-compose -p $(PROJECT_WORKSPACE) -f deploy/develop/docker-compose.yml
DOCKER_BUILDKIT := 1

GOLANGLINTCI_VERSION := latest
GOLANGLINTCI := $(TMP_VERSIONS)/golangci-lint/$(GOLANGLINTCI_VERSION)
$(GOLANGLINTCI):
	$(eval GOLANGLINTCI_TMP := $(shell mktemp -d))
	cd $(GOLANGLINTCI_TMP); go get github.com/golangci/golangci-lint/cmd/golangci-lint@$(GOLANGLINTCI_VERSION)
	@rm -rf $(GOLANGLINTCI_TMP)
	@rm -rf $(dir $(GOLANGLINTCI))
	@mkdir -p $(dir $(GOLANGLINTCI))
	@touch $(GOLANGLINTCI)


GOMOCK_VERSION := v1.4.4
GOMOCK := $(TMP_VERSIONS)/mockgen/$(GOMOCK_VERSION)
$(GOMOCK):
	$(eval GOMOCK_TMP := $(shell mktemp -d))
	cd $(GOMOCK_TMP); go get github.com/golang/mock/mockgen@$(GOMOCK_VERSION)
	@rm -rf $(GOMOCK_TMP)
	@rm -rf $(dir $(GOMOCK))
	@mkdir -p $(dir $(GOMOCK))
	@touch $(GOMOCK)


QTC_VERSION := latest
QTC := $(TMP_VERSIONS)/qtc/$(QTC_VERSION)
$(QTC):
	$(eval QTC_TMP := $(shell mktemp -d))
	cd $(QTC_TMP); go get github.com/valyala/quicktemplate/qtc@$(QTC_VERSION)
	@rm -rf $(QTC_TMP)
	@rm -rf $(dir $(QTC))
	@mkdir -p $(dir $(QTC))
	@touch $(QTC)


.PHONY: deps
deps: $(GOLANGLINTCI) $(GOMOCK) $(QTC)

.PHONY: all
all: lint cover

.PHONY: lint
lint: golint

.PHONY: golint
golint: $(GOLANGLINTCI)
	# golint -set_exit_status ./...
	golangci-lint run -v ./...

.PHONY: fmt
fmt: ## Run formatting code
	@echo "Fix formatting"
	@gofmt -w ${GO_FMT_FLAGS} $$(go list -f "{{ .Dir }}" ./...); if [ "$${errors}" != "" ]; then echo "$${errors}"; fi

.PHONY: test
test: ## Run unit tests
	go test -v -tags ${APP_TAGS} -race ./...

.PHONY: qtc
qtc: ## Build templates
	qtc -dir=private/templates

.PHONY: tidy
tidy:
	go mod tidy

.PHONY: cover
cover:
	@mkdir -p $(TMP_ETC)
	@rm -f $(TMP_ETC)/coverage.txt $(TMP_ETC)/coverage.html
	go test -race -coverprofile=$(TMP_ETC)/coverage.txt -coverpkg=./... ./...
	@go tool cover -html=$(TMP_ETC)/coverage.txt -o $(TMP_ETC)/coverage.html
	@echo
	@go tool cover -func=$(TMP_ETC)/coverage.txt | grep total
	@echo
	@echo Open the coverage report:
	@echo open $(TMP_ETC)/coverage.html

.PHONY: generate-code
generate-code: ## Run codegeneration procedure
	@echo "Generate code"
	@go generate ./...

.PHONY: build
build: ## Build application
	@echo "Build application"
	@rm -rf .build/sspserver
	GOOS=${BUILD_GOOS} GOARCH=${BUILD_GOARCH} CGO_ENABLED=${BUILD_CGO_ENABLED} \
		go build -ldflags "-X main.buildDate=`date -u +%Y%m%d.%H%M%S` -X main.buildCommit=${COMMIT_NUMBER}" \
			-tags ${APP_TAGS} -o ".build/sspserver" cmd/sspserver/main.go

.PHONY: docker-dev-build
docker-dev-build: build
	echo "Build develop docker image"
	DOCKER_BUILDKIT=${DOCKER_BUILDKIT} docker build -t ${CONTAINER_IMAGE} -f deploy/develop/Dockerfile .

.PHONY: run
run: docker-dev-build ## Run service by docker-compose
	@echo "Run sspserver service"
	$(DOCKER_COMPOSE) up sspserver

.PHONY: stop
stop: ## Stop all services
	@echo "Stop all services"
	$(DOCKER_COMPOSE) stop

.PHONY: help
help:
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

.DEFAULT_GOAL := help
