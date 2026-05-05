MATCH (s)-[r:IS_POWERED_FROM]->(t)
CALL apoc.refactor.setType(r, 'IS_POWERED_BY') YIELD input, output
RETURN count(output);

MATCH (s)-[r:IS_COOLED_FROM]->(t)
CALL apoc.refactor.setType(r, 'IS_COOLED_BY') YIELD input, output
RETURN count(output);
