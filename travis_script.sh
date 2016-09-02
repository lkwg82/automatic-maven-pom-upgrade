#!/bin/bash 

set -e

# run in tests in ubuntu 16.04
./docker/docker_run.sh

# remove compiled binaries as they doesnt have the right permisssions
sudo rm -rf bin

# only start release process of latest, when on master
if [ "${TRAVIS_BRANCH}" == "master" ]; then
  sudo apt-get update && sudo apt-get install -y upx
  env SKIP_TESTS=1 ./buildRelease.sh
  ./releaseLatest.sh
else
  echo "skipping latest release, as we are on branch '${TRAVIS_BRANCH}'"
fi
