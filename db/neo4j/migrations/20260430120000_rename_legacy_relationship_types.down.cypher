// See up.cypher for uniqueness invariant rationale; same reasoning applies
// to the rollback direction.

MATCH (s)-[r:IS_POWERED_FROM]->(t)
MERGE (s)-[r2:IS_POWERED_BY]->(t)
ON CREATE SET r2 = properties(r)
DELETE r;

MATCH (s)-[r:IS_COOLED_FROM]->(t)
MERGE (s)-[r2:IS_COOLED_BY]->(t)
ON CREATE SET r2 = properties(r)
DELETE r;
