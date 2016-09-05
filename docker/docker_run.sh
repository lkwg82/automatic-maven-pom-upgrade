#!/usr/bin/env bash

set -ex

bash -n $0

DOCKER_CONTAINER=$(docker build docker | tail -n1 | cut -d\  -f3)
project=src/github.com/lkwg82/automatic-maven-pom-upgrade

docker run -v $(pwd):/go/${project} \
            --workdir /go/${project} \
            -e RELEASE=${RELEASE:-""} \
            -ti ${DOCKER_CONTAINER}
