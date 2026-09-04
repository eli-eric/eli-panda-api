MATCH (researcher:Researcher)
WHERE researcher.researcherId IS NOT NULL
  AND trim(researcher.researcherId) <> ""
SET researcher.researcherIds = reduce(
  researcherIds = [],
  researcherId IN coalesce(researcher.researcherIds, []) + [researcher.researcherId] |
  CASE
    WHEN any(existingId IN researcherIds WHERE toUpper(trim(existingId)) = toUpper(trim(researcherId)))
      THEN researcherIds
    ELSE researcherIds + toUpper(trim(researcherId))
  END
);
