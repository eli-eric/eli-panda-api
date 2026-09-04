// Researcher identity audit — READ ONLY, writes nothing.
//
// RIV export sends exactly one ResearcherID per author, so a researcher who
// exists twice, or whose current ID is stale, produces a wrong government
// delivery that nothing downstream flags.
//
// Run against a database you are allowed to read:
//   docker exec panda-dev-neo4j cypher-shell -u neo4j -p '<password>' \
//     -f /dev/stdin < db/neo4j/reports/researcher-identity-audit.cypher
//
// Report the counts before anyone builds a merge tool for this.

// 1. Coverage — how much of the register carries each identifier.
MATCH (r:Researcher) WHERE r.deleted IS NULL OR r.deleted = false
RETURN "coverage" AS report,
  count(r) AS researchers,
  count(r.researcherId) AS withCurrentID,
  count(r.researcherIds) AS withIDHistory,
  count(r.identificationNumber) AS withIdentificationNumber,
  count(r.orcid) AS withOrcid;

// 2. Same human, two records — by identification number, the strongest key.
MATCH (r:Researcher)
WHERE (r.deleted IS NULL OR r.deleted = false)
  AND r.identificationNumber IS NOT NULL AND trim(r.identificationNumber) <> ""
WITH trim(r.identificationNumber) AS identificationNumber,
     collect(r.uid) AS uids,
     collect(r.lastName + ", " + r.firstName) AS names
WHERE size(uids) > 1
RETURN "duplicate identificationNumber" AS report, identificationNumber, names, uids
ORDER BY size(uids) DESC;

// 3. Same ORCID on more than one record.
MATCH (r:Researcher)
WHERE (r.deleted IS NULL OR r.deleted = false)
  AND r.orcid IS NOT NULL AND trim(r.orcid) <> ""
WITH trim(r.orcid) AS orcid, collect(r.uid) AS uids,
     collect(r.lastName + ", " + r.firstName) AS names
WHERE size(uids) > 1
RETURN "duplicate orcid" AS report, orcid, names, uids
ORDER BY size(uids) DESC;

// 4. Same name on more than one record — weakest signal, expect false positives.
MATCH (r:Researcher) WHERE r.deleted IS NULL OR r.deleted = false
WITH toLower(trim(r.firstName)) + "|" + toLower(trim(r.lastName)) AS nameKey,
     collect(r.uid) AS uids,
     collect(coalesce(r.identificationNumber, "<none>")) AS identificationNumbers
WHERE size(uids) > 1
RETURN "duplicate name" AS report, nameKey, identificationNumbers, uids
ORDER BY size(uids) DESC;

// 5. Researchers whose exported ID is not the newest they hold. The suffix year
//    is the issue year, so a plain string comparison of the last four characters
//    ranks them.
MATCH (r:Researcher)
WHERE (r.deleted IS NULL OR r.deleted = false) AND r.researcherIds IS NOT NULL
WITH r, [id IN r.researcherIds WHERE id =~ "^[A-Z]{1,3}-[0-9]{4}-[0-9]{4}$"] AS readable
WHERE size(readable) > 1
WITH r, readable, reduce(newest = "", id IN readable |
  CASE WHEN newest = "" OR right(id, 4) > right(newest, 4) THEN id ELSE newest END) AS newest
WHERE coalesce(toUpper(trim(r.researcherId)), "") <> newest
RETURN "stale current ResearcherID" AS report,
  r.lastName + ", " + r.firstName AS researcher,
  coalesce(r.researcherId, "<none>") AS exportedNow,
  newest AS newestOnFile,
  readable AS allOnFile
ORDER BY researcher;
