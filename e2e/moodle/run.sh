#!/bin/bash

source "../funcs.sh"

FROM_VERSION=10
TO_VERSION=13

# enable unofficial bash strict mode
set -o errexit
set -o nounset
set -o pipefail

# Network
stdout "creating docker network"
docker network create pgmesh-test

DOCKER_RUN="docker run -d --network=pgmesh-test"
PG_ENV=" -e POSTGRES_PASSWORD=q1w2e3r4 -e POSTGRES_DB=moodle"

# Postgres (prev version)
stdout "running postgres $FROM_VERSION"
$DOCKER_RUN --name pgmesh-test-pg${FROM_VERSION} ${PG_ENV} postgres:${FROM_VERSION} || fail "failed running postgres:$FROM_VERSION"

# Moodle
stdout "building moodle test image"
docker build --quiet . -f Dockerfile -t pgmesh-test-moodle || fail "failed building moodle test image"

stdout "running moodle test image"
$DOCKER_RUN -p 80:80 --name pgmesh-test-moodle pgmesh-test-moodle  2>&1 > /dev/null || fail "failed running moodle test image"
