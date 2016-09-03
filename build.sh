#!/bin/bash

set -e
#set -x

# syntax check
bash -n $0

if [ -z "$SKIP_TESTS" ]; then
    echo " run tests"
    go test -v ./lib/...
fi

mkdir -p bin
echo " building"

determineCurrentOS() {
    case $(uname | tr [:upper:] [:lower:]) in
        "darwin")
            export CURRENT_GOOS="darwin"
        ;;
        "linux")
            export CURRENT_GOOS="linux"
        ;;
    esac
}

determineCurrentArch() {
     set -x
     local arch=$(uname -p)
     set +x
     case  ${arch} in
        "x86_64")
            export CURRENT_GOARCH="amd64"
        ;;
        "unknown")
            echo "falling back to x86_64 architecture"
            export CURRENT_GOARCH="amd64"
        ;;
    esac
}

build() {
  export GOOS=$1
  export GOARCH=$2
  go build -ldflags="-s -w" -o bin/upgrade_${GOOS}_${GOARCH} main.go
  [ -n "$RELEASE" ] || ln -sf bin/upgrade_${GOOS}_${GOARCH} bin/upgrade
  unset GOOS
  unset GOARCH
  _upx $1 $2
}

_upx() {
  [ -n "$RELEASE" ] && upx --ultra-brute bin/upgrade_$1_$2 || echo -n ""
}

# cross compiling
# https://golang.org/doc/install/source#environment
determineCurrentOS
determineCurrentArch

if [ -n "$RELEASE" ]; then
            build darwin amd64
            build linux amd64
else
    if [ "$CURRENT_GOARCH" == "amd64" ]; then
        if [ "$CURRENT_GOOS" == "darwin" ]; then
            build darwin amd64
        elif [ "$CURRENT_GOOS" == "linux" ]; then
            build linux amd64
        fi
    else
        echo "dont support current arch: $CURRENT_GOARCH"
        exit 1
    fi
fi




