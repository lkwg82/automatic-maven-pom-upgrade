#!/bin/bash

set -e
#set -x

go fmt ./...
go test ./...

mkdir -p bin
go build -o bin/upgrade main.go