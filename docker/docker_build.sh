#!/usr/bin/env bash

set -e

bash -n $0

export DOCKER_CONTAINER=$(docker build docker | tail -n1 | cut -d\  -f3)