## Neo4j database

Neo4j database is used as a primary data storage.

We are using the open-source version of the Neo4j - community edition.

To keep the data in the database in sync with this API we use migrations ([golang-migrate](https://github.com/golang-migrate/migrate)).

Migrations are located in the [migrations folder](./migrations).

Migrations are created following this [tutorial](https://github.com/golang-migrate/migrate/blob/master/database/neo4j/TUTORIAL.md).

When this API is started, migrations are automatically applied.

To create new migration files, you can use create-new-migration.sh script in this directory like this:

`./create-new-migration.sh create_some_indexes`

The argument is a name of the new migration.

### Run Neo4j instance locally

Follow the local setup in the main [README](../../README.md). The local stack is defined in [`docker-compose-local.yml`](../../docker-compose-local.yml) and pins Neo4j Community 4.4 with APOC.

Start only the database with `make db-local-up`. Neo4j Browser is available at `http://localhost:7470`, and host-run API/tests connect to `bolt://localhost:7680`. The API container connects to `bolt://panda-dev-neo4j:7687` on the Compose network.

Starting Neo4j alone does not apply PANDA migrations. Migrations run when the API starts. Use `make db-local-down` to stop the database without deleting the bind-mounted local data.
