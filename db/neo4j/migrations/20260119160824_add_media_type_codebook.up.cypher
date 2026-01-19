// Create MediaType constraint and index
CREATE CONSTRAINT mediatype_uid_unique IF NOT EXISTS FOR (n:MediaType) REQUIRE n.uid IS UNIQUE;
CREATE INDEX mediatype_name_index IF NOT EXISTS FOR (n:MediaType) ON (n.name);

// Note: Seed data will be added manually or through the API
// This codebook starts empty
