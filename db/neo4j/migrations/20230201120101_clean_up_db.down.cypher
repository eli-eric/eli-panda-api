CALL apoc.periodic.iterate('MATCH (n) RETURN n', 'DETACH DELETE n', { batchSize:1000, iterateList: true });
