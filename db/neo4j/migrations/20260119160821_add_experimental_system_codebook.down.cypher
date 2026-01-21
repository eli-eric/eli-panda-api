// Remove relationships (string property already exists for backward compatibility)
MATCH (p:Publication)-[r:HAS_EXPERIMENTAL_SYSTEM]->(es:ExperimentalSystem)
DELETE r;

// Remove ExperimentalSystem nodes
MATCH (n:ExperimentalSystem) DETACH DELETE n;

// Drop constraints and indexes
DROP INDEX experimentalsystem_name_index IF EXISTS;
DROP CONSTRAINT experimentalsystem_uid_unique IF EXISTS;
