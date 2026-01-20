// Drop schemas in reverse order
DROP INDEX mediatype_name_index IF EXISTS;
DROP CONSTRAINT mediatype_uid_unique IF EXISTS;
DROP INDEX userexperiment_name_index IF EXISTS;
DROP CONSTRAINT userexperiment_uid_unique IF EXISTS;
DROP INDEX grant_name_index IF EXISTS;
DROP CONSTRAINT grant_uid_unique IF EXISTS;
DROP INDEX experimentalsystem_name_index IF EXISTS;
DROP CONSTRAINT experimentalsystem_uid_unique IF EXISTS;
