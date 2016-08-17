#!/usr/bin/env bash

source $(dirname $0)/include.flags.sh
source $(dirname $0)/include.git.sh
source $(dirname $0)/include.maven.sh

exit_if_git_is_dirty

maven_update_parent
git diff > diff.log
git branch --list "pr-autoupdate-parent"

# lars@lars-MS-7930:~/development/projects/github/spring-boot/spring-boot-cli$ alias gitreset='git reset --hard master && git checkout tags/v1.3.7.RELEASE'