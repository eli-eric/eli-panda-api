CREATE CONSTRAINT SystemType_uid_unique IF NOT EXISTS FOR (r:SystemType) REQUIRE r.uid IS UNIQUE;
