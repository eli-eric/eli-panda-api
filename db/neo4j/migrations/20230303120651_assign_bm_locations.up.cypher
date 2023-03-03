MATCH (n:Location) where n.facility = "B"
match(f:Facility{code:"B"})
merge(n)-[:BELONGS_TO_FACILITY]->(f);