CREATE CONSTRAINT SystemCriticality_uid_unique IF NOT EXISTS FOR (r:SystemCriticality) REQUIRE r.uid IS UNIQUE;
CREATE CONSTRAINT SystemCriticality_code_unique IF NOT EXISTS FOR (r:SystemCriticality) REQUIRE r.code IS UNIQUE;
CREATE CONSTRAINT SystemCriticality_name_unique IF NOT EXISTS FOR (r:SystemCriticality) REQUIRE r.name IS UNIQUE;