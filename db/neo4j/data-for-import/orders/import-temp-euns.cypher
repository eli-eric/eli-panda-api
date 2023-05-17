LOAD CSV WITH HEADERS FROM 'file:///var/lib/neo4j/import/euns-orders-old.csv' AS line
 WITH line
  MATCH(adminUser:User{username: "admin"})
  with line, adminUser
  MATCH(f:Facility{code: "B"}) 
  with line, f, adminUser  
  MATCH(ccg:CatalogueCategory{uid: '97598f04-948f-4da5-95b6-b2a44e0076db'})
  with line, f, adminUser, ccg
  OPTIONAL MATCH(os:OrderStatus{code: 'ORDERED'}) 
  with line, os, adminUser, ccg, f, apoc.text.split(line.deliveryTime, '/') as dateSplit
  with line, os, adminUser, ccg, f, datetime({year: toInteger(dateSplit[2]),month: toInteger(dateSplit[0]),day: toInteger(dateSplit[1]), hour: 12}) as deliveryTime
  MERGE (o:Order { 
    imported: true,
    deleted:false,
    deliveryStatus: 2,
    name: trim(line.orderNumber) + ' - imported',  
    orderNumber: trim(line.orderNumber), 
    orderDate: datetime({year: 2023, month: 5, day: 17, hour: 12}),
    lastUpdateTime: datetime({year: 2023, month: 5, day: 17, hour: 12}),
    lastUpdateBy: 'admin',
    notes: 'Imported from temporary excel sheet.'})
 with line, o, os, adminUser, f,deliveryTime, ccg
 set o.uid = apoc.create.uuid() 
with line,o, os, adminUser, f, deliveryTime, ccg,
case when line.catalogueNumber is null or line.catalogueNumber = 'NA' then trim(line.name) else trim(line.catalogueNumber) end as catalogueNumber,
case when line.serialNumber is null or line.serialNumber = 'NA' then '' else trim(line.serialNumber) end as serialNumber
MERGE(ci:CatalogueItem{catalogueNumber: catalogueNumber, name: trim(line.name), lastUpdateBy: 'admin', imported: true})-[:BELONGS_TO_CATEGORY]->(ccg)
with line,o, os, adminUser, f, ci,deliveryTime, serialNumber
set ci.uid = apoc.create.uuid()
with line,o, os, adminUser, f, ci,deliveryTime, serialNumber
  MERGE (o)-[ol:HAS_ORDER_LINE{ imported:true, isDelivered:true, deliveredTime: deliveryTime, lastUpdateBy: 'admin' }]->(itm:Item{
    imported: true, 
    eun: 'B' + trim(line.eun),
    name: trim(line.name),
    serialNumber: serialNumber   
   })-[:IS_BASED_ON]->(ci)
  with line, o, os, adminUser, f, itm
  set itm.uid = apoc.create.uuid()
  with line, o, os, adminUser, f
  MERGE(o)-[:WAS_UPDATED_BY{ at: datetime(), action: "INSERT" }]->(adminUser)
  MERGE(o)-[:BELONGS_TO_FACILITY]->(f)
  MERGE(o)-[:HAS_ORDER_STATUS]->(os);
