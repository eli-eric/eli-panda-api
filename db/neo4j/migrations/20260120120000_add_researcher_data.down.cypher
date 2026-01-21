// Remove all researchers added by this migration
MATCH (r:Researcher) DETACH DELETE r;
