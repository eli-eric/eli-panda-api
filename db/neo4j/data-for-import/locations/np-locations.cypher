:auto LOAD CSV WITH HEADERS FROM 'file:///var/lib/neo4j/import/np-locations.csv' AS line
with line
match(f:Facility{code:"N"})
with line, f
MERGE(level1:Location{code:line.level1})-[:BELONGS_TO_FACILITY]->(f)
MERGE(level2:Location{code:line.level1+line.level2})-[:BELONGS_TO_FACILITY]->(f)
MERGE(level3:Location{code: line.code})-[:BELONGS_TO_FACILITY]->(f)
MERGE(level1)-[:HAS_SUBLOCATION]->(level2)-[:HAS_SUBLOCATION]->(level3)
set
level1.name = line.level1,
level1.deleteCode = true,
level2.name = line.level2,
level2.deleteCode = true,
level3.name = line.level3,
level1.facility = "N",
level2.facility = "N",
level3.facility = "N"
    //add manualy uid!! and delete codes level1,level2
    //MATCH(l:Location) where l.uid is null set l.uid = apoc.create.uuid()
    //MATCH(l:Location) where l.deleteCode = true set l.code = null, l.deleteCode = null