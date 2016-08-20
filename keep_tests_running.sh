#!/bin/bash

set -e
#set -x

inotifywait -r -e close_write,move_self *.sh *.go lib;

echo
echo " ---- RUN ---- "
echo

./build.sh && ./bin/upgrade

# self restarting
exec $0
