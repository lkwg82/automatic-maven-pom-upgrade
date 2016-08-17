#!/usr/bin/env bash

function assertOneArg {
    if [ "$#" != 1 ]; then
        echo "expected one argument in \"$(basename ${BASH_SOURCE[1]}):${BASH_LINENO[0]}:${FUNCNAME[1]}()\""
        exit 1
    fi
}
