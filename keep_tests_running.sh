#!/bin/bash

set -e
#set -x

pushd () {
    command pushd "$@" > /dev/null
}

popd () {
    command popd "$@" > /dev/null 2>&1
}

echo "-------------------------------------------------------------------------------------------"
inotifywait -r -e close_write,move_self *.sh *.go lib;

echo
echo " ---- RUN ---- "
echo

mkdir -p test

executeProgramm() {
    pushd test && ../bin/upgrade
    popd
}

set +e
./build.sh && executeProgramm
set -e


# self restarting
exec $0
