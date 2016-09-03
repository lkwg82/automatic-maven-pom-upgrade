#!/bin/bash

set -e
#set -x

bash -n $0

case $(uname) in
 "Darwin")
     brew install fswatch
     watchCommand='fswatch -m fsevents_monitor -x  -1  -r *.go *.sh lib/*.go'
     ;;
 "Linux")
     watchCommand='inotifywait -r -e close_write,move_self *.sh *.go lib'
     ;;
esac

pushd () {
    command pushd "$@" > /dev/null
}

popd () {
    command popd "$@" > /dev/null 2>&1
}

echo "-------------------------------------------------------------------------------------------"
${watchCommand}

echo
echo " ---- RUN ---- "
echo

mkdir -p test

executeProgramm() {
    pushd test && ../bin/upgrade__
    popd
}

set +e
./build.sh && executeProgramm
set -e


# self restarting
exec $0
