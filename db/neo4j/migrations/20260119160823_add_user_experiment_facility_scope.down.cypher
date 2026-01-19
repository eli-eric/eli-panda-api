// Remove relationships (string property already exists for backward compatibility)
MATCH (p:Publication)-[r:HAS_USER_EXPERIMENT]->(ue:UserExperiment)
DELETE r;

// Note: Not deleting UserExperiment nodes as they may have been created before this migration
// Only removing the facility relationships and seed data added by this migration
