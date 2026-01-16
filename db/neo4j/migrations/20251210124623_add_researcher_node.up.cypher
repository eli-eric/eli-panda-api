CREATE CONSTRAINT researcher_uid_unique IF NOT EXISTS
FOR (r:Researcher) REQUIRE r.uid IS UNIQUE;
