# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

**Source of truth: [`AGENTS.md`](./AGENTS.md). Read it first — it has the full build/run/test
commands, architecture, auth model, helpers, conventions, and migration rules.**

Critical reference for any Cypher work: **[`db/schema-simple.json`](./db/schema-simple.json)**
— authoritative dump of every Neo4j label, property, and relationship endpoint set.
Consult it before writing or reviewing queries in `*-db-queries.go`, migrations, or
ad-hoc Cypher. Regenerate via `db/schema-simple-query.cypher` (needs APOC) after schema
migrations.
