#!/bin/bash

set -e  # Exit immediately if a command exits with a non-zero status
set -x  # Print commands and their arguments as they are executed

# Get the Go binary path
GOBIN=$(go env GOPATH)/bin

# Install protoc plugins if not already installed
go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest

# Add GOBIN to PATH
export PATH=$PATH:$GOBIN

# Check if protoc-gen-go and protoc-gen-go-grpc are installed
echo "Looking for protoc plugins in $GOBIN"
if [ -f "$GOBIN/protoc-gen-go" ]; then
    echo "Found protoc-gen-go"
else
    echo "protoc-gen-go not found in $GOBIN"
    ls -la $GOBIN
fi

# Generate Go code from proto files
echo "Generating Go code from proto files..."
protoc --proto_path=. \
       --go_out=.. --go_opt=paths=source_relative \
       --go-grpc_out=.. --go-grpc_opt=paths=source_relative \
       github/v1/github.proto github/v1/luminex.proto

echo "Proto generation completed." 