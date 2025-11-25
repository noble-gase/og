#!/bin/bash

set -ve

# proto
go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
go install github.com/noble-gase/og/cmd/protoc-gen-og@latest

# buf
go install github.com/bufbuild/buf/cmd/buf@latest

# swagger
go install github.com/go-swagger/go-swagger/cmd/swagger@latest
go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv2@latest
