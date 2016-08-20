#!/bin/bash

set -e
#set -x

pushd () {
    command pushd "$@" > /dev/null
}

popd () {
    command popd "$@" > /dev/null 2>&1
}

inotifywait -r -e close_write,move_self *.sh *.go lib;

echo
echo " ---- RUN ---- "
echo

mkdir -p test
./build.sh && pushd test && ../bin/upgrade --type=parent

popd

# self restarting
exec $0
