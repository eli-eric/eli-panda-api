MATCH (s:System {systemCode: 'FUR'}) WHERE NOT (s)-[:HAS_SUBSYSTEM]->() DETACH DELETE s;

MATCH (i:CatalogueItem)-[:BELONGS_TO_CATEGORY]->(:CatalogueCategory)<-[:HAS_SUBCATEGORY]-(:CatalogueCategory {code: 'furniture'})
WHERE NOT (i)<-[:IS_BASED_ON]-() DETACH DELETE i;

MATCH (:CatalogueCategory {code: 'furniture'})-[:HAS_SUBCATEGORY]->(c:CatalogueCategory)
WHERE NOT (c)<-[:BELONGS_TO_CATEGORY]-() DETACH DELETE c;

MATCH (fur:CatalogueCategory {code: 'furniture'})
WHERE NOT (fur)-[:HAS_SUBCATEGORY]->() AND NOT (fur)<-[:BELONGS_TO_CATEGORY]-() DETACH DELETE fur;
