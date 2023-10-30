:auto 
LOAD CSV WITH HEADERS FROM 'file:///var/lib/neo4j/import/qr_codes_import.csv' AS line
  WITH line where line.import = "y" and line.orderId is not null
CALL {
  WITH line
  MATCH(f:Facility{code: "B"})   
  MATCH(cat:CatalogueCategory{uid: "97598f04-948f-4da5-95b6-b2a44e0076db"})
  WITH line, f
  MERGE(o:Order{ orderNumber: line.orderId})
  MERGE(catItem:CatalogueItem{catalogueNumber: line.catalogueNumber, name: line.name })
  MERGE(itm:Item{ eun: line.eun })  
  MERGE(o)-[:BELONGS_TO_FACILITY]->(f)
  MERGE(o)-[:HAS_ORDER_LINE]->(itm)
  MERGE(itm)-[:IS_BASED_ON]->(catItem)
  WITH line, f, o, itm, catItem, case when line.serialNumber <> "" then line.serialNumber else itm.serialNumber end as itmSerialNumber, case when itm.uid is null then apoc.create.uuid() else itm.uid end as itmUid
  SET itm.name = line.name, itm.serialNumber = itmSerialNumber, itm.deleted = false, itm.uid = itmUid, itm.imported_qr_code = true
  WITH line, f, o, catItem, case when o.uid is null then apoc.create.uuid() else o.uid end as oUid
  SET o.uid = oUid, o.deleted = false, o.imported_qr_code = true
  WITH line, f, catItem, case when catItem.uid is null then apoc.create.uuid() else catItem.uid end as catItemUid
  SET catItem.uid = catItemUid, catItem.deleted = false, catItem.imported_qr_code = true
} IN TRANSACTIONS OF 500 ROWS