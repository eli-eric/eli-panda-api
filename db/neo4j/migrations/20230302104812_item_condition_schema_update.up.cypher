CREATE CONSTRAINT ItemCondition_uid_unique IF NOT EXISTS FOR (r:ItemCondition) REQUIRE r.uid IS UNIQUE;
CREATE CONSTRAINT ItemCondition_code_unique IF NOT EXISTS FOR (r:ItemCondition) REQUIRE r.code IS UNIQUE;
CREATE CONSTRAINT ItemCondition_name_unique IF NOT EXISTS FOR (r:ItemCondition) REQUIRE r.name IS UNIQUE;