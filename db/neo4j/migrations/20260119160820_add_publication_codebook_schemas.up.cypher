// Schema for ExperimentalSystem
CREATE CONSTRAINT experimentalsystem_uid_unique IF NOT EXISTS FOR (n:ExperimentalSystem) REQUIRE n.uid IS UNIQUE;
CREATE INDEX experimentalsystem_name_index IF NOT EXISTS FOR (n:ExperimentalSystem) ON (n.name);

// Schema for Grant
CREATE CONSTRAINT grant_uid_unique IF NOT EXISTS FOR (n:Grant) REQUIRE n.uid IS UNIQUE;
CREATE INDEX grant_name_index IF NOT EXISTS FOR (n:Grant) ON (n.name);

// Schema for UserExperiment (if not exist)
CREATE CONSTRAINT userexperiment_uid_unique IF NOT EXISTS FOR (n:UserExperiment) REQUIRE n.uid IS UNIQUE;
CREATE INDEX userexperiment_name_index IF NOT EXISTS FOR (n:UserExperiment) ON (n.name);

// Schema for MediaType
CREATE CONSTRAINT mediatype_uid_unique IF NOT EXISTS FOR (n:MediaType) REQUIRE n.uid IS UNIQUE;
CREATE INDEX mediatype_name_index IF NOT EXISTS FOR (n:MediaType) ON (n.name);
