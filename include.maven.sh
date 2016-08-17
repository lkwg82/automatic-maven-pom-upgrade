#!/usr/bin/env bash

source $(dirname $0)/include.flags.sh

# use the maven wrapper if available
if [ -f "mvnw" -a -x "mvnw" ];
then
  alias mvn='./mvnw'
fi

function maven_update_parent {
    mvn versions:update-parent -DgenerateBackupPoms=false
}