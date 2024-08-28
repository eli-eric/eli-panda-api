:auto LOAD CSV WITH HEADERS FROM 'file:///var/lib/neo4j/import/2024_08_27_PBS_import.csv' AS line
with line
where line.ImportInstr = "Yes" 
and line.pandaOrderGUID is not null
and (line.cataloguePartNumber is not null or line.PandaPartNumber is not null) 
and line.eun is not null  
with line
  CALL {
  WITH line
  MERGE(o:Order{ uid: line.pandaOrderGUID })
  MERGE(catItem:CatalogueItem{catalogueNumber: coalesce(line.PandaPartNumber, line.cataloguePartNumber) })
  MERGE(itm:Item{ eun: line.eun }) 
  MERGE(o)-[ol:HAS_ORDER_LINE]->(itm)
  MERGE(itm)-[:IS_BASED_ON]->(catItem) 
  WITH line, o, ol, itm, catItem,
  case when line.pbsTenderReference is not null then "
  PBS tender reference: " + line.pbsTenderReference else "" end as tenderReferenceNote,
  case when line.pbsSupplier is not null then "
  PBS supplier: " + line.pbsSupplier else "" end as supplierNote,
  case when line.notes is not null then "
  PBS notes: " + line.notes else "" end as notesNote,
  case when line.description is not null then "
  PBS description: " + line.description else "" end as descriptionNote,
  case when line.sectionName is not null then "
  PBS section name: " + line.sectionName else "" end as sectionNameNote,
  case when line.destionation is not null then "
  PBS destination: " + line.destionation else "" end as destinationNote,
  case when line.cartName is not null then "
  PBS cart name: " + line.cartName else "" end as cartNameNote,
  case when line.itemPriceCZK is not null then "CZK" else null end as currencyCZK,
  case when line.itemPriceEUR is not null then  "EUR" else null end as currencyEUR,
  case when line.pbsNumber is not null then "
  PBS number: " + line.pbsNumber else "" end as pbsNumberNote,
  case when itm.notes is not null then itm.notes else "" end as itmNotes,
  case when o.description is not null then o.description else "" end as oDescription
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
    itm.parentSystemImportUid = line.TargetSystemGUID, 
    itm.itemUsageImport = line.ItemUsage,
    itm.itemConditionImport = line.ItemCondition,
    itm.supplierImport = line.pbsSupplier, 
    itm.lastUpdateTime = datetime(), 
    itm.lastUpdatedBy = "admin",  
    itm.notes = pbsNumberNote + supplierNote + notesNote + descriptionNote + sectionNameNote + destinationNote + cartNameNote,
    itm.hasImageImport = line.hasImage,
    itm.filesCountImport = line.filesCount,
    itm.quantityImport = toInteger(line.quantity),    
    o.touch_by_import = true,     
    o.lastUpdatedBy = "admin",
    o.description = oDescription + tenderReferenceNote,
    o.orderNumber = coalesce(o.orderNumber, line.pbsFis),
    o.requestNumber = coalesce(o.requestNumber, line.pbsVerso),
    catItem.name = coalesce(catItem.name, coalesce(line.PandaPartNumber, line.cataloguePartNumber)),
    catItem.uid = coalesce(catItem.uid, apoc.create.uuid()),
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