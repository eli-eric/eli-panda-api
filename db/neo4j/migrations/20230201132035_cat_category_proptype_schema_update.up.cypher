CREATE CONSTRAINT CatalogueCategoryPropertyType_uid_unique IF NOT EXISTS FOR (r:CatalogueCategoryPropertyType) REQUIRE r.uid IS UNIQUE;

CREATE CONSTRAINT CatalogueCategoryPropertyType_code_unique IF NOT EXISTS FOR (r:CatalogueCategoryPropertyType) REQUIRE r.code IS UNIQUE;

CREATE CONSTRAINT CatalogueCategoryPropertyType_name_unique IF NOT EXISTS FOR (r:CatalogueCategoryPropertyType) REQUIRE r.name IS UNIQUE;
