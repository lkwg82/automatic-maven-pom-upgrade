#!/bin/bash

set -e
#set -x

while true; do
	inotifywait -r -e close_write,move_self *.sh *.go lib;

    echo
    echo " ---- RUN ---- "
    echo

	set +e
	echo running build
	./build.sh || echo "ERROR"

	echo "executing"
	./bin/upgrade || echo "ERROR"
	set -e
done
