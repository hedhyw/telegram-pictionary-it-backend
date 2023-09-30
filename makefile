GOLANG_CI_LINT_VER:=v1.54.2

include .env
export

all: lint test

run:
	go run cmd/server/main.go  | tee app.log
.PHONY: run

vendor:
	go mod tidy
	go mod vendor
.PHONY: vendor

test:
	go test -race ./...
.PHONY: test

lint: bin/golangci-lint
	./bin/golangci-lint run
.PHONY: lint

bin/golangci-lint:
	curl \
		-sSfL \
		https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh \
		| sh -s $(GOLANG_CI_LINT_VER)
