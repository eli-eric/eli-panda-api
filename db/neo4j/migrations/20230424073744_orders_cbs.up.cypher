CREATE CONSTRAINT OrderStatus_uid_unique IF NOT EXISTS FOR (r:OrderStatus) REQUIRE r.uid IS UNIQUE;
CREATE CONSTRAINT OrderStatus_code_unique IF NOT EXISTS FOR (r:OrderStatus) REQUIRE r.code IS UNIQUE;
CREATE CONSTRAINT OrderStatus_name_unique IF NOT EXISTS FOR (r:OrderStatus) REQUIRE r.name IS UNIQUE;

CREATE CONSTRAINT Supplier_uid_unique IF NOT EXISTS FOR (r:Supplier) REQUIRE r.uid IS UNIQUE;
CREATE CONSTRAINT Supplier_name_unique IF NOT EXISTS FOR (r:Supplier) REQUIRE r.name IS UNIQUE;
