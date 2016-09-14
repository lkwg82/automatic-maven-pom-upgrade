#!/bin/bash

set -e

env | sort | grep ^TRAVIS

# only start release process of latest, when on master
if [ "${TRAVIS_PULL_REQUEST}" == "false" ]; then
  ./buildRelease.sh
  ./releaseLatest.sh
else
  ./build.sh
  echo "skipping latest release, as we are on pull request '${TRAVIS_PULL_REQUEST}'"
fi
