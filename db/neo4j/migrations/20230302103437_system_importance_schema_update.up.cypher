CREATE CONSTRAINT SystemImportance_uid_unique IF NOT EXISTS FOR (r:SystemImportance) REQUIRE r.uid IS UNIQUE;
CREATE CONSTRAINT SystemImportance_code_unique IF NOT EXISTS FOR (r:SystemImportance) REQUIRE r.code IS UNIQUE;
CREATE CONSTRAINT SystemImportance_name_unique IF NOT EXISTS FOR (r:SystemImportance) REQUIRE r.name IS UNIQUE;