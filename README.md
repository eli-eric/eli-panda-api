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

`docker-compose -f docker-compose-dev.yml up -d --build`

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

## HTTP Statuses

### Create

Success - 201 Created - Return created object
Failure - 400 Invalid request - Return details about the failure
Async fire and forget operation - 202 Accepted - Optionally return url for polling status

### Update

Success - 200 Ok - Return the updated object
Success - 204 NoContent
Failure - 404 NotFound - The targeted entity identifier does not exist
Failure - 400 Invalid request - Return details about the failure
Async fire and forget operation - 202 Accepted - Optionally return url for polling status

### Patch

Success - 200 Ok - Return the patched object
Success - 204 NoContent
Failure - 404 NotFound - The targeted entity identifier does not exist
Failure - 400 Invalid request - Return details about the failure
Async fire and forget operation - 202 Accepted - Optionally return url for polling status

### Delete

Success - 200 Ok - No content
Success - 200 Ok - When element attempting to be deleted does not exist
Async fire and forget operation - 202 Accepted - Optionally return url for polling status

### Get

Success - 200 Ok - With the list of resulting entities matching the search criteria
Success - 200 Ok - With an empty array

### Get specific

Success - 200 Ok - The entity matching the identifier specified is returned as content
Failure - 404 NotFound - No content

### Action

Success - 200 Ok - Return content where appropriate
Success - 204 NoContent
Failure - 400 - Return details about the failure
Async fire and forget operation - 202 Accepted - Optionally return url for polling status

### Generic results

Authorization error 401 Unauthorized
Authentication error 403 Forbidden
For methods not supported 405
Generic server error 500
