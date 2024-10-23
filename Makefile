include .env
export

DEF_APP_TAGS := postgres,nats,redisps,allplatform,fsloader,dbloader,htmltemplates,jaeger,migrate
APP_TAGS  ?= $(or ${APP_BUILD_TAGS},${DEF_APP_TAGS})

include deploy/build.mk

PROJECT_WORKSPACE := adnet-project
PROJECT_NAME ?= sspserver
DOCKER_COMPOSE := docker compose -p $(PROJECT_WORKSPACE) -f deploy/develop/docker-compose.yml
DOCKER_CONTAINER_IMAGE := ${PROJECT_WORKSPACE}/${PROJECT_NAME}
DOCKER_CONTAINER_MUGRATE_IMAGE := ${DOCKER_CONTAINER_IMAGE}:migrate-latest
DOCKER_EVENTSTREAM_CONTAINER_IMAGE := ${PROJECT_WORKSPACE}/eventstream

.PHONY: all
all: lint cover

.PHONY: lint
lint: golint

.PHONY: golint
golint:
	# golint -set_exit_status ./...
	golangci-lint run -v ./...

.PHONY: fmt
fmt: ## Run formatting code
	@echo "Fix formatting"
	@gofmt -w ${GO_FMT_FLAGS} $$(go list -f "{{ .Dir }}" ./...); if [ "$${errors}" != "" ]; then echo "$${errors}"; fi

.PHONY: test
test: ## Run unit tests
	go test -v -tags "${APP_TAGS}" -race ./...

.PHONY: qtc
qtc: ## Build templates
	go run github.com/valyala/quicktemplate/qtc -dir=private/templates

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
	@rm -rf .build
	@$(call do_build,"cmd/sspserver/main.go",sspserver)

.PHONY: build-docker-dev
build-docker-dev: build ## Build docker image for development
	echo "Build develop docker image"
	DOCKER_BUILDKIT=${DOCKER_BUILDKIT} docker build -t ${DOCKER_CONTAINER_IMAGE} -f deploy/develop/Dockerfile .

.PHONY: run
run: build-docker-dev ## Run service by docker-compose
	@echo "Run sspserver service on http://localhost:${DOCKER_SERVER_HTTP_PORT}"
	$(DOCKER_COMPOSE) up sspserver

.PHONY: stop
stop: ## Stop all services
	@echo "Stop all services"
	$(DOCKER_COMPOSE) stop

.PHONY: dbcli
dbcli: ## Open development database
	$(DOCKER_COMPOSE) exec $(DOCKER_DATABASE_NAME) psql -U $(DATABASE_USER) $(DATABASE_DB)

.PHONY: chi
chi: ## Run clickhouse client
	${K8C} -n adlab-statistic exec -it chi-statistic-statistic-0-0-0 -- clickhouse-client

.PHONY: help
help:
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' Makefile | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

.DEFAULT_GOAL := help
