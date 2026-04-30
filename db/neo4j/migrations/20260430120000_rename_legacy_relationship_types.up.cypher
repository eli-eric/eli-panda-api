MATCH (s)-[r:IS_POWERED_BY]->(t)
CREATE (s)-[r2:IS_POWERED_FROM]->(t)
SET r2 = properties(r)
DELETE r;

MATCH (s)-[r:IS_COOLED_BY]->(t)
CREATE (s)-[r2:IS_COOLED_FROM]->(t)
SET r2 = properties(r)
DELETE r;
