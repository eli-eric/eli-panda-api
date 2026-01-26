// Remove Grant nodes
MATCH (g:Grant) DETACH DELETE g;

// Remove GrantGroup nodes
MATCH (gg:GrantGroup) DETACH DELETE gg;
