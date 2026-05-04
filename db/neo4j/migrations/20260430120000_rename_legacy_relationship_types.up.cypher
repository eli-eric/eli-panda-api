MATCH (s)-[r:IS_POWERED_BY]->(t)
MERGE (s)-[r2:IS_POWERED_FROM]->(t)
ON CREATE SET r2 = properties(r)
DELETE r;

MATCH (s)-[r:IS_COOLED_BY]->(t)
MERGE (s)-[r2:IS_COOLED_FROM]->(t)
ON CREATE SET r2 = properties(r)
DELETE r;
