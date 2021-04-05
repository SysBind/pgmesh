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

DOCKER_RUN="docker run --network=pgmesh-test"
DOCKER_EXC="docker exec"
PG_ENV=" -e POSTGRES_PASSWORD=q1w2e3r4 -e POSTGRES_DB=moodle"
PG_CLIENT_ENV="-e PGHOST=pgmesh-test-pg${FROM_VERSION} -e PGUSER=postgres -e PGPASSWORD=q1w2e3r4 -e PGDATABASE=moodle -e PGNEWHOST=pgmesh-test-pg${TO_VERSION}"

# Postgres (prev version)
stdout "running postgres $FROM_VERSION"
$DOCKER_RUN -d  --name pgmesh-test-pg${FROM_VERSION} ${PG_ENV} -v "$PWD/postgresql.conf":/etc/postgresql/postgresql.conf postgres:${FROM_VERSION} -c 'config_file=/etc/postgresql/postgresql.conf' || fail "failed running postgres:$FROM_VERSION"

while ! docker logs pgmesh-test-pg${FROM_VERSION} | grep "database system is ready to accept connections"; do    
    stdout "waitin for postgres"
    sleep 2s
done

# Moodle
stdout "building moodle test image"
docker build --quiet . -f Dockerfile -t pgmesh-test-moodle || fail "failed building moodle test image"

stdout "running moodle test image"
$DOCKER_RUN -d $PG_CLIENT_ENV -p 80:80 --name pgmesh-test-moodle pgmesh-test-moodle  2>&1 > /dev/null || fail "failed running moodle test image"

stdout "testing postgres connexion"
if ! $DOCKER_EXC pgmesh-test-moodle psql -c "\l"; then
    fail "Database Connexion Failed"
fi

stdout "running moodle database install"
$DOCKER_EXC pgmesh-test-moodle php admin/cli/install_database.php --agree-license --adminpass=q1w2e3r4

stdout "completing admin registration"
$DOCKER_EXC pgmesh-test-moodle psql -c "UPDATE mdl_user SET firstname = 'Admin',lastname = 'User',email = 'test@user.com',maildisplay = '1',country = 'IL',timezone = '99' WHERE id=2"
$DOCKER_EXC pgmesh-test-moodle psql -c "INSERT INTO mdl_user_preferences (name,value,userid) VALUES('email_bounce_count', '1',2)"
$DOCKER_EXC pgmesh-test-moodle psql -c "INSERT INTO mdl_user_preferences (name,value,userid) VALUES('email_send_count', '1', 2)"

# Postgres (new version)
stdout "running postgres $TO_VERSION"
$DOCKER_RUN -d --name pgmesh-test-pg${TO_VERSION} ${PG_ENV} postgres:${TO_VERSION} || fail "failed running postgres:$TO_VERSION"

while ! docker logs pgmesh-test-pg${TO_VERSION} | grep "database system is ready to accept connections"; do    
    stdout "waitin for postgres"
    sleep 2s
done


# Run pgmesh to establish logical replication from old to new Postgres
$DOCKER_RUN --name pgmesh-test-pgmesh pgmesh pubsub --source-host=pgmesh-test-pg${FROM_VERSION} --dest-host=pgmesh-test-pg${TO_VERSION} --source-database=moodle --dest-database=moodle || fail "failed running pgmesh pubsub"


# Logical Replication Established - Generate Some Traffic
stdout "generating test data"
$DOCKER_EXC pgmesh-test-moodle php admin/tool/generator/cli/maketestsite.php --bypasscheck --size XS

