#!/bin/bash

set -e
#set -x

go get github.com/droundy/goopt
go get github.com/stretchr/testify/assert

go build main.go
#go fmt ./...
go test ./...

mkdir -p bin
go build -o bin/upgrade main.go