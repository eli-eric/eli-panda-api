DROP INDEX Team_name_index IF EXISTS;
DROP CONSTRAINT Team_uid_unique IF EXISTS;
MATCH (r:Role{code:'teams-view'}) DETACH DELETE r;
MATCH (r:Role{code:'teams-edit'}) DETACH DELETE r;
