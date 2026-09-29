MATCH (sys:System)-[:HAS_SUBSYSTEM]-(:System {systemCode: 'FUR'})
WHERE sys.systemCode STARTS WITH 'FUR-'
OPTIONAL MATCH (sys)-[:CONTAINS_ITEM]->(itm:Item)
DETACH DELETE itm, sys;
