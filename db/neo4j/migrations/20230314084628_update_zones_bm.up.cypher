MATCH (n:Zone)<-[r:HAS_ZONE]-(f:Facility) detach delete r;
MATCH (z:Zone) MATCH(f:Facility{code:"B"}) WITH z,f MERGE (z)-[:BELONGS_TO_FACILITY]->(f);