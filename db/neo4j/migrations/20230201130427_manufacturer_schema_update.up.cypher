CREATE CONSTRAINT manufacturer_name_unique IF NOT EXISTS FOR (r:Manufacturer) REQUIRE r.name IS UNIQUE;
