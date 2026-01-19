// Create ExperimentalSystem constraint and index
CREATE CONSTRAINT experimentalsystem_uid_unique IF NOT EXISTS FOR (n:ExperimentalSystem) REQUIRE n.uid IS UNIQUE;
CREATE INDEX experimentalsystem_name_index IF NOT EXISTS FOR (n:ExperimentalSystem) ON (n.name);

// Seed data for ELI-Beamlines facility (28 entries)
MATCH (f:Facility {code: "B"})
MERGE (es:ExperimentalSystem {code: "ELIps"}) ON CREATE SET es.uid = apoc.create.uuid(), es.name = "ELIps"
MERGE (es)-[:BELONGS_TO_FACILITY]->(f);

MATCH (f:Facility {code: "B"})
MERGE (es:ExperimentalSystem {code: "ALFA"}) ON CREATE SET es.uid = apoc.create.uuid(), es.name = "ALFA"
MERGE (es)-[:BELONGS_TO_FACILITY]->(f);

MATCH (f:Facility {code: "B"})
MERGE (es:ExperimentalSystem {code: "ELBA"}) ON CREATE SET es.uid = apoc.create.uuid(), es.name = "ELBA"
MERGE (es)-[:BELONGS_TO_FACILITY]->(f);

MATCH (f:Facility {code: "B"})
MERGE (es:ExperimentalSystem {code: "ELIMAIA"}) ON CREATE SET es.uid = apoc.create.uuid(), es.name = "ELIMAIA"
MERGE (es)-[:BELONGS_TO_FACILITY]->(f);

MATCH (f:Facility {code: "B"})
MERGE (es:ExperimentalSystem {code: "ELIMAIA-ELIMED"}) ON CREATE SET es.uid = apoc.create.uuid(), es.name = "ELIMAIA-ELIMED"
MERGE (es)-[:BELONGS_TO_FACILITY]->(f);

MATCH (f:Facility {code: "B"})
MERGE (es:ExperimentalSystem {code: "FSRS"}) ON CREATE SET es.uid = apoc.create.uuid(), es.name = "FSRS"
MERGE (es)-[:BELONGS_TO_FACILITY]->(f);

MATCH (f:Facility {code: "B"})
MERGE (es:ExperimentalSystem {code: "FSRS & TA"}) ON CREATE SET es.uid = apoc.create.uuid(), es.name = "FSRS & TA"
MERGE (es)-[:BELONGS_TO_FACILITY]->(f);

MATCH (f:Facility {code: "B"})
MERGE (es:ExperimentalSystem {code: "Gammatron"}) ON CREATE SET es.uid = apoc.create.uuid(), es.name = "Gammatron"
MERGE (es)-[:BELONGS_TO_FACILITY]->(f);

MATCH (f:Facility {code: "B"})
MERGE (es:ExperimentalSystem {code: "HHG"}) ON CREATE SET es.uid = apoc.create.uuid(), es.name = "HHG"
MERGE (es)-[:BELONGS_TO_FACILITY]->(f);

MATCH (f:Facility {code: "B"})
MERGE (es:ExperimentalSystem {code: "HHG-MAC"}) ON CREATE SET es.uid = apoc.create.uuid(), es.name = "HHG-MAC"
MERGE (es)-[:BELONGS_TO_FACILITY]->(f);

MATCH (f:Facility {code: "B"})
MERGE (es:ExperimentalSystem {code: "L1 - ALFA"}) ON CREATE SET es.uid = apoc.create.uuid(), es.name = "L1 - ALFA"
MERGE (es)-[:BELONGS_TO_FACILITY]->(f);

MATCH (f:Facility {code: "B"})
MERGE (es:ExperimentalSystem {code: "L1 - HHG"}) ON CREATE SET es.uid = apoc.create.uuid(), es.name = "L1 - HHG"
MERGE (es)-[:BELONGS_TO_FACILITY]->(f);

MATCH (f:Facility {code: "B"})
MERGE (es:ExperimentalSystem {code: "L1 - HHG - MAC"}) ON CREATE SET es.uid = apoc.create.uuid(), es.name = "L1 - HHG - MAC"
MERGE (es)-[:BELONGS_TO_FACILITY]->(f);

MATCH (f:Facility {code: "B"})
MERGE (es:ExperimentalSystem {code: "L1 - PXS"}) ON CREATE SET es.uid = apoc.create.uuid(), es.name = "L1 - PXS"
MERGE (es)-[:BELONGS_TO_FACILITY]->(f);

MATCH (f:Facility {code: "B"})
MERGE (es:ExperimentalSystem {code: "L1 - PXS - TREX"}) ON CREATE SET es.uid = apoc.create.uuid(), es.name = "L1 - PXS - TREX"
MERGE (es)-[:BELONGS_TO_FACILITY]->(f);

MATCH (f:Facility {code: "B"})
MERGE (es:ExperimentalSystem {code: "L3 - ELBA"}) ON CREATE SET es.uid = apoc.create.uuid(), es.name = "L3 - ELBA"
MERGE (es)-[:BELONGS_TO_FACILITY]->(f);

MATCH (f:Facility {code: "B"})
MERGE (es:ExperimentalSystem {code: "L3 - ELIMAIA"}) ON CREATE SET es.uid = apoc.create.uuid(), es.name = "L3 - ELIMAIA"
MERGE (es)-[:BELONGS_TO_FACILITY]->(f);

MATCH (f:Facility {code: "B"})
MERGE (es:ExperimentalSystem {code: "L3 - ELIMAIA - ELIMED"}) ON CREATE SET es.uid = apoc.create.uuid(), es.name = "L3 - ELIMAIA - ELIMED"
MERGE (es)-[:BELONGS_TO_FACILITY]->(f);

MATCH (f:Facility {code: "B"})
MERGE (es:ExperimentalSystem {code: "L4 - P3"}) ON CREATE SET es.uid = apoc.create.uuid(), es.name = "L4 - P3"
MERGE (es)-[:BELONGS_TO_FACILITY]->(f);

MATCH (f:Facility {code: "B"})
MERGE (es:ExperimentalSystem {code: "L4n - P3"}) ON CREATE SET es.uid = apoc.create.uuid(), es.name = "L4n - P3"
MERGE (es)-[:BELONGS_TO_FACILITY]->(f);

MATCH (f:Facility {code: "B"})
MERGE (es:ExperimentalSystem {code: "Legend - HHG - MAC"}) ON CREATE SET es.uid = apoc.create.uuid(), es.name = "Legend - HHG - MAC"
MERGE (es)-[:BELONGS_TO_FACILITY]->(f);

MATCH (f:Facility {code: "B"})
MERGE (es:ExperimentalSystem {code: "Legend - PXS - TREX"}) ON CREATE SET es.uid = apoc.create.uuid(), es.name = "Legend - PXS - TREX"
MERGE (es)-[:BELONGS_TO_FACILITY]->(f);

MATCH (f:Facility {code: "B"})
MERGE (es:ExperimentalSystem {code: "MAC"}) ON CREATE SET es.uid = apoc.create.uuid(), es.name = "MAC"
MERGE (es)-[:BELONGS_TO_FACILITY]->(f);

MATCH (f:Facility {code: "B"})
MERGE (es:ExperimentalSystem {code: "PXS - TREX"}) ON CREATE SET es.uid = apoc.create.uuid(), es.name = "PXS - TREX"
MERGE (es)-[:BELONGS_TO_FACILITY]->(f);

MATCH (f:Facility {code: "B"})
MERGE (es:ExperimentalSystem {code: "TA"}) ON CREATE SET es.uid = apoc.create.uuid(), es.name = "TA"
MERGE (es)-[:BELONGS_TO_FACILITY]->(f);

MATCH (f:Facility {code: "B"})
MERGE (es:ExperimentalSystem {code: "TCT"}) ON CREATE SET es.uid = apoc.create.uuid(), es.name = "TCT"
MERGE (es)-[:BELONGS_TO_FACILITY]->(f);

MATCH (f:Facility {code: "B"})
MERGE (es:ExperimentalSystem {code: "trELIps"}) ON CREATE SET es.uid = apoc.create.uuid(), es.name = "trELIps"
MERGE (es)-[:BELONGS_TO_FACILITY]->(f);

MATCH (f:Facility {code: "B"})
MERGE (es:ExperimentalSystem {code: "trELIPs & TCT"}) ON CREATE SET es.uid = apoc.create.uuid(), es.name = "trELIPs & TCT"
MERGE (es)-[:BELONGS_TO_FACILITY]->(f);

MATCH (f:Facility {code: "B"})
MERGE (es:ExperimentalSystem {code: "TREX"}) ON CREATE SET es.uid = apoc.create.uuid(), es.name = "TREX"
MERGE (es)-[:BELONGS_TO_FACILITY]->(f);

// Data migration: Create relationships from existing Publication.experimentalSystem strings
// Note: keeping the string property for backward compatibility
MATCH (p:Publication) WHERE p.experimentalSystem IS NOT NULL AND p.experimentalSystem <> ""
MATCH (es:ExperimentalSystem {name: p.experimentalSystem})
MERGE (p)-[:HAS_EXPERIMENTAL_SYSTEM]->(es);
