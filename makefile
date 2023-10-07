GOLANG_CI_LINT_VER?=v1.54.2
OUT_BUILD?=./bin/server
IMAGE?=example_replace_it

ifneq (,$(wildcard ./.env))
	include .env
	export
endif

all: lint test

build:
	go build -o ${OUT_BUILD} cmd/server/main.go
.PHONY: build

run:
	go run cmd/server/main.go | tee app.log
.PHONY: run

vendor:
	go mod tidy
	go mod vendor
.PHONY: vendor

test:
	go test -timeout 30s -race ./...
.PHONY: test

lint: bin/golangci-lint
	./bin/golangci-lint run
.PHONY: lint

image:
	docker build -t "${IMAGE}" .
	docker push "${IMAGE}"
.PHONY: image

bin/golangci-lint:
	curl \
		-sSfL \
		https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh \
		| sh -s $(GOLANG_CI_LINT_VER)
