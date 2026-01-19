// Create Grant constraint and index
CREATE CONSTRAINT grant_uid_unique IF NOT EXISTS FOR (n:Grant) REQUIRE n.uid IS UNIQUE;
CREATE INDEX grant_name_index IF NOT EXISTS FOR (n:Grant) ON (n.name);

// Seed data for ELI-Beamlines facility (10 entries)
// MSM - Ministerstvo skolstvi, mladeze a telovychovy
MATCH (f:Facility {code: "B"})
MERGE (g:Grant {code: "MSM-8K0101NFRASTRUCTURE"})
ON CREATE SET g.uid = apoc.create.uuid(), g.name = "8K0101NFRASTRUCTURE - faze 2 (2016 - 2018)"
MERGE (g)-[:BELONGS_TO_FACILITY]->(f);

MATCH (f:Facility {code: "B"})
MERGE (g:Grant {code: "MSM-EF16_013/0001793"})
ON CREATE SET g.uid = apoc.create.uuid(), g.name = "EF16_013/0001793 - Pokrocile simulacni nastroje pro ELI Beamlines (2017 - 2019)"
MERGE (g)-[:BELONGS_TO_FACILITY]->(f);

MATCH (f:Facility {code: "B"})
MERGE (g:Grant {code: "MSM-EF16_019/0000789"})
ON CREATE SET g.uid = apoc.create.uuid(), g.name = "EF16_019/0000789 - Pokrocily vyzkum s vyuzitim fotonu a castic vytvorenych vysoce intenzivnimi lasery (2018 - 2023) ADONIS"
MERGE (g)-[:BELONGS_TO_FACILITY]->(f);

// GA0 - Grantova agentura Ceske republiky
MATCH (f:Facility {code: "B"})
MERGE (g:Grant {code: "GA0-GA20-19854S"})
ON CREATE SET g.uid = apoc.create.uuid(), g.name = "GA20-19854S - Studium urychleni castic v astrofyzikalnich vytryscich (2020 - 2023)"
MERGE (g)-[:BELONGS_TO_FACILITY]->(f);

MATCH (f:Facility {code: "B"})
MERGE (g:Grant {code: "GA0-GA21-05180S"})
ON CREATE SET g.uid = apoc.create.uuid(), g.name = "GA21-05180S - Prenos naboje v chromofor-proteinovych komplexech tryptofanovymi drahami (2021 - 2024)"
MERGE (g)-[:BELONGS_TO_FACILITY]->(f);

MATCH (f:Facility {code: "B"})
MERGE (g:Grant {code: "GA0-GA24-10671S"})
ON CREATE SET g.uid = apoc.create.uuid(), g.name = "GA24-10671S - Biosenzor pro prime sledovani dynamiky bunecne aktivity (2024 - 2026)"
MERGE (g)-[:BELONGS_TO_FACILITY]->(f);

MATCH (f:Facility {code: "B"})
MERGE (g:Grant {code: "GA0-GA24-11398S"})
ON CREATE SET g.uid = apoc.create.uuid(), g.name = "GA24-11398S - Plazmatem podporena synteza hybridnich nanomaterialu pro laserem rizenou proton-borovou jadernou fuzi (2024 - 2026)"
MERGE (g)-[:BELONGS_TO_FACILITY]->(f);

MATCH (f:Facility {code: "B"})
MERGE (g:Grant {code: "GA0-GF22-42890L"})
ON CREATE SET g.uid = apoc.create.uuid(), g.name = "GF22-42890L - Studium generace zareni gama doprovazejiciho interakci plazmatu s laserem o vysoke intenzite na ELI Beamlines (2022 - 2025)"
MERGE (g)-[:BELONGS_TO_FACILITY]->(f);

MATCH (f:Facility {code: "B"})
MERGE (g:Grant {code: "GA0-GF22-42963L"})
ON CREATE SET g.uid = apoc.create.uuid(), g.name = "GF22-42963L - Fyzika plazmatu silnych QED poli v laboratorich vybavenych lasery s vykonem PW a vice (2022 - 2025)"
MERGE (g)-[:BELONGS_TO_FACILITY]->(f);

MATCH (f:Facility {code: "B"})
MERGE (g:Grant {code: "GA0-GF24-14395L"})
ON CREATE SET g.uid = apoc.create.uuid(), g.name = "GF24-14395L - Koherentni zareni elektronu pri interakci s intenzivnimi laserovymi pulzy (2024 - 2027)"
MERGE (g)-[:BELONGS_TO_FACILITY]->(f);

// Data migration: Create relationships from existing Publication.grant strings
// Note: keeping the string property for backward compatibility
MATCH (p:Publication) WHERE p.grant IS NOT NULL AND p.grant <> ""
MATCH (g:Grant {name: p.grant})
MERGE (p)-[:HAS_GRANT]->(g);
