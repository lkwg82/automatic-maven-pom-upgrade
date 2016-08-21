#!/usr/bin/env bash
set -e

cid=$(docker build docker | tail -n1 | cut -d\  -f3)

project=src/github.com/lkwg82/automatic-maven-pom-upgrade

docker run -v $(pwd):/go/${project} --workdir /go/${project} -ti ${cid}