## Neo4j database

Neo4j database is used as a primary data storage.

We are using the open-source version of the Neo4j - community edition.

To keep the data in the database in sync with this API we are using migrations([miagrate library](https://github.com/golang-migrate/migrate)).

Migrations are located in [migrations folder](https://github.com/eli-eric/eli-panda-api/blob/main/db/neo4j/migrations).

Migrations are created following this [tutorial](https://github.com/golang-migrate/migrate/blob/master/database/neo4j/TUTORIAL.md).

When this API is started, migrations are automatically applied.

To create new migration files, you can use create-new-migration.sh script in this directory like this:

`./create-new-migration.sh create_some_indexes`

The argument is a name of the new migration.

### Run Neo4j instance locally

Please follow the instructions in the main [README](https://github.com/eli-eric/eli-panda-api/blob/main/README.md) file.
We have this [development docker-compose file](https://github.com/eli-eric/eli-panda-api/blob/main/docker-compose-dev.yml) with all services needed to run this API in development mode successfully.
