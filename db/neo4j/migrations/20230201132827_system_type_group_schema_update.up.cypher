CREATE CONSTRAINT SystemTypeGroup_uid_unique IF NOT EXISTS FOR (r:SystemTypeGroup) REQUIRE r.uid IS UNIQUE;
