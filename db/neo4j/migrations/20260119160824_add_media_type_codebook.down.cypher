// Remove relationships
MATCH (p:Publication)-[r:HAS_MEDIA_TYPE]->(mt:MediaType)
DELETE r;

// Remove MediaType nodes
MATCH (n:MediaType) DETACH DELETE n;

// Drop constraints and indexes
DROP INDEX mediatype_name_index IF EXISTS;
DROP CONSTRAINT mediatype_uid_unique IF EXISTS;
