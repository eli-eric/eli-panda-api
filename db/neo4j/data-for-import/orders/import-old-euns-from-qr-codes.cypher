:auto LOAD CSV WITH HEADERS FROM 'file:///var/lib/neo4j/import/qr_codes_euns_24_10_2024.csv' AS line
  WITH line where tolower(line.import) = "y" and line.orderId is not null and line.catalogueNumber is not null and line.eun is not null
CALL {
  WITH line
  MATCH(f:Facility{code: "B"})     
  WITH line, f
  MERGE(o:Order{ orderNumber:trim(line.orderId)})
  MERGE(catItem:CatalogueItem{catalogueNumber: trim(line.catalogueNumber) })
  MERGE(itm:Item{ eun: trim(line.eun) })  
  MERGE(o)-[:BELONGS_TO_FACILITY]->(f)
  MERGE(o)-[ol:HAS_ORDER_LINE]->(itm)
  WITH line, f, o, ol, itm, catItem limit 1 
  MERGE(itm)-[:IS_BASED_ON]->(catItem) 
  WITH line, f, o,ol, itm, catItem, case when line.serialNumber <> "" then trim(line.serialNumber) else itm.serialNumber end as itmSerialNumber, case when itm.uid is null then apoc.create.uuid() else itm.uid end as itmUid
  SET 
    ol.isDelivered = true,
    itm.name = case when catItem.name = "" or catItem.name is null then line.name else catItem.name end, 
    itm.serialNumber = itmSerialNumber,
    itm.deleted = false, 
    itm.uid = itmUid, 
    itm.imported_qr_code = true,
    itm.parentSystemImportUid = trim(line.systemGuid), 
    itm.supplierImport = trim(line.manufacturer), 
    itm.lastUpdateTime = datetime(), 
    itm.lastUpdatedBy = "admin"
  WITH line, f, o, catItem, case when o.uid is null then apoc.create.uuid() else o.uid end as oUid, case when o.name is null then "Order " + line.orderId else o.name end as oName
  SET 
  o.uid = oUid, 
  o.deleted = false, 
  o.name = oName,
  o.imported_qr_code = true, 
  o.lastUpdateTime = datetime("2023-01-01"), 
  o.lastUpdatedBy = "admin"
  WITH line, f, catItem, case when catItem.uid is null then apoc.create.uuid() else catItem.uid end as catItemUid
  SET 
  catItem.name = case when catItem.name = "" or catItem.name is null then line.name else catItem.name end,
  catItem.uid = catItemUid, 
  catItem.deleted = false, 
  catItem.imported_qr_code = true, 
  catItem.lastUpdateTime = datetime(), 
  catItem.lastUpdatedBy = "admin"
} IN TRANSACTIONS OF 500 ROWS;

//dodat is delivered

MATCH(cat:CatalogueCategory{uid: "97598f04-948f-4da5-95b6-b2a44e0076db"})
MATCH(ci:CatalogueItem) where not (ci)-[:BELONGS_TO_CATEGORY]->() 
MERGE(ci)-[:BELONGS_TO_CATEGORY]->(cat);

MATCH(itm:Item) where not ()-[:CONTAINS_ITEM]->(itm) and itm.parentSystemImportUid is not null and itm.parentSystemImportUid <> " "
WITH itm
MATCH(parentSystem:System{uid: itm.parentSystemImportUid})
MERGE(parentSystem)-[:HAS_SUBSYSTEM]->(s:System)-[:CONTAINS_ITEM]->(itm)
SET 
s.uid = case when s.uid is null then apoc.create.uuid() else s.uid end,
s.name = itm.name,
s.deleted = false,
s.systemLevel = coalesce(s.systemLevel, "SUBSYSTEMS_AND_PARTS"),
s.imported_qr_code = true 
WITH s
MATCH(f:Facility{code: "B"})
MERGE(s)-[:BELONGS_TO_FACILITY]->(f)
WITH s
MATCH(u:User{username: "admin"})
CREATE(s)-[:WAS_UPDATED_BY{ at: datetime(), action: "IMPORT" }]->(u);

// set order delivery statuses
  MATCH (o:Order{imported_qr_code: true})
	WITH o
	MATCH(o)-[olAll:HAS_ORDER_LINE]->()
	WITH count(olAll) as totalLines, o
	OPTIONAL MATCH(o)-[olDelivered:HAS_ORDER_LINE{isDelivered: true}]->()
	WITH totalLines, count(olDelivered) as deliveredLines, o
	SET o.deliveryStatus = case when deliveredLines = 0 then 0 when deliveredLines = totalLines then 2 else 1 end;

// setup ItemUsage as a "Stock Item" if itemUsageImport is null and was touched by import
MATCH(itm:Item) where itm.imported_qr_code = true and not (itm)-[:HAS_ITEM_USAGE]->()
WITH itm
MATCH(iu:ItemUsage{name: "Stock Item"})
MERGE(itm)-[:HAS_ITEM_USAGE]->(iu);

// set order names if not set
MATCH(o:Order) WHERE o.name IS NULL and o.orderNumber IS NOT NULL SET o.name = "Order " + o.orderNumber;