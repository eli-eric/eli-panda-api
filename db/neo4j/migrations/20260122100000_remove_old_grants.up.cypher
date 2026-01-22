// Remove existing Publication-Grant relationships
MATCH (p:Publication)-[r:HAS_GRANT]->(g:Grant)
DELETE r;

// Remove all existing Grant nodes
MATCH (n:Grant) DETACH DELETE n;
