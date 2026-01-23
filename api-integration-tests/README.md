# API Integration Tests

This folder contains **manual API integration test scenarios** for the PANDA REST API.

The goal is to validate real behavior through HTTP (routing + handlers + services + Neo4j queries) in a **feature-based / vertical-slice** structure.

At the moment these scenarios are intended for **semi-automatic/manual checks**:

- You **prepare (seed) data manually** in Neo4j (e.g. Neo4j Browser).
- Then you run a single command (`newman`) that executes a Postman collection against a running API.

In the future we can upgrade scenarios to automatically run seed/cleanup scripts, but that is intentionally **out of scope for now**.

## Prerequisites

- PANDA REST API is running and reachable (e.g. `http://localhost:50000/v1`).
- Neo4j is running and reachable (used by the API).
- `newman` is installed (Postman CLI runner).

Optional (recommended): `jq` for inspecting responses during debugging.

## Shared environment settings (root)

All scenarios should use the same environment variables, defined once in the root of this folder.

Create a local file (not committed) like:

```bash
api-integration-tests/env.local.sh
```

Example content:

```bash
export PANDA_API_BASE_URL="http://localhost:50000/v1"
export PANDA_API_BEARER_TOKEN="<paste-token-here>"
```

Usage:

```bash
source api-integration-tests/env.local.sh
bash api-integration-tests/ELIPANDA-480-copy-systems/run.sh
```

Notes:

- `PANDA_API_BASE_URL` should point to the REST API base (including `/v1`).
- `PANDA_API_BEARER_TOKEN` is passed to Newman and used by requests as a Bearer token.

## Folder structure (feature-based)

Each scenario lives in its own folder (usually mapped to a ticket):

```
api-integration-tests/
	ELIPANDA-<ticket>-<short-slug>/
		README.md
		run.sh
		postman/
			collection.json
			environment.json        # optional
		seed.cypher               # optional, for manual copy-paste/import
		cleanup.cypher            # optional, for manual copy-paste/import
```

### Scenario `README.md`

Each scenario README should contain:

- What the scenario is verifying (business behavior)
- Required manual seeding steps (what to run in Neo4j Browser)
- Which endpoints are called
- Expected outcomes (status codes + key response fields)

## `run.sh` requirements

For consistency, each scenario `run.sh` should contain **only a Newman call** (no seeding, no docker orchestration).

It should assume these env vars are already set:

- `PANDA_API_BASE_URL`
- `PANDA_API_BEARER_TOKEN`

Recommended flags:

- `--env-var baseUrl="$PANDA_API_BASE_URL"`
- `--env-var bearerToken="$PANDA_API_BEARER_TOKEN"`

If you need multiple tokens/users, prefer defining extra env vars and passing them as additional `--env-var` values.

## Manual seeding and cleanup

If a scenario includes `seed.cypher` / `cleanup.cypher`, treat them as **helper scripts** for humans:

- Copy/paste into Neo4j Browser, or
- Import via `cypher-shell` manually.

They are intentionally not executed by `run.sh` yet.
