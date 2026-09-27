MATCH (r:Researcher)
WHERE r.researcherIds IS NOT NULL
REMOVE r.researcherIds
