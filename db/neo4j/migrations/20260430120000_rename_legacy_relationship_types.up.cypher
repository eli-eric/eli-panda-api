MATCH (s)-[r:IS_POWERED_BY]->(t)
CALL apoc.refactor.setType(r, 'IS_POWERED_FROM') YIELD input, output
RETURN count(output);

MATCH (s)-[r:IS_COOLED_BY]->(t)
CALL apoc.refactor.setType(r, 'IS_COOLED_FROM') YIELD input, output
RETURN count(output);
