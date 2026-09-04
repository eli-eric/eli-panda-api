# Publications: server-side filtering and sorting

`GET /v1/publications` and `GET /v1/publications/export` accept `search`, `sorting` and
`columnFilter`. All three are applied in Cypher, so paging, `totalCount` and the CSV export describe
the same set of rows.

Implementation: `services/publications-service/publications-db-queries.go`
(`ApplyPublicationFilters`), wired into `buildPublicationsQuery` and `buildPublicationsCountQuery`.

## How filters are applied

Every active filter appends one AND-ed predicate to the `WHERE` that both query builders already
open on `(n:Publication)`. Two consequences worth knowing:

- **The list query and the count query share the predicates verbatim**, which is what keeps
  `totalCount` in agreement with the page. A filter added to one and not the other would silently
  desynchronise them.
- **Every predicate holds on the bare publication match**, so codebook and relationship filters are
  written as `EXISTS { MATCH (n)-[:REL]->(x) WHERE ... }` rather than relying on the `OPTIONAL MATCH`
  aliases the list query happens to bind. The count query has no such aliases.

Values are always parameterized. Only ids drawn from the maps in `publications-db-queries.go` are
ever interpolated into the query string.

## Filter ids

Ids are the frontend column ids from `publications.columns.tsx`. Filters with no column (book,
conference and department fields) use the property name.

| Kind | Value shape | Ids |
|---|---|---|
| Text | `"laser"` — case-insensitive CONTAINS | `title`, `code`, `doi`, `allAuthors`, `eliAuthors`, `keywords`, `longJournalTitle`, `shortJournalTitle`, `abstract`, `citeAs`, `wosNumber`, `issn`, `eissn`, `eidScopus`, `oecdFord`, `note`, `otherGrants`, `webLink`, `publisher`, `publishPlace`, `isbn`, `bookTitle`, `editionVolume`, `proceedingsIsbn`, `conferencePlace`, `pages` |
| List | `["2024","2025"]` — matched with IN | `yearOfPublication`, `eliPublication` (YES/NO), `quartil`, `quartilBasis`, `language` |
| Numeric range | `{"min":1.5,"max":9}` — either bound optional | `impactFactor`, `allAuthorsCount`, `eliAuthorsCount`, `pagesCount`, `bookPagesCount`, `volume`, `issue` |
| Date range | `{"min":"2024-01","max":"2024-12"}` | `dateOfPublication`, `conferenceDate` |
| Codebook | `{"uid":"…"}` from a combobox, or `["uid","uid"]` from a checkbox group | `mediaType`, `openAccessType`, `publishingCountry`, `userCall`, `userExperiment`, `experimentalSystem`, `publishFormat`, `conferenceScope`, `department` |
| Relationship list | `["uid","uid"]`, or a single `{"uid":"…"}` | `grant`, `eliResearchers` |

Empty values — a blank string, an empty list, a range with no bounds, a codebook with no uid — are
ignored rather than narrowing the result set.

### Dates are compared as text

`dateOfPublication` and `conferenceDate` are stored as strings, and the register holds `YYYY`,
`YYYY-MM` and `YYYY-MM-DD` side by side. Because those forms are zero-padded, lexical comparison
orders and bounds them correctly, and a partial value like `2024-06` behaves as expected. Parsing
them into dates would have to guess at the missing parts.

### Department is denormalized

A `Publication` has no `HAS_DEPARTMENT` relationship. Departments live in the
`authorsDepartmentsArray` property as `uid||name||count` entries, so the filter is
`ANY(entry IN n.authorsDepartmentsArray WHERE entry STARTS WITH $uid)`. Introducing a real
relationship would let this become an `EXISTS` match like the others, and is the obvious later
optimization.

## Sorting

`publicationSortFields` in `publications-service.go` maps a column id to a property path valid where
the `ORDER BY` runs. Codebook columns sort by the related node's name — `mediaType` becomes
`mediaTypeCb.name` — because that is what the table displays.

**Unknown ids are dropped rather than turned into a property path.** The previous behaviour appended
`"n." + fieldID`, which produced a meaningless sort for a relationship column and an invalid path for
anything else. If no usable id remains, the default `ORDER BY n.updatedAt DESC` applies.

`eliResearchers` and `grant` are deliberately absent: both are assembled in `CALL` subqueries that
run after the sort, so there is nothing to order by.

## Malformed input

The shared readers in `helpers/database.go` assert on the shape of a filter value without checking
it — `GetFilterValueCodebook` (line 940) requires both `uid` and `name`, and
`GetFilterValueListString` (line 948) requires `[]interface{}` of strings. Either panics when a
hand-written `columnFilter` disagrees.

The publications filters therefore read values through local, defensive equivalents
(`filterCodebookUID`, `filterStringList`, `filterString`) that skip anything unusable. Hardening the
shared helpers would fix this for orders, systems and catalogue too, and is worth its own ticket.

## Not covered

- **Facility (CZ/HU).** `Publication` has no `BELONGS_TO_FACILITY` relationship and ALPS data is not
  migrated. Revisit once ALPS publications are in the database.
- **Filter bounds endpoint.** Distinct years and min/max impact factor for the frontend inputs would
  follow `GetMinAndMaxOrderLinePriceQuery`.

## Tests

`services/publications-service/publications-filters_test.go` — every filter in isolation against
seeded data, several combined, empty values ignored, `totalCount` equal to the returned page,
sorting by a relationship column, unknown sort ids rejected, values parameterized rather than
interpolated, and malformed filter values not panicking.
