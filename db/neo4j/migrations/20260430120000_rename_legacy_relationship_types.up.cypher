// Uniqueness invariant: at most one IS_POWERED_BY/IS_COOLED_BY edge exists
// between any pair of systems, enforced by the batch create endpoint
// (services/systems-service: CreateBatchRelationships uses
// OPTIONAL MATCH ... existFwd IS NULL guard before CREATE).
// MERGE here is therefore safe to dedupe; if duplicates do exist (e.g. from
// manual data entry), they collapse into one new edge and inherit the
// properties of whichever legacy edge was matched first.

MATCH (s)-[r:IS_POWERED_BY]->(t)
MERGE (s)-[r2:IS_POWERED_FROM]->(t)
ON CREATE SET r2 = properties(r)
DELETE r;

MATCH (s)-[r:IS_COOLED_BY]->(t)
MERGE (s)-[r2:IS_COOLED_FROM]->(t)
ON CREATE SET r2 = properties(r)
DELETE r;
