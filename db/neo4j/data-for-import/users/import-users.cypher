:auto LOAD CSV WITH HEADERS FROM 'file:///var/lib/neo4j/import/ucty.csv' AS line
WITH line
CALL {
  WITH line
  MATCH (f:Facility{ code: line.facility })
  MERGE (u:User{ uid: line.uid })
   SET u.username = line.username,
  u.firstName = line.firstName,
  u.lastName = line.lastName,
  u.email = line.username,
  u.passwordHash = line.pwdHash,
  u.isEnabled = true
  MERGE (u)-[:BELONGS_TO_FACILITY]->(f)
  WITH split(line.roles, "|") AS roles, u
  MATCH (r:Role)
  WHERE r.code IN roles
  WITH r, u
  MERGE (u)-[:HAS_ROLE]->(r)
  } IN TRANSACTIONS OF 500 ROWS
