// ==================== GrantGroup seed data ====================

MATCH (f:Facility {code: "B"})
MERGE (gg:GrantGroup {code: "MSM"})
ON CREATE SET gg.uid = apoc.create.uuid(), gg.name = "Ministerstvo skolstvi, mladeze a telovychovy"
MERGE (gg)-[:BELONGS_TO_FACILITY]->(f);

MATCH (f:Facility {code: "B"})
MERGE (gg:GrantGroup {code: "GA0"})
ON CREATE SET gg.uid = apoc.create.uuid(), gg.name = "Grantova agentura Ceske republiky"
MERGE (gg)-[:BELONGS_TO_FACILITY]->(f);

MATCH (f:Facility {code: "B"})
MERGE (gg:GrantGroup {code: "TA0"})
ON CREATE SET gg.uid = apoc.create.uuid(), gg.name = "Technologicka agentura CR"
MERGE (gg)-[:BELONGS_TO_FACILITY]->(f);

MATCH (f:Facility {code: "B"})
MERGE (gg:GrantGroup {code: "OTHER"})
ON CREATE SET gg.uid = apoc.create.uuid(), gg.name = "Projekty dalsi"
MERGE (gg)-[:BELONGS_TO_FACILITY]->(f);

// ==================== MSM - Ministerstvo skolstvi, mladeze a telovychovy ====================

MATCH (f:Facility {code: "B"})
MATCH (gg:GrantGroup {code: "MSM"})
MERGE (g:Grant {code: "8K0101/0000449"})
ON CREATE SET g.uid = apoc.create.uuid(), g.name = "High Field Initiative (Vyzkum velmi intenzivnich poli) (2017 - 2023) HIFI"
MERGE (g)-[:BELONGS_TO_FACILITY]->(f)
MERGE (g)-[:BELONGS_TO_GROUP]->(gg);

MATCH (f:Facility {code: "B"})
MATCH (gg:GrantGroup {code: "MSM"})
MERGE (g:Grant {code: "EF15_008/0000162"})
ON CREATE SET g.uid = apoc.create.uuid(), g.name = "ELI - EXTREME LIGHT INFRASTRUCTURE - faze 2 (2016 - 2018)"
MERGE (g)-[:BELONGS_TO_FACILITY]->(f)
MERGE (g)-[:BELONGS_TO_GROUP]->(gg);

MATCH (f:Facility {code: "B"})
MATCH (gg:GrantGroup {code: "MSM"})
MERGE (g:Grant {code: "EF16_013/0001793"})
ON CREATE SET g.uid = apoc.create.uuid(), g.name = "Pokrocile simulacni nastroje pro ELI Beamlines (2017 - 2019)"
MERGE (g)-[:BELONGS_TO_FACILITY]->(f)
MERGE (g)-[:BELONGS_TO_GROUP]->(gg);

MATCH (f:Facility {code: "B"})
MATCH (gg:GrantGroup {code: "MSM"})
MERGE (g:Grant {code: "EF16_019/0000789"})
ON CREATE SET g.uid = apoc.create.uuid(), g.name = "Pokrocily vyzkum s vyuzitim fotonu a castic vytvorenych vysoce intenzivnimi lasery (2018 - 2023) ADONIS"
MERGE (g)-[:BELONGS_TO_FACILITY]->(f)
MERGE (g)-[:BELONGS_TO_GROUP]->(gg);

// ==================== GA0 - Grantova agentura Ceske republiky ====================

MATCH (f:Facility {code: "B"})
MATCH (gg:GrantGroup {code: "GA0"})
MERGE (g:Grant {code: "GA20-19854S"})
ON CREATE SET g.uid = apoc.create.uuid(), g.name = "Studium urychleni castic v astrofyzikalnich vytryscich (2020 - 2023)"
MERGE (g)-[:BELONGS_TO_FACILITY]->(f)
MERGE (g)-[:BELONGS_TO_GROUP]->(gg);

MATCH (f:Facility {code: "B"})
MATCH (gg:GrantGroup {code: "GA0"})
MERGE (g:Grant {code: "GA21-05180S"})
ON CREATE SET g.uid = apoc.create.uuid(), g.name = "Prenos naboje v chromofor-proteinovychlzy (2024 - 2027)"
MERGE (g)-[:BELONGS_TO_FACILITY]->(f)
MERGE (g)-[:BELONGS_TO_GROUP]->(gg);

MATCH (f:Facility {code: "B"})
MATCH (gg:GrantGroup {code: "GA0"})
MERGE (g:Grant {code: "GM21-09692M"})
ON CREATE SET g.uid = apoc.create.uuid(), g.name = "Stanoveni kvantovych limitu v biomolekulach pomoci entanglovanych fotonu generovanych z navazaneho kofaktoru, modelovano na OCP proteinu. (2021 - 2025)"
MERGE (g)-[:BELONGS_TO_FACILITY]->(f)
MERGE (g)-[:BELONGS_TO_GROUP]->(gg);

MATCH (f:Facility {code: "B"})
MATCH (gg:GrantGroup {code: "GA0"})
MERGE (g:Grant {code: "GA21-19779S"})
ON CREATE SET g.uid = apoc.create.uuid(), g.name = "GA21-19779S"
MERGE (g)-[:BELONGS_TO_FACILITY]->(f)
MERGE (g)-[:BELONGS_TO_GROUP]->(gg);

MATCH (f:Facility {code: "B"})
MATCH (gg:GrantGroup {code: "GA0"})
MERGE (g:Grant {code: "GA22-20012S"})
ON CREATE SET g.uid = apoc.create.uuid(), g.name = "GA22-20012S"
MERGE (g)-[:BELONGS_TO_FACILITY]->(f)
MERGE (g)-[:BELONGS_TO_GROUP]->(gg);

MATCH (f:Facility {code: "B"})
MATCH (gg:GrantGroup {code: "GA0"})
MERGE (g:Grant {code: "GF22-06059L"})
ON CREATE SET g.uid = apoc.create.uuid(), g.name = "NSF-GACR collaborative Grant No. 2206059"
MERGE (g)-[:BELONGS_TO_FACILITY]->(f)
MERGE (g)-[:BELONGS_TO_GROUP]->(gg);

// ==================== TA0 - Technologicka agentura CR ====================

MATCH (f:Facility {code: "B"})
MATCH (gg:GrantGroup {code: "TA0"})
MERGE (g:Grant {code: "TQ11000025"})
ON CREATE SET g.uid = apoc.create.uuid(), g.name = "Optimalizace PoC procesu v ELI Beamlines (2024 - 2028)"
MERGE (g)-[:BELONGS_TO_FACILITY]->(f)
MERGE (g)-[:BELONGS_TO_GROUP]->(gg);

// ==================== OTHER - Projekty dalsi ====================

MATCH (f:Facility {code: "B"})
MATCH (gg:GrantGroup {code: "OTHER"})
MERGE (g:Grant {code: "Horizon"})
ON CREATE SET g.uid = apoc.create.uuid(), g.name = "Horizon"
MERGE (g)-[:BELONGS_TO_FACILITY]->(f)
MERGE (g)-[:BELONGS_TO_GROUP]->(gg);

MATCH (f:Facility {code: "B"})
MATCH (gg:GrantGroup {code: "OTHER"})
MERGE (g:Grant {code: "COST"})
ON CREATE SET g.uid = apoc.create.uuid(), g.name = "COST"
MERGE (g)-[:BELONGS_TO_FACILITY]->(f)
MERGE (g)-[:BELONGS_TO_GROUP]->(gg);

MATCH (f:Facility {code: "B"})
MATCH (gg:GrantGroup {code: "OTHER"})
MERGE (g:Grant {code: "Impulse"})
ON CREATE SET g.uid = apoc.create.uuid(), g.name = "Impulse"
MERGE (g)-[:BELONGS_TO_FACILITY]->(f)
MERGE (g)-[:BELONGS_TO_GROUP]->(gg);
