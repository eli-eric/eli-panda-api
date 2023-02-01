MATCH (r:SystemType)
DETACH DELETE (r);

MATCH (r:SystemTypeGroup)
DETACH DELETE (r);
