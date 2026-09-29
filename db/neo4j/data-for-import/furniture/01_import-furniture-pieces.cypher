:auto LOAD CSV WITH HEADERS FROM 'file:///var/lib/neo4j/import/furniture-pieces.csv' AS line
CALL {
  WITH line
  MATCH (f:Facility {code: 'B'})
  MATCH (parent:System {systemCode: 'FUR'})
  MATCH (loc:Location {code: line.roomCode, facility: 'B'})
  MATCH (ci:CatalogueItem {catalogueNumber: line.catalogueNumber})
  MERGE (sys:System {systemCode: line.systemCode})
    ON CREATE SET sys.uid = apoc.create.uuid(),
                  sys.name = ci.name,
                  sys.systemLevel = 'SUBSYSTEMS_AND_PARTS',
                  sys.isTechnologicalUnit = false,
                  sys.deleted = false,
                  sys.lastUpdateTime = datetime(),
                  sys.lastUpdateBy = 'furniture-import'
  MERGE (sys)-[:BELONGS_TO_FACILITY]->(f)
  MERGE (parent)-[:HAS_SUBSYSTEM]->(sys)
  MERGE (sys)-[:HAS_LOCATION]->(loc)
  MERGE (sys)-[:CONTAINS_ITEM]->(itm:Item {name: sys.systemCode})
    ON CREATE SET itm.uid = apoc.create.uuid(), itm.lastUpdateTime = datetime()
  MERGE (itm)-[:IS_BASED_ON]->(ci)
} IN TRANSACTIONS OF 500 ROWS;
