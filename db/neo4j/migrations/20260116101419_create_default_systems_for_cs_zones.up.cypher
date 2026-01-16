//fix root system "Control and facility systems" in ELI-Beamlines facility 
//uid: ef37cdb5-6e5e-4e71-91ce-1c62390c69d6 
//fixedSystemStructureCode: "control-and-facility-systems"
//systemLevel: "TECHNOLOGY_UNIT"
//deleted: false
MATCH(fB:Facility{code:"B"})
MERGE(cafs:System{uid:"ef37cdb5-6e5e-4e71-91ce-1c62390c69d6"})
MERGE (cafs)-[:BELONGS_TO_FACILITY]->(fB)
SET cafs.name="Control and facility systems",
    cafs.fixedSystemStructureCode="control-and-facility-systems",
    cafs.systemLevel="TECHNOLOGY_UNIT",
    cafs.deleted=false;

//fixed subsystem "CS Zones" under "Control and facility systems" in ELI-Beamlines facility
// uid: 996ea2fb-c48c-40f2-b312-2676ece39953
// code: "cs-zones"
// name: "CS Zones"
// fixedSystemStructureCode: "cs-zones"
// systemLevel: "KEY_SYSTEMS"
// deleted: false
MATCH(fB:Facility{code:"B"})
MATCH(cafs:System{uid:"ef37cdb5-6e5e-4e71-91ce-1c62390c69d6"})
MERGE(csz:System{uid:"996ea2fb-c48c-40f2-b312-2676ece39953"})
MERGE (csz)-[:BELONGS_TO_FACILITY]->(fB)
MERGE (cafs)-[:HAS_SUBSYSTEM]->(csz)
SET
    csz.name="CS Zones",
    csz.fixedSystemStructureCode="cs-zones",
    csz.systemLevel="KEY_SYSTEMS",
    csz.deleted=false;

// now for each root Zone create default System under "CS Zones" subsystem
// root zones are: match(zone:Zone)-[:BELONGS_TO_FACILITY]->(:Facility{code:"B"}) WHERE NOT ()-[:HAS_SUBZONE]->(zone)
// name: zone.code + " " + zone.name
// uid: generate a new random UUID for each System if not exists
// systemLevel: "SUBSYSTEMS_AND_PARTS"
// fixedSystemStructureCode: "cs-zone" + "-" + zone.code
// deleted: false
// belongs to Facility ELI-Beamlines
// join to parent subsystem "CS Zones"
// create realationship between Zone and System: "HAS_DEFAULT_PARENT_SYSTEM"
MATCH(fB:Facility{code:"B"})
MATCH(csz:System{uid:"996ea2fb-c48c-40f2-b312-2676ece39953"})
MATCH(zone:Zone)-[:BELONGS_TO_FACILITY]->(fB)
WHERE NOT ()-[:HAS_SUBZONE]->(zone)
MERGE (sys:System{fixedSystemStructureCode:"cs-zone-" + zone.code})
ON CREATE SET sys.uid = apoc.create.uuid()
MERGE (sys)-[:BELONGS_TO_FACILITY]->(fB)
MERGE (csz)-[:HAS_SUBSYSTEM]->(sys)
MERGE (zone)-[:HAS_DEFAULT_PARENT_SYSTEM]->(sys)
SET sys.fixedSystemStructureCode = "cs-zone-" + zone.code,
    sys.name = zone.code + " - " + zone.name,
    sys.systemLevel = "SUBSYSTEMS_AND_PARTS",
    sys.deleted = false;
