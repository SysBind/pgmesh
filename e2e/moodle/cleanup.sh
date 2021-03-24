#!/bin/bash
docker ps -a | grep pgmesh-test- | cut -f1 -d" " | xargs docker stop
docker ps -a | grep pgmesh-test- | cut -f1 -d" " | xargs docker rm

docker network rm pgmesh-test
