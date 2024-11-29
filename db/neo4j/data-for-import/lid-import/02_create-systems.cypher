MATCH (lid:LidDevice), (sys:System)
WHERE sys.systemCode = lid.code
WITH collect(lid.code) as existingCodes
MATCH (lid:LidDevice) where not lid.code in existingCodes
WITH lid
MATCH (parent:System{uid:"3558f48b-36f2-4646-90ad-b280b61a09a4"})
MATCH(f:Facility{code: "B"})
MERGE(s:System{systemCode: lid.pandaCode})
SET 
s.uid = coalesce(s.uid, apoc.create.uuid()), 
s.name = lid.pandaCode, 
s.systemLevel = "SUBSYSTEMS_AND_PARTS",
s.description = lid.description, 
s.deleted = false,
s.importedFromLID = true,
s.lidZone = lid.zone,
s.lidType = lid.type,
s.lidSerial = lid.serial,
s.lidCode = lid.code
MERGE(parent)-[:HAS_SUBSYSTEM]->(s)
MERGE(s)-[:BELONGS_TO_FACILITY]->(f);

// set HAS_ZONE relationship
MATCH (s:System) where s.lidZone is not null
MATCH (z:Zone{code: apoc.text.lpad(s.lidZone,2,'0')})
MERGE (s)-[:HAS_ZONE]->(z);

// set HAS_SYSTEM_TYPE relationship
MATCH (s:System) where s.lidType is not null
MATCH (st:SystemType{code: s.lidType})
MERGE (s)-[:HAS_SYSTEM_TYPE]->(st);

// set was_updated_by relationship
MATCH (s:System) where s.importedFromLID = true
MATCH(u:User{username: "admin"})
CREATE(s)-[:WAS_UPDATED_BY{ at: datetime(), action: "IMPORT" }]->(u);

// add old code to the description
MATCH (s:System) where s.importedFromLID = true
SET s.description = coalesce(s.description,"") + " (old code: " + s.lidCode + ")";