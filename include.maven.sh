#!/usr/bin/env bash

source $(dirname $0)/include.flags.sh

__maven_versions_plugin_version=2.3
__maven_plugin="org.codehaus.mojo:versions-maven-plugin:${__maven_versions_plugin_version}"

# use the maven wrapper if available
if [ -f "mvnw" -a -x "mvnw" ];
then
  alias mvn='./mvnw'
fi

function __maven_update_parent {
    echo -n "checking for parent pom update ... "
    mvn ${__maven_plugin}:update-parent -DgenerateBackupPoms=false \
        | tee -a debug.log > maven.log 2>&1 \
        && echo "ok" && return 0 \
        || echo "fail" && return 1
}

function maven_update_parent {
    __maven_update_parent &&
        cat maven.log | grep "Updating parent from" | sed -e 's#^\[INFO\] ##' > commit.msg || \
            cat maven.log | grep -E "WARNING" maven.log | sed -e 's#WARNING#INFO#'\
    || cat maven.log
}