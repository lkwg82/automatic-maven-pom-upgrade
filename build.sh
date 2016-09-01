#!/bin/bash

set -e
#set -x

mkdir -p logs

echo " run tests"
go test -v ./...

mkdir -p bin
echo " building"

build() {
  go build -ldflags="-s -w" -o bin/upgrade_${GOOS}_${GOARCH} main.go

  [ -n "$RELEASE" ] && upx --ultra-brute bin/upgrade_${GOOS}_${GOARCH} || echo -n ""
}

build

# cross compiling
# https://golang.org/doc/install/source#environment
export GOOS=linux
export GOARCH=amd64
build

export GOOS=darwin
export GOARCH=amd64
build

unset GOOS
unset GOARCH
