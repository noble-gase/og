#!/bin/sh
set -e

buf generate

find api -name "*.pb.go" | xargs -I{} protoc-go-inject-tag -input={}

echo "done"
