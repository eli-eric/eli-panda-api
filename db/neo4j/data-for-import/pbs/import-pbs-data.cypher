:auto LOAD CSV WITH HEADERS FROM 'file:///var/lib/neo4j/import/pbs_eun_data.csv' AS line
  WITH line 
  CALL {
    WITH line
    CREATE(pbs:PbsData{
      eun: line.eun, 
      name: line.name,
      pbsVerso: line.pbsVerso,
      pbsFis: line.pbsFis,
      pbsTenderReference: line.pbsTenderReference,
      pbsSupplier: line.pbsSupplier,
      cataloguePartNumber: line.cataloguePartNumber,
      itemPriceCzk: line.itemPriceCzk,
      itemPriceEUR: line.itemPriceEUR,
      pbsNumber: line.pbsNumber,
      cartName: line.cartName,
      sectionName: line.sectionName,
      quantity: line.quantity,
      filesCount: line.filesCount
      })
  } IN TRANSACTIONS OF 500 ROWS;
  