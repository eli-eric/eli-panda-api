// Description: Import device codes from CSV file by EUN and add missing system type and zone/subzone

//we start by importing missing system types
:auto LOAD CSV WITH HEADERS FROM 'file:///var/lib/neo4j/import/device-codes.csv' AS line
with line WHERE line.SystemTypeCode IS  NOT NULL
MATCH(itm:Item{ eun: line.EUN })<-[:CONTAINS_ITEM]-(sys) 
WHERE sys.systemCode IS NULL OR sys.systemCode = "" 
WITH line, sys WHERE NOT (sys)-[:HAS_SYSTEM_TYPE]->()
MATCH(st:SystemType{code: line.SystemTypeCode})
MERGE(sys)-[:HAS_SYSTEM_TYPE]->(st);

// then we import missing zones and subzones
:auto LOAD CSV WITH HEADERS FROM 'file:///var/lib/neo4j/import/device-codes.csv' AS line
with line WHERE line.Zone IS NOT NULL AND line.SubZone IS NOT NULL
MATCH(itm:Item{ eun: line.EUN })<-[:CONTAINS_ITEM]-(sys) 
WHERE sys.systemCode IS NULL OR sys.systemCode = "" 
WITH line, sys WHERE NOT (sys)-[:HAS_ZONE]->()
MATCH(pz:Zone{code: line.Zone})-[:HAS_SUBZONE]->(z:Zone{code:line.SubZone})
MERGE(sys)-[:HAS_ZONE]->(z);

// finally we import the device codes
:auto LOAD CSV WITH HEADERS FROM 'file:///var/lib/neo4j/import/device-codes.csv' AS line
with line WHERE line.Alias IS NOT NULL
MATCH(itm:Item{ eun: line.EUN })<-[:CONTAINS_ITEM]-(sys) 
WHERE sys.systemCode IS NULL OR sys.systemCode = "" 
WITH line, sys 
SET sys.systemCode = line.Alias;