#!/usr/bin/env bash

set -e
set -x

go get github.com/aktau/github-release

function tool() {
 local cmd=$1
 shift
 github-release ${cmd} --user lkwg82 --repo automatic-maven-pom-upgrade --tag "latest" $@
}

tool delete || echo "no latest release yet"

git config user.email "builds@travis-ci.com"
git config user.name "Travis CI"
git remote add origin2 git@github.com:${TRAVIS_REPO_SLUG}.git
find .
ssh-add travis

git push --delete origin2 latest || echo "no latest tag"

DESCRIPTION=$(echo -n "will_be_released_with_each_successful_commit__${TRAVIS_COMMIT} at ${TRAVIS_STACK_TIMESTAMP}" | sed -e 's# #_#g')
tool release --pre-release --description ${DESCRIPTION}
tool upload --name upgrade_linux_amd64 -f bin/upgrade_linux_amd64
tool upload --name upgrade_darwin_amd64 -f bin/upgrade_darwin_amd64