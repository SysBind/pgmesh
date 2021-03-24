#!/bin/bash

FROM_VERSION=10
TO_VERSION=13

# enable unofficial bash strict mode
set -o errexit
set -o nounset
set -o pipefail

# Network
docker network create pgmesh-test

DOCKER_RUN="docker run -d --network=pgmesh-test"

# Postgres (prev version)
$DOCKER_RUN --name pgmesh-test-pg${FROM_VERSION} -e POSTGRES_PASSWORD=q1w2e3r4 postgres:${FROM_VERSION}

# Moodle
docker build . -f Dockerfile -t pgmesh-test-moodle
$DOCKER_RUN -p 80:80 --name pgmesh-test-moodle pgmesh-test-moodle
