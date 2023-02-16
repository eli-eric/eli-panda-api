# ELI - PANDA REST API

The only way how to access data in PANDA database.

Using [Echo](https://echo.labstack.com/) - High performance, extensible, minimalist Go web framework - for now it is one classic REST API but we expect to switching to microservices in the future if it makes sens. Now the effort is to stick to the style of "Vertical sliced architecture" to make it easy to switch to microservices.

## How to run the development server

### Using docker compose

Create import folder to import test data into the container:

`mkdir db/neo4j/dev-instance/import`

Copy test data to the import folder:

`cp db/neo4j/data-for-import/test-data.cypher db/neo4j/dev-instance/import`

Run database and API in docker:

`docker-compose -f docker-compose-dev.yml up -d`

It will run all necessary services including neo4j and REST API itself.

Please be patient, neo4j takes some time to start due to additional plugins and then the API takes some time to apply migrations.

After startup you should reach the API docs on: [localhost:50000](http://localhost:50000)

The REST API running on: [localhost:50000/v1](http://localhost:50000/v1)

Neo4j browser on: [localhost:7470](http://localhost:7470) (change Connect URL to: neo4j://localhost:7680, please) with credentials neo4j/elipanda2022

You can end the dev server and all the related docker services by running:

`docker-compose -f docker-compose-dev.yml down`

### Populate databse with test data

`docker exec -it panda-dev-neo4j cypher-shell -u neo4j -p 'elipanda2022' -f import/test-data.cypher`

### Runing localy

comming soon :)
