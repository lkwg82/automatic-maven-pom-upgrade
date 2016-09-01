#!/usr/bin/env bash

set -e
set -x

go get github.com/aktau/github-release

function tool() {
 local cmd=$1
 shift
 github-release ${cmd} --user lkwg82 --repo automatic-maven-pom-upgrade --tag "latest" $@
}

tool delete 2>/dev/null || echo "no latest release yet"

tool release --name "latest" --pre-release --description "will_be_released_with_each_successful_commit"
tool upload --name upgrade_linux_amd64 -f bin/upgrade_linux_amd64
tool upload --name upgrade_darwin_amd64 -f bin/upgrade_darwin_amd64