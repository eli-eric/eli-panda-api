// we want clean database at start, so we delete all nodes and realationships
// delete all nodes and relationships
CALL apoc.periodic.iterate('MATCH (n) RETURN n', 'DETACH DELETE n', { batchSize:1000, iterateList: true });
