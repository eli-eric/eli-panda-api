:auto 
LOAD CSV WITH HEADERS FROM 'file:///var/lib/neo4j/import/employees-bm.csv' AS line
 WITH line
CALL {
  with line
  MATCH(adminUser:User{username: "admin"})
  with line, adminUser
  MATCH(f:Facility{code: "B"}) 
  with line, f, adminUser  
  CREATE (empl:Employee {
  uid: apoc.create.uuid(),
  firstName: line.firstname,
  lastName: line.lastname,   
  employeeNumber: line.employeeNumber,
  deleted: false,
  lastUpdatedTime: datetime(), 
  lastUpdatedBy: adminUser.username })
  with empl,adminUser, f
  MERGE(empl)-[:WAS_UPDATED_BY{at: datetime()}]->(adminUser)
  MERGE(empl)-[:AFFILIATED_WITH_FACILITY]->(f)  
} IN TRANSACTIONS OF 500 ROWS