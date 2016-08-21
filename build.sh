#!/bin/bash

set -e
#set -x

mkdir -p logs

go get github.com/droundy/goopt
go get github.com/rafecolton/go-fileutils
go get github.com/stretchr/testify/assert

echo " run tests"
go test ./...
#> logs/test.log && rm logs/test.log || bash -c 'cat logs/test.log && exit 1'
#go fmt ./...

mkdir -p bin
echo " building"
go build -ldflags="-s -w" -o bin/upgrade.ldflags main.go > logs/build.log && rm logs/build.log || bash -c 'cat logs/build.log && exit 1'
