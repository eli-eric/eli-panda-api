// ELIPANDA-501: most scientists own several WoS ResearcherIDs; the single
// researcherId property silently misses authors. Carry every existing primary
// id into the new researcherIds array so matching has one list to consult.
MATCH (r:Researcher)
WHERE r.researcherId IS NOT NULL AND trim(r.researcherId) <> '' AND r.researcherIds IS NULL
SET r.researcherIds = [trim(r.researcherId)]
