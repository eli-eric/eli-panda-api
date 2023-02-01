CREATE CONSTRAINT facility_uid_unique IF NOT EXISTS
FOR (f:Facility) REQUIRE f.uid IS UNIQUE;

CREATE CONSTRAINT facility_name_unique IF NOT EXISTS
FOR (f:Facility) REQUIRE f.name IS UNIQUE;

CREATE CONSTRAINT facility_code_unique IF NOT EXISTS
FOR (f:Facility) REQUIRE f.code IS UNIQUE;
