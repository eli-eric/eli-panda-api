:auto 
LOAD CSV WITH HEADERS FROM 'file:///var/lib/neo4j/import/test-orders.csv' AS line
 WITH line where line.codeTA <> "100" and line.supplier <> "Cestovni prikaz"
CALL {
  with line
  MATCH(adminUser:User{username: "admin"})
  with line, adminUser
  MATCH(f:Facility{code: "B"}) 
  with line, f, adminUser
  OPTIONAL MATCH(s:Supplier{shortName: line.supplier})
  with line,s, adminUser
  OPTIONAL MATCH(os:OrderStatus{code: line.orderStatus}) 
  with line, s, os, adminUser
  CREATE (o:Order {uid: apoc.create.uuid(),
  name: line.notes,
  orderNumber: line.orderNumber,   
  orderDate: datetime(line.orderDate),
  contractNumber: line.contractNumber,  
  createdBy: line.createdBy,  
  deleted: false,
  testRecord: true})
  with line,s,o, os, adminUser, f
  MERGE(o)-[:UPDATED_BY{updated: datetime()}]->(adminUser)
  MERGE(o)-[:BELONGS_TO_FACILITY]->(f)
  MERGE(o)-[:HAS_SUPPLIER]->(s)
  MERGE(o)-[:HAS_ORDER_STATUS]->(os)
} IN TRANSACTIONS OF 500 ROWS