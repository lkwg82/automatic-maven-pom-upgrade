#!/bin/bash

set -e
#set -x

inotifywait -r -e close_write,move_self *.sh *.go lib;

echo
echo " ---- RUN ---- "
echo

set +e
echo running build
./build.sh || echo "ERROR"

echo "executing"
set -x
./bin/upgrade || echo "ERROR"
set +x
set -e

# self restarting
exec $0
