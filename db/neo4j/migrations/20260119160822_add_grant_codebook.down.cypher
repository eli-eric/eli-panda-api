// Remove relationships (string property already exists for backward compatibility)
MATCH (p:Publication)-[r:HAS_GRANT]->(g:Grant)
DELETE r;

// Remove Grant nodes
MATCH (n:Grant) DETACH DELETE n;

// Drop constraints and indexes
DROP INDEX grant_name_index IF EXISTS;
DROP CONSTRAINT grant_uid_unique IF EXISTS;
