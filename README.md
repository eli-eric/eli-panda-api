# ELI PANDA REST API

REST API for the ELI PANDA maintenance and operations platform.

The service is built with [Echo](https://echo.labstack.com/) and follows a vertical-slice service structure to keep domain logic grouped and maintainable.

**Production Swagger docs:** https://panda-api.eli-laser.eu/swagger/index.html

## Tech Stack

- Go `1.22`
- Echo web framework
- Neo4j as the primary datastore (with startup migrations)
- JWT-based authentication and role-based authorization
- Swagger / OpenAPI documentation generated via `swag`

## Project Structure

- `server.go` – application entrypoint and middleware/bootstrap wiring
- `services/` – domain slices (`systems`, `catalogue`, `orders`, `publications`, `room-cards`, `security`, etc.)
- `middlewares/` – CORS, logging, recovery, auth middleware
- `db/neo4j/` – migration and local data/import assets
- `docs/` and `open-api-specification/` – generated API docs

## Quick Start (Docker)

1. Create an ignored local environment file:

   ```bash
   cp example.env .env
   ```

   The defaults connect host-run tools to the Neo4j Bolt port published at `localhost:7680`. Leave the WoS key empty when running mocked tests. If you need a live lookup, obtain the key through ELI's secret-management process and put it only in `.env`; never commit it.

2. (Optional) Prepare the broader test-data fixture before starting Neo4j for the first time:

   ```bash
   mkdir -p db/neo4j/dev-instance/import
   cp db/neo4j/data-for-import/test-data.cypher db/neo4j/dev-instance/import
   ```

   The broader fixture is optional. The container owns the bind-mounted import directory after its first start.

3. Start the API and Neo4j 4.4:

   ```bash
   docker compose -f docker-compose-local.yml up -d --build
   ```

   To start only Neo4j, for example before running Go tests on the host:

   ```bash
   make db-local-up
   ```

   Normal schema and reference-data migrations run automatically when the API starts. If you prepared the optional broader fixture, load it after Neo4j is healthy:

   ```bash
   docker exec -it panda-dev-neo4j cypher-shell -u neo4j -p 'elipanda2022' -f import/test-data.cypher
   ```

4. Open services:
   - Swagger UI: [http://localhost:50000/swagger/index.html](http://localhost:50000/swagger/index.html)
   - API base path: [http://localhost:50000/v1](http://localhost:50000/v1)
   - Neo4j Browser: [http://localhost:7470](http://localhost:7470) (connect using `neo4j://localhost:7680`)

Stop local stack:

```bash
docker compose -f docker-compose-local.yml down
```

### Web of Science locally

`docker-compose-local.yml` forwards `API_INTEGRATION_B_WOS_STARTER_API_URL` and `API_INTEGRATION_B_WOS_STARTER_API_KEY` from the ignored `.env` file into the API container. The URL may point to Clarivate only when a valid key is available, or to an explicitly configured mock server. Automated tests must mock the WoS HTTP server and must not require or call with a live key.

The WoS import endpoints require the `publications-edit` role. Local migrations create publication roles but do not grant them to a user. Grant the local-only test account the publication roles, replacing `test` with your local username if needed:

```bash
docker exec panda-dev-neo4j cypher-shell -u neo4j -p 'elipanda2022' \
  'MATCH (u:User {username: "test"}), (r:Role)
   WHERE r.code IN ["publications-view", "publications-edit"]
   MERGE (u)-[:HAS_ROLE]->(r)
   RETURN u.username, collect(r.code)'
```

Authenticate again after changing roles so the new JWT contains them. Do not apply local role grants to shared development or production databases.

### Running tests

The repository requires Go 1.22. Some service tests use the Neo4j connection from `.env`, so start the local database first and ensure `NEO4J_PORT=7680`. WoS tests use a local HTTP test server; no Clarivate key is required.

```bash
make db-local-up
make test
```

Use only a disposable local database for tests. The shared test setup does not perform automatic database-wide cleanup.

## Local Development (without Docker)

1. Copy environment file and update values if needed:

   ```bash
   cp example.env .env
   ```

2. Install dependencies:

   ```bash
   make install
   ```

   `make swagger` and `make run` use `go run github.com/swaggo/swag/cmd/swag@v1.16.3`, so a separate global `swag` installation is not required.

   If you want the standalone CLI anyway:

   ```bash
   go install github.com/swaggo/swag/cmd/swag@v1.16.3
   ```

3. Run API (includes Swagger generation):

   ```bash
   make run
   ```

## Useful Commands

- `make db-local-up` – start the local Neo4j 4.4 service
- `make db-local-down` – stop the local Neo4j service without deleting its data
- `make db-local-status` – show the local Neo4j container status
- `make swagger` – regenerate Swagger and OpenAPI files
- `make build` – build the API binary
- `make test` – run all tests (`go test ./...`)
- `make run` – generate docs and start the server

## Authentication

Most endpoints under `/v1` are protected by JWT middleware and role checks.
Use `/v1/authenticate` to obtain a token, then send it in the `Authorization` header as `Bearer <token>`.

## Troubleshooting

- **Neo4j is not ready yet:** first startup can take longer because plugins initialize and migrations run.
- **Swagger not updated after handler changes:** run `make swagger` and restart the API.
- **Cannot connect to local Neo4j Browser data:** in Neo4j Browser use `neo4j://localhost:7680` (not the default `7687` mapping from host).
- **401/403 responses:** verify JWT token validity and that the user has required role permissions for the endpoint.

## Contributing

1. Clone the `eli-eric/eli-panda-api` repository directly. Do not create a fork for ELI work.
2. Sync the default `dev` branch and create a Jira-prefixed branch from it:

   ```bash
   git switch dev
   git pull --ff-only origin dev
   git switch -c ELIPANDA-123-short-description
   ```

3. Make your changes and run checks locally:

   ```bash
   make test
   make swagger
   ```

4. Ensure touched Go files are formatted (`gofmt`) and commit with a clear message.
5. Push the Jira branch to the ELI repository and open a Pull Request **to `dev`**. Never push directly to `dev`, `main`, or `production`, and do not merge without review.
6. In the PR description, include:
   - what changed and why,
   - how you tested it,
   - whether API docs (Swagger) were regenerated.

## License

See [LICENSE](LICENSE).
