// Description: countries
// // exmaple:
// "country","latitude","longitude","name"
// "AD",42.546245,1.601554,"Andorra"
//we start by importing missing system types
:auto LOAD CSV WITH HEADERS FROM 'file:///var/lib/neo4j/import/countries.csv' AS line
with line
MERGE (c:Country{code: line.country})
SET c.uid = coalesce(c.uid, apoc.create.uuid()),
 c.latitude = line.latitude, 
 c.longitude = line.longitude, 
 c.name = line.name;