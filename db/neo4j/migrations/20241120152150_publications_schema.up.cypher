CREATE CONSTRAINT publication_uid_unique IF NOT EXISTS
FOR (f:Publication) REQUIRE f.uid IS UNIQUE;