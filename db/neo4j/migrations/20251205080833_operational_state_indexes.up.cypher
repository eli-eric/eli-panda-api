// Create constraints and indexes for OperationalState
CREATE CONSTRAINT OperationalState_uid_unique IF NOT EXISTS FOR (r:OperationalState) REQUIRE r.uid IS UNIQUE;
CREATE CONSTRAINT OperationalState_code_unique IF NOT EXISTS FOR (r:OperationalState) REQUIRE r.code IS UNIQUE;
CREATE CONSTRAINT OperationalState_name_unique IF NOT EXISTS FOR (r:OperationalState) REQUIRE r.name IS UNIQUE;

// Create index for better search performance
CREATE INDEX OperationalState_name_index IF NOT EXISTS FOR (r:OperationalState) ON (r.name);
CREATE INDEX OperationalState_code_index IF NOT EXISTS FOR (r:OperationalState) ON (r.code);