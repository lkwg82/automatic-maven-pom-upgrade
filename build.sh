#!/bin/bash

set -e

go test ./...
go build -o upgrade main.go
