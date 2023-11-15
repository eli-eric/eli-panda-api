:auto LOAD CSV WITH HEADERS FROM 'file:///var/lib/neo4j/import/np-employees.csv' AS line
with line
match(f:Facility{code:"N"})
with split(trim(line.name)," ") as names, line, f
with line,f, case when size(names) = 2 then names[0] when size(names) = 3 then names[1] + " " + names[2] end as firstName, case when size(names) = 2 then names[1] when size(names) = 3 then names[0] end as lastName
MERGE(empl:Employee{employeeNumber: line.email})
set 
empl.uid = apoc.create.uuid(), 
empl.firstName = firstName,
empl.lastName = lastName,
empl.deleted = false,
empl.fullName = lastName + " " + firstName,
empl.lastUpdateBy = "admin",
empl.phone = line.phone,
empl.npGroup = line.group,
empl.position = line.position
MERGE(empl)-[:AFFILIATED_WITH_FACILITY]->(f)
