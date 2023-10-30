:auto 
LOAD CSV WITH HEADERS FROM 'file:///var/lib/neo4j/import/qr_codes_import.csv' AS line
  WITH line where line.import = "y" and line.orderId is not null
CALL {
  WITH line
  MATCH(f:Facility{code: "B"})     
  WITH line, f
  MERGE(o:Order{ orderNumber: line.orderId})
  MERGE(catItem:CatalogueItem{catalogueNumber: line.catalogueNumber })
  MERGE(itm:Item{ eun: line.eun })  
  MERGE(o)-[:BELONGS_TO_FACILITY]->(f)
  MERGE(o)-[:HAS_ORDER_LINE{isDelivered: true}]->(itm)
  MERGE(itm)-[:IS_BASED_ON]->(catItem) 
  WITH line, f, o, itm, catItem, case when line.serialNumber <> "" then line.serialNumber else itm.serialNumber end as itmSerialNumber, case when itm.uid is null then apoc.create.uuid() else itm.uid end as itmUid
  SET itm.name = line.name, itm.serialNumber = itmSerialNumber, itm.deleted = false, itm.uid = itmUid, itm.imported_qr_code = true, itm.parentSystemImportUid = line.systemGuid, itm.supplierImport = line.manufacturer, itm.lastUpdateTime = datetime(), itm.lastUpdatedBy = "admin"
  WITH line, f, o, catItem, case when o.uid is null then apoc.create.uuid() else o.uid end as oUid
  SET o.uid = oUid, o.deleted = false, o.imported_qr_code = true, o.lastUpdateTime = datetime("2023-01-01"), o.lastUpdatedBy = "admin"
  WITH line, f, catItem, case when catItem.uid is null then apoc.create.uuid() else catItem.uid end as catItemUid
  SET catItem.uid = catItemUid, catItem.deleted = false, catItem.imported_qr_code = true, catItem.lastUpdateTime = datetime(), catItem.lastUpdatedBy = "admin"
} IN TRANSACTIONS OF 500 ROWS;

//dodat is delivered

MATCH(cat:CatalogueCategory{uid: "97598f04-948f-4da5-95b6-b2a44e0076db"})
MATCH(ci:CatalogueItem) where not (ci)-[:BELONGS_TO_CATEGORY]->() 
MERGE(ci)-[:BELONGS_TO_CATEGORY]->(cat);

MATCH(itm:Item) where not ()-[:CONTAINS_ITEM]->(itm) and itm.parentSystemImportUid is not null and itm.parentSystemImportUid <> " "
WITH itm
MATCH(parentSystem:System{uid: itm.parentSystemImportUid})
MERGE(parentSystem)-[:HAS_SUBSYSTEM]->(s:System{uid: apoc.create.uuid(), name: itm.name, deleted: false, imported_qr_code: true })-[:CONTAINS_ITEM]->(itm)
WITH s
MATCH(f:Facility{code: "B"})
MERGE(s)-[:BELONGS_TO_FACILITY]->(f);
