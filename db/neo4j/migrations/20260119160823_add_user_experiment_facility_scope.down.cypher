// Remove relationships from publications (string property already exists for backward compatibility)
MATCH (p:Publication)-[r:HAS_USER_EXPERIMENT]->(ue:UserExperiment)
DELETE r;

// Remove seeded UserExperiment nodes (by code) - these were added by this migration
// UPM-* codes
MATCH (ue:UserExperiment) WHERE ue.code STARTS WITH "UPM-" DETACH DELETE ue;
// ELIUPM-* codes
MATCH (ue:UserExperiment) WHERE ue.code STARTS WITH "ELIUPM" DETACH DELETE ue;

// Note: Not removing constraints/indexes as UserExperiment may have existed before
// Only removing the seed data added by this migration
