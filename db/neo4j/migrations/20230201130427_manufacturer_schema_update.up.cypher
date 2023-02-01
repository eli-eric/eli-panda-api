CREATE CONSTRAINT manufacturer_uid_unique IF NOT EXISTS FOR (r:Manufacturer) REQUIRE r.uid IS UNIQUE;

CREATE CONSTRAINT manufacturer_code_unique IF NOT EXISTS FOR (r:Manufacturer) REQUIRE r.code IS UNIQUE;

CREATE CONSTRAINT manufacturer_name_unique IF NOT EXISTS FOR (r:Manufacturer) REQUIRE r.name IS UNIQUE;
