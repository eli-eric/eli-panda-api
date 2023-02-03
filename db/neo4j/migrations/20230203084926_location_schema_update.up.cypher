CREATE CONSTRAINT Location_uid_unique IF NOT EXISTS FOR (r:Location) REQUIRE r.uid IS UNIQUE;
