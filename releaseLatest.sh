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

DESCRIPTION=$(echo -n "will_be_released_with_each_successful_commit_on_master_$(git rev-parse HEAD) at $(date)" | sed -e 's# #_#g')
tool release --pre-release --description ${DESCRIPTION}
tool upload --name upgrade_linux_amd64 -f bin/upgrade_linux_amd64
tool upload --name upgrade_darwin_amd64 -f bin/upgrade_darwin_amd64
tool upload --name upgrade_darwin_386 -f bin/upgrade_darwin_386