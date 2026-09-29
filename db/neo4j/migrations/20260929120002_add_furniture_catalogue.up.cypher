MATCH (root:CatalogueCategory {code: 'building-components'})
MERGE (fur:CatalogueCategory {code: 'furniture'})
  ON CREATE SET fur.uid = apoc.create.uuid(), fur.name = 'Furniture'
MERGE (root)-[:HAS_SUBCATEGORY]->(fur);

MATCH (fur:CatalogueCategory {code: 'furniture'})
MERGE (c:CatalogueCategory {code: 'furniture-chairs'})
  ON CREATE SET c.uid = apoc.create.uuid(), c.name = 'Chairs'
MERGE (fur)-[:HAS_SUBCATEGORY]->(c);

MATCH (fur:CatalogueCategory {code: 'furniture'})
MERGE (c:CatalogueCategory {code: 'furniture-tables'})
  ON CREATE SET c.uid = apoc.create.uuid(), c.name = 'Tables and desks'
MERGE (fur)-[:HAS_SUBCATEGORY]->(c);

MATCH (fur:CatalogueCategory {code: 'furniture'})
MERGE (c:CatalogueCategory {code: 'furniture-cabinets'})
  ON CREATE SET c.uid = apoc.create.uuid(), c.name = 'Cabinets'
MERGE (fur)-[:HAS_SUBCATEGORY]->(c);

MATCH (fur:CatalogueCategory {code: 'furniture'})
MERGE (c:CatalogueCategory {code: 'furniture-containers'})
  ON CREATE SET c.uid = apoc.create.uuid(), c.name = 'Mobile containers'
MERGE (fur)-[:HAS_SUBCATEGORY]->(c);

MATCH (c:CatalogueCategory {code: 'furniture-chairs'})
MERGE (i:CatalogueItem {catalogueNumber: 'Z1'})
  ON CREATE SET i.uid = apoc.create.uuid(), i.name = 'Office chair Z1',
                i.description = 'Židle kancelářská Z1', i.lastUpdateTime = datetime()
MERGE (i)-[:BELONGS_TO_CATEGORY]->(c);

MATCH (c:CatalogueCategory {code: 'furniture-chairs'})
MERGE (i:CatalogueItem {catalogueNumber: 'Z2'})
  ON CREATE SET i.uid = apoc.create.uuid(), i.name = 'Conference chair Z2',
                i.description = 'Židle konferenční Z2', i.lastUpdateTime = datetime()
MERGE (i)-[:BELONGS_TO_CATEGORY]->(c);

MATCH (c:CatalogueCategory {code: 'furniture-chairs'})
MERGE (i:CatalogueItem {catalogueNumber: 'Z3'})
  ON CREATE SET i.uid = apoc.create.uuid(), i.name = 'Office chair Z3',
                i.description = 'Židle kancelářská Z3', i.lastUpdateTime = datetime()
MERGE (i)-[:BELONGS_TO_CATEGORY]->(c);

MATCH (c:CatalogueCategory {code: 'furniture-chairs'})
MERGE (i:CatalogueItem {catalogueNumber: 'Z4'})
  ON CREATE SET i.uid = apoc.create.uuid(), i.name = 'Office chair Z4',
                i.description = 'Židle kancelářská Z4', i.lastUpdateTime = datetime()
MERGE (i)-[:BELONGS_TO_CATEGORY]->(c);

MATCH (c:CatalogueCategory {code: 'furniture-chairs'})
MERGE (i:CatalogueItem {catalogueNumber: 'Z5'})
  ON CREATE SET i.uid = apoc.create.uuid(), i.name = 'Dining chair Z5',
                i.description = 'Židle jídelní Z5', i.lastUpdateTime = datetime()
MERGE (i)-[:BELONGS_TO_CATEGORY]->(c);

MATCH (c:CatalogueCategory {code: 'furniture-tables'})
MERGE (i:CatalogueItem {catalogueNumber: 'T1'})
  ON CREATE SET i.uid = apoc.create.uuid(), i.name = 'Working table T1',
                i.description = 'Stůl pracovní T1 / 1800/900/720-740 mm', i.lastUpdateTime = datetime()
MERGE (i)-[:BELONGS_TO_CATEGORY]->(c);

MATCH (c:CatalogueCategory {code: 'furniture-tables'})
MERGE (i:CatalogueItem {catalogueNumber: 'T2'})
  ON CREATE SET i.uid = apoc.create.uuid(), i.name = 'Conference table T2',
                i.description = 'Stůl konferenční T2 / diam. 900/755 mm', i.lastUpdateTime = datetime()
MERGE (i)-[:BELONGS_TO_CATEGORY]->(c);

MATCH (c:CatalogueCategory {code: 'furniture-tables'})
MERGE (i:CatalogueItem {catalogueNumber: 'T3'})
  ON CREATE SET i.uid = apoc.create.uuid(), i.name = 'Working table T3',
                i.description = 'Stůl pracovní T3 / 1600/900/720-740 mm', i.lastUpdateTime = datetime()
MERGE (i)-[:BELONGS_TO_CATEGORY]->(c);

MATCH (c:CatalogueCategory {code: 'furniture-tables'})
MERGE (i:CatalogueItem {catalogueNumber: 'T4'})
  ON CREATE SET i.uid = apoc.create.uuid(), i.name = 'Conference table T4',
                i.description = 'Stůl konferenční T4 / 900/900/720-740 mm', i.lastUpdateTime = datetime()
MERGE (i)-[:BELONGS_TO_CATEGORY]->(c);

MATCH (c:CatalogueCategory {code: 'furniture-tables'})
MERGE (i:CatalogueItem {catalogueNumber: 'T5'})
  ON CREATE SET i.uid = apoc.create.uuid(), i.name = 'Working table T5',
                i.description = 'Stůl pracovní T5 / 1600/800/720-740 mm', i.lastUpdateTime = datetime()
MERGE (i)-[:BELONGS_TO_CATEGORY]->(c);

MATCH (c:CatalogueCategory {code: 'furniture-tables'})
MERGE (i:CatalogueItem {catalogueNumber: 'T6'})
  ON CREATE SET i.uid = apoc.create.uuid(), i.name = 'Sitting desk T6',
                i.description = 'Přísedová deska T6 / 1600/600/720-740 mm', i.lastUpdateTime = datetime()
MERGE (i)-[:BELONGS_TO_CATEGORY]->(c);

MATCH (c:CatalogueCategory {code: 'furniture-tables'})
MERGE (i:CatalogueItem {catalogueNumber: 'T7'})
  ON CREATE SET i.uid = apoc.create.uuid(), i.name = 'Conference table T7',
                i.description = 'Stůl konferenční T7 / 800/800 mm', i.lastUpdateTime = datetime()
MERGE (i)-[:BELONGS_TO_CATEGORY]->(c);

MATCH (c:CatalogueCategory {code: 'furniture-tables'})
MERGE (i:CatalogueItem {catalogueNumber: 'T14'})
  ON CREATE SET i.uid = apoc.create.uuid(), i.name = 'Conference table T14',
                i.description = 'Stůl konferenční T14 / 4*T5 3200/1600/720-740 mm', i.lastUpdateTime = datetime()
MERGE (i)-[:BELONGS_TO_CATEGORY]->(c);

MATCH (c:CatalogueCategory {code: 'furniture-tables'})
MERGE (i:CatalogueItem {catalogueNumber: 'T1-ND'})
  ON CREATE SET i.uid = apoc.create.uuid(), i.name = 'Working table T1 (new design)',
                i.description = 'Stůl pracovní T1 / 1800/900/720-740 mm / Nový design', i.lastUpdateTime = datetime()
MERGE (i)-[:BELONGS_TO_CATEGORY]->(c);

MATCH (c:CatalogueCategory {code: 'furniture-tables'})
MERGE (i:CatalogueItem {catalogueNumber: 'T2-ND'})
  ON CREATE SET i.uid = apoc.create.uuid(), i.name = 'Conference table T2 (new design)',
                i.description = 'Stůl konferenční T2 / diam. 900/755 mm / Nový design', i.lastUpdateTime = datetime()
MERGE (i)-[:BELONGS_TO_CATEGORY]->(c);

MATCH (c:CatalogueCategory {code: 'furniture-tables'})
MERGE (i:CatalogueItem {catalogueNumber: 'T10'})
  ON CREATE SET i.uid = apoc.create.uuid(), i.name = 'Conference table T10',
                i.description = 'Stůl konferenční T10 / 1600/800/740 mm', i.lastUpdateTime = datetime()
MERGE (i)-[:BELONGS_TO_CATEGORY]->(c);

MATCH (c:CatalogueCategory {code: 'furniture-tables'})
MERGE (i:CatalogueItem {catalogueNumber: 'T8'})
  ON CREATE SET i.uid = apoc.create.uuid(), i.name = 'Conference table T8',
                i.description = 'Stůl konferenční T8 / 2200/1000/740 mm', i.lastUpdateTime = datetime()
MERGE (i)-[:BELONGS_TO_CATEGORY]->(c);

MATCH (c:CatalogueCategory {code: 'furniture-containers'})
MERGE (i:CatalogueItem {catalogueNumber: 'K1'})
  ON CREATE SET i.uid = apoc.create.uuid(), i.name = 'Mobile container 3+1 drawer K1',
                i.description = 'Kontejner mobilní 3+1 zásuvka K1 / 800/450/1800 mm', i.lastUpdateTime = datetime()
MERGE (i)-[:BELONGS_TO_CATEGORY]->(c);

MATCH (c:CatalogueCategory {code: 'furniture-cabinets'})
MERGE (i:CatalogueItem {catalogueNumber: 'SV80'})
  ON CREATE SET i.uid = apoc.create.uuid(), i.name = 'High shelf cabinet SV80',
                i.description = 'Skříň vysoká policová SV80 / 800/450/1800 mm', i.lastUpdateTime = datetime()
MERGE (i)-[:BELONGS_TO_CATEGORY]->(c);

MATCH (c:CatalogueCategory {code: 'furniture-cabinets'})
MERGE (i:CatalogueItem {catalogueNumber: 'SV80s'})
  ON CREATE SET i.uid = apoc.create.uuid(), i.name = 'High wardrobe cabinet SV80s',
                i.description = 'Skříň vysoká šatní SV80s / 800/450/1800 mm', i.lastUpdateTime = datetime()
MERGE (i)-[:BELONGS_TO_CATEGORY]->(c);

MATCH (c:CatalogueCategory {code: 'furniture-cabinets'})
MERGE (i:CatalogueItem {catalogueNumber: 'SV81'})
  ON CREATE SET i.uid = apoc.create.uuid(), i.name = 'High combined cabinet SV81',
                i.description = 'Skříň vysoká kombinovaná SV81 / 800/450/1800 mm', i.lastUpdateTime = datetime()
MERGE (i)-[:BELONGS_TO_CATEGORY]->(c);

MATCH (c:CatalogueCategory {code: 'furniture-cabinets'})
MERGE (i:CatalogueItem {catalogueNumber: 'SN80'})
  ON CREATE SET i.uid = apoc.create.uuid(), i.name = 'Low cabinet shelf SN80',
                i.description = 'Skříň nízká policová SN80 / 800/450/720-740 mm', i.lastUpdateTime = datetime()
MERGE (i)-[:BELONGS_TO_CATEGORY]->(c);

MATCH (f:Facility {code: 'B'})
MATCH (d:System {name: 'Facility governance and support systems', systemLevel: 'SYSTEM_DOMAIN'})-[:BELONGS_TO_FACILITY]->(f)
MERGE (s:System {systemCode: 'FUR'})
  ON CREATE SET s.uid = apoc.create.uuid(), s.name = 'Furniture',
                s.systemLevel = 'TECHNOLOGY_UNIT', s.isTechnologicalUnit = true,
                s.deleted = false, s.lastUpdateTime = datetime(),
                s.lastUpdateBy = 'furniture-import'
MERGE (s)-[:BELONGS_TO_FACILITY]->(f)
MERGE (d)-[:HAS_SUBSYSTEM]->(s);
