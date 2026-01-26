# Copy System API plan (Systems service)

- Date: 2026-01-21 14:23
- Branch: dev
- Scope: REST API only (eli-panda-api)

## Goal

Add a "copy system" operation to Systems service. The endpoint creates NEW System node(s) by cloning only selected properties/relationships from a source system (or its children) and attaches the created root clone(s) under an existing destination parent system.

## Endpoint

- Method + path: `PUT /v1/systems/copy`
- Auth: `ROLE_SYSTEMS_EDIT`
- Content-Type: `application/json`

## Request

JSON body (models.SystemCopyRequest):

- `sourceSystemUid` (string, required)
- `destinationSystemUid` (string, required) — existing destination parent
- `copyOnlySourceSystemChildren` (boolean)
  - `false`: copy source system itself
  - `true`: copy only the source system's direct children
- `copyRecursive` (boolean)
  - `false`: copy only the selected root(s)
  - `true`: copy the selected root(s) and their subtree via `HAS_SUBSYSTEM`

## What is copied

Only:

- Properties: `name`, `systemLevel`
- Relationships: `HAS_SYSTEM_TYPE`, `HAS_SUBSYSTEM`

Not copied (intentionally): `systemCode`, location/zone, responsible/owner, physical items, files/links, other relationships.

## Preconditions / safety

- Both `sourceSystemUid` and `destinationSystemUid` must exist.
- Both must belong to the same `facilityCode` from auth context.

## Response

- Success: `201` with JSON array `[]string` containing created _root_ system UID(s)
  - 1 UID when copying source itself
  - 0..N UIDs when copying only children (N = number of source direct children)

## Errors

- `400` invalid request (missing required fields)
- `404` source or destination system not found
- `500` unexpected failures

## Implementation notes

- Route: services/systems-service/systems-routes.go
- Handler + Swagger annotations: services/systems-service/systems-handlers.go
- Service method + validation: services/systems-service/systems-service.go
- Cypher builders: services/systems-service/systems-db-queries.go

### Neo4j copy strategy (single write query)

- Match `Facility` by `facilityCode`.
- Match `User` by `userUID` for audit field `lastUpdatedBy`.
- Determine copy roots:
  - if `copyOnlySourceSystemChildren=true`: direct children of source
  - else: the source itself
- For each root:
  - if `copyRecursive=true`: collect nodes in `HAS_SUBSYSTEM*0..50`
  - else: copy only the root node
  - Create cloned System nodes with new UID (via `apoc.create.uuid()`), copy only `name` and `systemLevel`, set `deleted=false` + `lastUpdateTime` + `lastUpdatedBy`, and create `BELONGS_TO_FACILITY`.
  - Recreate `HAS_SYSTEM_TYPE` for each cloned node.
  - Recreate internal `HAS_SUBSYSTEM` edges among cloned nodes.
  - Attach the cloned root node under destination via `HAS_SUBSYSTEM`.

## Local verification

- Run `make swagger`
- Run `make build`
