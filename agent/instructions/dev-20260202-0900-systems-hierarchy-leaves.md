# Systems Hierarchy Explorer API (Systems service)

- Date: 2026-02-02
- Branch: dev
- Scope: REST API only (eli-panda-api)

## Goal

Provide two lightweight read endpoints for the new System Hierarchy Explorer UI:

1. A **parents-only hierarchy tree** (loads first, cheap)
2. A **leaf systems table** for a selected parent (loaded on selection, paginated)

## Endpoints

### 1) Hierarchy tree

- Method + path: `GET /v1/systems/hierarchy`
- Auth: `ROLE_SYSTEMS_VIEW`
- Content-Type: `application/json`

#### Response

- `200` with JSON array `[]models.SystemHierarchyNode`

`SystemHierarchyNode` is intentionally minimal:

- `uid` (string)
- `name` (string)
- `systemCode` (string, optional)
- `hasLeafChildren` (bool)
- `children` (`[]SystemHierarchyNode`) — **contains only parent nodes**

Notes:

- The tree includes only systems that have at least one subsystem.
- `children` only contains children that are also parents (i.e., themselves have subsystems).
- `hasLeafChildren=true` indicates the node has at least one direct child that is a leaf; the UI can use this to show that leaf systems are available even if `children` is empty.

### 2) Leaves for parent

- Method + path: `GET /v1/system/{uid}/leaves`
- Auth: `ROLE_SYSTEMS_VIEW`
- Content-Type: `application/json`

#### Query params

- `pagination` (required): JSON string, e.g. `{ "page": 1, "pageSize": 50 }`
  - Fallback supported: `page` + `pageSize`
- `sorting` (optional): JSON string array, e.g. `[{"id":"name","desc":false}]`
- `search` (optional): string (matches system name and systemCode)
- `columnFilter` (optional): JSON string array of `{id,value,...}`
  - Supported filters (currently): `zone`, `systemType`

#### Response

- `200` with `helpers.PaginationResult[models.System]`

This reuses the existing `models.System` projection used by the systems list endpoint (includes `systemType`, `zone`, `location`, `physicalItem`, etc.).

#### Semantics

- Returns **direct child** leaf systems of the given parent:
  - child: `(parent)-[:HAS_SUBSYSTEM]->(sys)`
  - leaf: `NOT (sys)-[:HAS_SUBSYSTEM]->(:System{deleted:false})`

If the UI later needs **recursive descendant leaves**, switch the traversal to `[:HAS_SUBSYSTEM*1..50]` and keep the same “leaf” condition.

#### Sorting safety

Sorting is handled via an allowlist to avoid Cypher injection.
Allowed sort IDs:

- `name`
- `systemCode` / `code`
- `systemType`
- `zone`
- `location`
- `eun`

## Implementation locations

- Routes: `services/systems-service/systems-routes.go`
- Handlers + Swagger annotations: `services/systems-service/systems-handlers.go`
- Service methods: `services/systems-service/systems-service.go`
- DTOs: `services/systems-service/models/model_system.go`
- Cypher builders: `services/systems-service/systems-db-queries.go`

## Local verification

- Run `gofmt` on touched Go files
- Run `go test ./...`
- Run `make swagger` if Swagger docs are updated
