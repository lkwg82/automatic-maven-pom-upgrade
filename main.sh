#!/usr/bin/env bash

source $(dirname $0)/include.flags.sh
source $(dirname $0)/include.git.sh
source $(dirname $0)/include.maven.sh

exit_if_git_is_dirty

branch="parent"
if [ "$(git_branch_exists ${branch})" == "1" ]; then
   git_branch_checkout_existing ${branch}
else
   git_branch_checkout_new ${branch}
fi

# for development cycles
#git reset --hard master
git reset --hard tags/v1.3.7.RELEASE

# git_merge_updates_from_master

maven_update_parent
git_try_committing_changes

# alias gitreset='git clean -f && git reset --hard tags/v1.3.7.RELEASE'
# ../github/spring-boot/spring-boot-cli