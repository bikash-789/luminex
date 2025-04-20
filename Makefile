GOHOSTOS:=$(shell go env GOHOSTOS)
GOPATH:=$(shell go env GOPATH)
VERSION=$(shell git describe --tags --always)
GOVERSION := $(shell go version | cut -d " " -f 3 | cut -c 3-)
GOROOT:=$(shell go env GOROOT)
PROJECT_DIR = $(shell pwd)
PROJECT_BIN = $(PROJECT_DIR)/bin
GOLANGCI_LINT = $(PROJECT_BIN)/golangci-lint


ifeq ($(GOHOSTOS), windows)
	INTERNAL_PROTO_FILES=$(shell $(Git_Bash) -c "find internal -name '*.proto'")
	API_PROTO_FILES=$(shell $(Git_Bash) -c "find api -name '*.proto'")
else
	INTERNAL_PROTO_FILES=$(shell find internal -name '*.proto')
	API_PROTO_FILES=$(shell find api -name '*.proto')
endif


.PHONY: init
init:
	@echo ''
	@echo 'Init:'
	@echo ''
	@echo 'Installing dependencies:'
	go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
	go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
	go install github.com/go-kratos/kratos/cmd/protoc-gen-go-http/v2@latest
	go install github.com/go-kratos/kratos/cmd/kratos/v2@latest
	go install github.com/google/wire/cmd/wire@latest

lint:
	go install github.com/golangci/golangci-lint/cmd/golang-lint@v1.46.2
	golangci-lint run


.PHONY: config
config:
	protoc --proto_path=./internal \
	       --proto_path=./third_party \
 	       --go_out=paths=source_relative:./internal \
	       $(INTERNAL_PROTO_FILES)

.PHONY: generate
generate:
	go mod tidy
	go get github.com/google/wire/cmd/wire@latest
	go generate ./...


.PHONY: build
# build
build:
	mkdir -p bin/ && go build -ldflags "-X main.Version=$(VERSION)" -o ./bin/ ./cmd/server

.PHONY: all
# generate all
all:
	make config;
	make generate;

.PHONY: clean
# clean
clean:
	rm -rf bin/ 