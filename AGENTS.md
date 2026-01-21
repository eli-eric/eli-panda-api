# Repository Guidelines

## Project Structure & Module Organization
- `server.go` is the main entrypoint for the Echo-based API server.
- `services/` hosts domain services (for example, publications-related code lives under `services/publications-service/`).
- `middlewares/`, `helpers/`, `ioutils/`, and `shared/` contain cross-cutting middleware and reusable utilities.
- `db/neo4j/` stores Neo4j fixtures and import data used by local/dev setups.
- `docs/swagger.yaml` is the Swagger source; `open-api-specification/` contains the generated spec and assets.
- `config/` and `example.env` define runtime configuration and environment variable expectations.

## Build, Test, and Development Commands
- `make install`: download and verify Go module dependencies.
- `make swagger`: run `swag init -g server.go` and refresh `open-api-specification/panda-api.yaml`.
- `make run`: generate Swagger then run the API with `go run server.go`.
- `make build`: build the API binary (also regenerates Swagger).
- `make test`: run all Go tests via `go test ./...`.
- Local stack: `docker-compose -f docker-compose-local.yml up -d --build` (stop with `docker-compose -f docker-compose-local.yml down`).

## Coding Style & Naming Conventions
- Use standard Go formatting (run `gofmt` on touched files); Go uses tabs for indentation.
- Prefer idiomatic Go naming: mixedCaps for identifiers, short lowercase package names.
- Keep generated API docs in sync by running `make swagger` whenever you change handlers or Swagger annotations.

## Testing Guidelines
- Tests live next to their packages and follow Go’s `_test.go` naming (examples in `services/` and `helpers/`).
- Run `make test` before submitting; add focused unit tests for new behavior.

## Commit & Pull Request Guidelines
- Commit messages are typically imperative; some use Conventional Commit prefixes like `fix:` or `refactor:`.
- Include relevant issue IDs when applicable (for example, `ELIPANDA-455`).
- PRs should describe what changed, how it was tested, and note any API surface updates (with regenerated Swagger).

## Configuration & Local Data
- Start from `example.env` and keep secrets out of version control.
- For Neo4j test data, see `db/neo4j/data-for-import/` and the README’s import steps.
