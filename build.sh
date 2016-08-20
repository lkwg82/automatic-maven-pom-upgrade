#!/bin/bash

set -e
#set -x

go build main.go
go fmt ./...
go test ./...

mkdir -p bin
go build -o bin/upgrade main.go