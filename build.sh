#!/bin/bash

set -e
#set -x

mkdir -p logs

if [ -z "$SKIP_TESTS" ]; then
    echo " run tests"
    go test -v ./lib/...
fi

mkdir -p bin
echo " building"

build() {
  go build -ldflags="-s -w" -o bin/upgrade_${GOOS}_${GOARCH} main.go
}

_upx() {
  [ -n "$RELEASE" ] && upx --ultra-brute bin/upgrade_${GOOS}_${GOARCH} || echo -n ""
}

build

# cross compiling
# https://golang.org/doc/install/source#environment
export GOOS=linux
export GOARCH=amd64
build
_upx

export GOOS=darwin
export GOARCH=amd64
build
_upx

unset GOOS
unset GOARCH
