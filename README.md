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

1. Prepare import folder and test data:

   ```bash
   mkdir -p db/neo4j/dev-instance/import
   cp db/neo4j/data-for-import/test-data.cypher db/neo4j/dev-instance/import
   ```

2. Start API + Neo4j:

   ```bash
   docker-compose -f docker-compose-local.yml up -d --build
   ```

3. (Optional) Load test data:

   ```bash
   docker exec -it panda-dev-neo4j cypher-shell -u neo4j -p 'elipanda2022' -f import/test-data.cypher
   ```

4. Open services:
   - Swagger UI: [http://localhost:50000/swagger/index.html](http://localhost:50000/swagger/index.html)
   - API base path: [http://localhost:50000/v1](http://localhost:50000/v1)
   - Neo4j Browser: [http://localhost:7470](http://localhost:7470) (connect using `neo4j://localhost:7680`)

Stop local stack:

```bash
docker-compose -f docker-compose-local.yml down
```

## Local Development (without Docker)

1. Copy environment file and update values if needed:

   ```bash
   cp example.env .env
   ```

2. Install dependencies:

   ```bash
   make install
   ```

3. Run API (includes Swagger generation):

   ```bash
   make run
   ```

## Useful Commands

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

1. **Fork and clone** the repository.
2. **Sync with `dev` branch** and create a feature branch from it:

   ```bash
   git checkout dev
   git pull origin dev
   git checkout -b feat/short-description
   ```

3. Make your changes and run checks locally:

   ```bash
   make test
   make swagger
   ```

4. Ensure touched Go files are formatted (`gofmt`) and commit with a clear message.
5. Push your branch and open a Pull Request **to `dev`**.
6. In the PR description, include:
   - what changed and why,
   - how you tested it,
   - whether API docs (Swagger) were regenerated.

## License

See [LICENSE](LICENSE).
