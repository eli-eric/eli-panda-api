:auto LOAD CSV WITH HEADERS FROM 'file:///var/lib/neo4j/import/pbs_eun_data.csv' AS line
  WITH line 
  CALL {
    WITH line
    OPTIONAL MATCH (itm:Item) WHERE toLower(trim(itm.eun)) = toLower(trim(line.eun))
    OPTIONAL MATCH(o:Order) WHERE toLower(trim(o.orderNumber)) = toLower(trim(line.pbsFis)) OR toLower(trim(o.orderNumber)) = toLower(trim(line.pbsVerso))
    CREATE(pbs:PbsData{
      eun: line.eun, 
      name: line.name,
      pbsVerso: line.pbsVerso,
      pbsFis: line.pbsFis,
      pbsTenderReference: line.pbsTenderReference,
      pbsSupplier: line.pbsSupplier,
      cataloguePartNumber: line.cataloguePartNumber,
      serialNumber: line.serialNumber,
      itemPriceCzk: line.itemPriceCzk,
      itemPriceEUR: line.itemPriceEUR,
      pbsNumber: line.pbsNumber,
      cartName: line.cartName,
      sectionName: line.sectionName,
      quantity: line.quantity,
      filesCount: line.filesCount,
      hasImage: line.hasImage,
      deliveredDate: line.deliveredDate,
      description: line.description,
      notes: line.notes,
      destination: line.destination,
      isEunInPanda: CASE WHEN itm IS NOT NULL THEN true ELSE false END,
      isOrderInPanda: CASE WHEN o IS NOT NULL THEN true ELSE false END,
      pandaOrderGUID: CASE WHEN o IS NOT NULL THEN o.uid ELSE null END
      })
  } IN TRANSACTIONS OF 200 ROWS;
  
  