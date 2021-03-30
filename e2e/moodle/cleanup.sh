#!/bin/bash

source "../funcs.sh"

if docker ps -a | grep pgmesh-test- 2>&1 > /dev/null; then
    stdout "cleaning up containers"
    docker ps -a | grep pgmesh-test- | cut -f1 -d" " | xargs docker stop 2>&1 > /dev/null
    docker ps -a | grep pgmesh-test- | cut -f1 -d" " | xargs docker rm 2>&1 > /dev/null
fi

if docker network list | grep pgmesh-test 2>&1 > /dev/null; then
    stdout "removing pgmesh-test network"
    docker network rm pgmesh-test
fi
