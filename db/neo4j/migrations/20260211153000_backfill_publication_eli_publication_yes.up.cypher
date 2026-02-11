// Backfill eliPublication for active publications
MATCH (p:Publication)
WHERE (p.deleted IS NULL OR p.deleted = false)
SET p.eliPublication = "YES";
