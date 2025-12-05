// Drop indexes
DROP INDEX OperationalState_name_index IF EXISTS;
DROP INDEX OperationalState_code_index IF EXISTS;

// Drop constraints
DROP CONSTRAINT OperationalState_uid_unique IF EXISTS;
DROP CONSTRAINT OperationalState_code_unique IF EXISTS;
DROP CONSTRAINT OperationalState_name_unique IF EXISTS;