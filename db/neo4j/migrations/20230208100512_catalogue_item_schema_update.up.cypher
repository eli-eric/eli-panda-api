CREATE CONSTRAINT CatalogueItem_uid_unique IF NOT EXISTS FOR (r:CatalogueItem) REQUIRE r.uid IS UNIQUE;
