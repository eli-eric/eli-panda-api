#!/bin/bash
# manual migration down - all

migrationsDir=$(pwd)/migrations

docker run -v $migrationsDir:/migrations \
--network host migrate/migrate \
-path=/migrations/ -database neo4j://neo4j:elipanda2022@localhost:7600?x-multi-statement=true down -all
