# System Codes API: include systemType + sorting (Systems service)

- Date: 2026-01-28
- Branch: dev
- Scope: REST API only (eli-panda-api)

## Goal

Extend the Control Systems “system codes” endpoints to include `systemType` (as a codebook object with `uid/name/code`) in the response payload and allow sorting by `systemType`.

## Affected endpoints

- `GET /v1/systems/system-codes`
- `GET /v1/systems/system-codes/preview`
- `POST /v1/systems/system-codes`

## Response change

Extend `models.SystemCodesResult`:

- Add `systemType` (optional)
  - Shape: `{ uid, name, code }`

### Mapping rule

The DB layer returns a Cypher map literal `RETURN { ... } AS result`. This is converted via map → JSON → struct.

- The map key must match the JSON tag exactly: `systemType`
- Nested fields must match `Codebook` JSON tags: `uid`, `name`, `code`

## Query changes

Update Cypher builders in `services/systems-service/systems-db-queries.go`:

- `GetSystemsForControlsSystemsQuery`
  - Ensure `st:SystemType` is available in scope (`MATCH`/`OPTIONAL MATCH` already present)
  - Add `systemType: case when st is not null then {uid: st.uid, name: st.name, code: st.code} else null end` to the returned result map

- `GetNewSystemCodesPreviewQuery`
  - Add `systemType: {uid: st.uid, name: st.name, code: st.code}` to the returned result map

- `SaveNewSystemCodesQuery`
  - Add `systemType: {uid: st.uid, name: st.name, code: st.code}` to the returned result map

## Sorting

Update `GetControlSystemsOrderByClauses` whitelist:

- Add allowed sort key: `systemType`
- Sort target: `result.systemType.name`

This prevents Cypher injection by only allowing known columns.

## Implementation locations

- DTO: `services/systems-service/models/model_system.go` (`SystemCodesResult`)
- Queries + sorting: `services/systems-service/systems-db-queries.go`

## Local verification

- Run `gofmt` on touched Go files
- Run `make test`
- Run `make swagger` (keeps `docs/swagger.yaml` + `open-api-specification/panda-api.yaml` in sync)
