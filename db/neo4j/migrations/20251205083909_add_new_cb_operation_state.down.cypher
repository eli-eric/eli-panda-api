// Remove OperationalState nodes
MATCH (r:OperationalState)
DETACH DELETE r;

