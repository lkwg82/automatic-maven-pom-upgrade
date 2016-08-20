#!/bin/bash

set -e

go test ./...
mkdir -p bin
go build -o bin/upgrade main.go
