MATCH(empl:Employee) WITH empl.lastName+empl.firstName as emplNames, empl
MATCH(usr:User) WHERE usr.lastName+usr.firstName in emplNames
MERGE (empl)-[:HAS_USER]->(usr);