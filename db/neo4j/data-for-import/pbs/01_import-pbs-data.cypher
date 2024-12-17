// this is the main import script for PBS data
// use this script to import PBS data into the system
// example of the import file is in the file: db/neo4j/data-for-import/pbs/example.csv
:auto LOAD CSV WITH HEADERS FROM 'file:///var/lib/neo4j/import/chnagethefilename.csv' AS line
with line
where line.ImportInstr = "Yes" 
and line.pandaOrderGUID is not null
and (line.cataloguePartNumber is not null or line.PandaPartNumber is not null) 
and line.eun is not null  
with line
  CALL {
  WITH line
  MERGE(o:Order{ uid: trim(line.pandaOrderGUID) })  
  MERGE(itm:Item{ eun: trim(line.eun) })   
  MERGE(o)-[ol:HAS_ORDER_LINE]->(itm)
  MERGE(catItem:CatalogueItem{catalogueNumber: coalesce(trim(line.PandaPartNumber), trim(line.cataloguePartNumber))})
  MERGE(itm)-[:IS_BASED_ON]->(catItem) 
  WITH line, o, ol, itm, catItem,
  case when line.pbsTenderReference is not null then "
  Tender reference imported from PBS: " + line.pbsTenderReference else "" end as tenderReferenceNote,
  case when line.pbsSupplier is not null then "
  PBS supplier: " + line.pbsSupplier else "" end as supplierNote,
  case when line.notes is not null then "
  PBS notes: " + line.notes else "" end as notesNote,
  case when line.description is not null then "
  PBS description: " + line.description else "" end as descriptionNote,
  case when line.sectionName is not null then "
  PBS section name: " + line.sectionName else "" end as sectionNameNote,
  case when line.destination is not null then "
  PBS destination: " + line.destination else "" end as destinationNote,
  case when line.cartName is not null then "
  PBS cart name: " + line.cartName else "" end as cartNameNote,
  case when line.itemPriceCZK is not null then "CZK" else null end as currencyCZK,
  case when line.itemPriceEUR is not null and line.itemPriceEUR <> "0" then "EUR" else "CZK" end as currencyEUR,
  case when line.pbsNumber is not null then "
  PBS number: " + line.pbsNumber else "" end as pbsNumberNote
  SET 
    ol.isDelivered = true,
    ol.deliveredTimeImport = line.deliveredDate,
    ol.touch_by_import = true,
    ol.price = case when currencyEUR is not null then toFloat(line.itemPriceEUR) else toFloat(line.itemPriceCZK) end,
    ol.currency = coalesce(currencyEUR, currencyCZK),
    itm.name = coalesce(itm.name,catItem.name, line.name),
    itm.serialNumber = coalesce(itm.serialNumber, line.serialNumber),
    itm.deleted = false, 
    itm.uid = coalesce(itm.uid, apoc.create.uuid()),
    itm.touch_by_import = true,
    itm.parentSystemImportUid = line.TargetsystemGUID, 
    itm.itemUsageImport = line.ItemUsage,
    itm.itemConditionImport = line.ItemCondition,
    itm.supplierImport = line.pbsSupplier, 
    itm.lastUpdateTime = datetime(), 
    itm.lastUpdatedBy = "admin",  
    itm.importedNotes = pbsNumberNote + supplierNote + notesNote + descriptionNote + destinationNote + cartNameNote, 
    itm.hasImageImport = line.hasImage,
    itm.filesCountImport = line.filesCount,
    itm.quantityImport = toInteger(line.quantity),    
    o.touch_by_import = true,     
    o.lastUpdatedBy = "admin",
    o.importedNotes =  tenderReferenceNote,
    o.orderNumber = coalesce(o.orderNumber, line.pbsFis),
    o.requestNumber = coalesce(o.requestNumber, line.pbsVerso),
    catItem.name = coalesce(catItem.name, coalesce(line.PandaPartNumber, line.cataloguePartNumber)),
    itm.name = coalesce(itm.name,catItem.name, line.name),
    catItem.uid = coalesce(catItem.uid, apoc.create.uuid()),
    catItem.catalogueNumber = coalesce(catItem.catalogueNumber, coalesce(line.PandaPartNumber, line.cataloguePartNumber)),
    catItem.deleted = false, 
    catItem.touch_by_import = true, 
    catItem.lastUpdateTime = datetime(), 
    catItem.lastUpdatedBy = "admin"

  } IN TRANSACTIONS OF 200 ROWS;
  

  // next phase is to setup System if exists parentSystemImportUid on the Item and no System is connected to the Item
MATCH(itm:Item) where not ()-[:CONTAINS_ITEM]->(itm) and itm.parentSystemImportUid is not null and itm.parentSystemImportUid <> " "
WITH itm
MATCH(parentSystem:System{uid: itm.parentSystemImportUid})
MERGE(parentSystem)-[:HAS_SUBSYSTEM]->(s:System)-[:CONTAINS_ITEM]->(itm)
SET 
s.uid = coalesce(s.uid, apoc.create.uuid()),
s.name = itm.name,
s.deleted = false,
s.systemLevel = coalesce(s.systemLevel, "SUBSYSTEMS_AND_PARTS"),
s.touch_by_import = true 
WITH s
MATCH(f:Facility{code: "B"})
MERGE(s)-[:BELONGS_TO_FACILITY]->(f);
  
// setup ItemUsage and ItemCondition and is not connected to the Item
MATCH(itm:Item) where itm.itemUsageImport is not null and not (itm)-[:HAS_ITEM_USAGE]->()
WITH itm
MATCH(iu:ItemUsage{name: itm.itemUsageImport})
MERGE(itm)-[:HAS_ITEM_USAGE]->(iu);

MATCH(itm:Item) where itm.itemConditionImport is not null and not (itm)-[:HAS_CONDITION_STATUS]->()
WITH itm
MATCH(ic:ItemCondition{name: itm.itemConditionImport})
MERGE(itm)-[:HAS_CONDITION_STATUS]->(ic);

// setup ItemUsage as a "Stock Item" if itemUsageImport is null and was touched by import
MATCH(itm:Item) where itm.touch_by_import = true and not (itm)-[:HAS_ITEM_USAGE]->()
WITH itm
MATCH(iu:ItemUsage{name: "Stock Item"})
MERGE(itm)-[:HAS_ITEM_USAGE]->(iu);

//connect catalogue categor yto catalogue item if no connection exists
MATCH(cat:CatalogueCategory{uid: "97598f04-948f-4da5-95b6-b2a44e0076db"})
MATCH(ci:CatalogueItem) where not (ci)-[:BELONGS_TO_CATEGORY]->()
MERGE(ci)-[:BELONGS_TO_CATEGORY]->(cat);

// match all items and order and set notes from importedNotes to notes
MATCH(itm:Item) WHERE itm.importedNotes is not null and itm.importedNotes <> ""
with itm, case when itm.notes is null then "" else itm.notes end as notes
SET itm.notes = notes + itm.importedNotes;

MATCH(o:Order) WHERE o.importedNotes is not null and o.importedNotes <> ""
with o, case when o.notes is null then "" else o.notes end as notes
SET o.notes = notes + o.importedNotes;

match(itm)-[r:IS_BASED_ON]->(catItem) 
with itm, count(r) as counts where counts > 1
match(itm)-[r:IS_BASED_ON]->(catItem{touch_by_import: true}) 
delete r;

// set order delivery statuses
  MATCH (o:Order{touch_by_import: true})
	WITH o
	MATCH(o)-[olAll:HAS_ORDER_LINE]->()
	WITH count(olAll) as totalLines, o
	OPTIONAL MATCH(o)-[olDelivered:HAS_ORDER_LINE{isDelivered: true}]->()
	WITH totalLines, count(olDelivered) as deliveredLines, o
	SET o.deliveryStatus = case when deliveredLines = 0 then 0 when deliveredLines = totalLines then 2 else 1 end;

  // set order names if not set
  MATCH(o:Order) WHERE o.name IS NULL and o.orderNumber IS NOT NULL SET o.name = "Order " + o.orderNumber;