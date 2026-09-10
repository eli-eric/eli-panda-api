CREATE CONSTRAINT Team_uid_unique IF NOT EXISTS
FOR (t:Team)
REQUIRE t.uid IS UNIQUE;
CREATE INDEX Team_name_index IF NOT EXISTS
FOR (t:Team)
ON (t.name);