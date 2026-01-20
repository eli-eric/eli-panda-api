// MediaType codebook - schema created in 20260119160820
// Seed data for media types
MATCH (f:Facility {code: "B"})
MERGE (mt:MediaType {code: "J"})
ON CREATE SET mt.uid = apoc.create.uuid(), mt.name = "J - Peer Reviewed Article"
MERGE (mt)-[:BELONGS_TO_FACILITY]->(f);

MATCH (f:Facility {code: "B"})
MERGE (mt:MediaType {code: "C"})
ON CREATE SET mt.uid = apoc.create.uuid(), mt.name = "C - Chapter in Book or Technical Paper"
MERGE (mt)-[:BELONGS_TO_FACILITY]->(f);

MATCH (f:Facility {code: "B"})
MERGE (mt:MediaType {code: "D"})
ON CREATE SET mt.uid = apoc.create.uuid(), mt.name = "D - Article in Proceedings"
MERGE (mt)-[:BELONGS_TO_FACILITY]->(f);
