#!/bin/bash

set -e
#set -x

mkdir -p logs

go get github.com/alexcesaro/log
go get github.com/alexcesaro/log/golog
go get github.com/droundy/goopt
go get github.com/go-errors/errors
go get github.com/rafecolton/go-fileutils
go get github.com/stretchr/testify/assert

echo " run tests"
go test -v ./...
#> logs/test.log && rm logs/test.log || bash -c 'cat logs/test.log && exit 1'
#go fmt ./...

# cross compiling
# https://golang.org/doc/install/source#environment

mkdir -p bin
echo " building"

build() {
  go build -ldflags="-s -w" -o bin/upgrade_${GOOS}_${GOARCH} main.go
  upx -9 bin/upgrade_${GOOS}_${GOARCH}
}

export GOOS=linux
export GOARCH=amd64
build

export GOOS=darwin
export GOARCH=amd64
build

#mkdir -p dist
#cp bin/upgrade dist/upgrade-$(git rev-parse HEAD)-linux-amd64