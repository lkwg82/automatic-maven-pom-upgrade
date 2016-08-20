#!/bin/bash

set -e

while true; do 
	inotifywait -r -e close_write,move_self *.go lib; 
	./build.sh && ./upgrade; 
done
