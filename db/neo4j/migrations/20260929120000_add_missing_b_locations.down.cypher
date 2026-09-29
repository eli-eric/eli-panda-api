MATCH (l:Location {facility: 'B'})
WHERE l.code IN ['LB.00.17', 'LB.00.18', 'LB.1.40', 'M.00.31', 'O.2.27']
  AND NOT (l)-[:HAS_SUBLOCATION]->()
  AND NOT (l)-[:HAS_ROOM_CARD]->()
  AND NOT (:System)-[:HAS_LOCATION]->(l)
DETACH DELETE l;
