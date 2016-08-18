#!/usr/bin/env bash

source $(dirname $0)/include.flags.sh
source $(dirname $0)/include.assert.sh

__git_tag_prefix="autoupdate"

# Returns "*" if the current git branch is dirty.
function git_is_dirty {
  [[ $(git diff --shortstat 2> /dev/null | tail -n1) != "" ]] && echo -n "*"
}

function exit_if_git_is_dirty {
    echo -n "checking if repository is fresh and clean ... "
    if [ "$(git_is_dirty)" == "*" ]; then
        echo "dirty"
        exit 1
    else
        echo "clean"
    fi
}

function git_tag_set {
    assertOneArg $@
    git tag ${__git_tag_prefix}_$1
}

function git_tag_delete {
    assertOneArg $@
    git tag --delete ${__git_tag_prefix}_$1
}

function git_tag_exists {
    assertOneArg $@
    git tag --list "${__git_tag_prefix}_$1" | wc -l
}

function git_branch_exists {
    assertOneArg $@
    git branch --list "${__git_tag_prefix}_$1" | wc -l
}

function git_branch_checkout_existing {
    assertOneArg $@
    local branch="${__git_tag_prefix}_$1"
    git checkout ${branch}
}

function git_branch_checkout_new {
    assertOneArg $@
    local branch="${__git_tag_prefix}_$1"
    git checkout -b ${branch}
}

function git_merge_updates_from_master {
    if [[ "$(git branch | grep \* | cut -d\  -f2)" =~ "^${__git_tag_prefix}_.*" ]]; then
        echo "not on a \"${__git_tag_prefix}_*\" branch"
        exit 1
    fi
    echo -n "merging updates from master ... "
    git merge --no-ff master --commit --message "updates from master" >> debug.log && echo "done" \
        || git merge --abort && echo "conflict"
}

# argument: branch to merge grom
function git_check_for_conflicting_merge {
    assertOneArg $@
    echo -n "checking for merge conflict ... "
    git format-patch $1 --stdout | git apply --check - >> debug.log 2>&1 && echo ok || echo conflict
}

function git_try_committing_changes {
    echo "check for committable changes"
    if [ -f "commit.msg" ]; then
        git diff pom.xml > diff.log
        echo -n "committing ... "
        git commit --file commit.msg pom.xml >> debug.log 2>&1 && echo ok && return 0 \
            || echo fail && exit 1
    fi
}
# proceed when executed directly
if [ $(basename $0) == "include.git.sh" ]; then
  set -x
  echo $0

  branch="parent2"
  if [ "$(git_branch_exists ${branch})" == "1" ]; then
       git_branch_checkout_existing ${branch}
  else
       git_branch_checkout_new ${branch}
  fi

  git_merge_updates_from_master
fi