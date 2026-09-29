MATCH (f:Facility {code: 'B'})
MATCH (p:Location {facility: 'B', name: '00-Ground floor (lasers)'}) WHERE p.code IS NULL
MERGE (l:Location {code: 'LB.00.17', facility: 'B'})
  ON CREATE SET l.uid = apoc.create.uuid(), l.name = 'LB.00.17'
MERGE (p)-[:HAS_SUBLOCATION]->(l)
MERGE (l)-[:BELONGS_TO_FACILITY]->(f);

MATCH (f:Facility {code: 'B'})
MATCH (p:Location {facility: 'B', name: '00-Ground floor (lasers)'}) WHERE p.code IS NULL
MERGE (l:Location {code: 'LB.00.18', facility: 'B'})
  ON CREATE SET l.uid = apoc.create.uuid(), l.name = 'LB.00.18'
MERGE (p)-[:HAS_SUBLOCATION]->(l)
MERGE (l)-[:BELONGS_TO_FACILITY]->(f);

MATCH (f:Facility {code: 'B'})
MATCH (p:Location {facility: 'B', name: '1-first floor'}) WHERE p.code IS NULL
MERGE (l:Location {code: 'LB.1.40', facility: 'B'})
  ON CREATE SET l.uid = apoc.create.uuid(), l.name = 'LB.1.40'
MERGE (p)-[:HAS_SUBLOCATION]->(l)
MERGE (l)-[:BELONGS_TO_FACILITY]->(f);

MATCH (f:Facility {code: 'B'})
MATCH (p:Location {facility: 'B', name: '00-Ground floor'}) WHERE p.code IS NULL
MERGE (l:Location {code: 'M.00.31', facility: 'B'})
  ON CREATE SET l.uid = apoc.create.uuid(), l.name = 'M.00.31'
MERGE (p)-[:HAS_SUBLOCATION]->(l)
MERGE (l)-[:BELONGS_TO_FACILITY]->(f);

MATCH (f:Facility {code: 'B'})
MATCH (p:Location {facility: 'B', name: '2-Second floor'}) WHERE p.code IS NULL
MERGE (l:Location {code: 'O.2.27', facility: 'B'})
  ON CREATE SET l.uid = apoc.create.uuid(), l.name = 'O.2.27'
MERGE (p)-[:HAS_SUBLOCATION]->(l)
MERGE (l)-[:BELONGS_TO_FACILITY]->(f);
