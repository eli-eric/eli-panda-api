// Rollback eliPublication backfill
MATCH (p:Publication)
WHERE (p.deleted IS NULL OR p.deleted = false)
REMOVE p.eliPublication;
