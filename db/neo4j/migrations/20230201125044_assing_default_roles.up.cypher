MATCH (admin:User{ uid:"71864520-9e86-427c-901c-0c220f951775" })

MATCH (basics:Role{ uid:'6b535e07-2a1c-4498-a1b3-8d63f9275519' })
MATCH (catalogueView:Role{ uid:'85f4310f-044e-45ef-a146-f531e2d89efe' })
MATCH (catalogueEdit:Role{ uid:'b73df5c2-cb61-45af-96b4-5ca93c32e852' })
MATCH (catCategoryEdit:Role{ uid:'b6585153-67ae-4c2e-8d3e-f9d7be519216' })
MATCH (systemsView:Role{ uid:'4a2c7e3a-62f3-4836-93d2-3af202acfcee' })
MATCH (systemsEdit:Role{ uid:'186dee7b-3afa-46ad-84dd-6b62443fec49' })

CREATE (admin)-[:HAS_ROLE]->(basics)
CREATE (admin)-[:HAS_ROLE]->(catalogueView)
CREATE (admin)-[:HAS_ROLE]->(catalogueEdit)
CREATE (admin)-[:HAS_ROLE]->(catCategoryEdit)
CREATE (admin)-[:HAS_ROLE]->(systemsView)
CREATE (admin)-[:HAS_ROLE]->(systemsEdit)
