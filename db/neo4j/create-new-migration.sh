#!/bin/bash

migrationsDir=$(pwd)/migrations

docker run -v $migrationsDir:/migrations \
-u $(id -u ${USER}):$(id -g ${USER}) \
migrate/migrate \
create -dir migrations -ext cypher $1