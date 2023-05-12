:auto 
LOAD CSV WITH HEADERS FROM 'file:///var/lib/neo4j/import/suppliers.csv' AS line
 WITH line
CALL {  
  WITH line   
  CREATE (o:Order {uid: apoc.create.uuid(),
  name: line.Name,
  shortName: line.ShortName,   
  address: line.Address,
  stateCode: line.StateCode,
  CIN: line.CIN,
  VAT: line.VAT,
  deleted: false})
} IN TRANSACTIONS OF 500 ROWS;