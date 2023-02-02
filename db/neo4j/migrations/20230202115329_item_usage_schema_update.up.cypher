CREATE CONSTRAINT ItemUsage_uid_unique IF NOT EXISTS FOR (r:ItemUsage) REQUIRE r.uid IS UNIQUE;

CREATE CONSTRAINT ItemUsage_code_unique IF NOT EXISTS FOR (r:ItemUsage) REQUIRE r.code IS UNIQUE;

CREATE CONSTRAINT ItemUsage_name_unique IF NOT EXISTS FOR (r:ItemUsage) REQUIRE r.name IS UNIQUE;
