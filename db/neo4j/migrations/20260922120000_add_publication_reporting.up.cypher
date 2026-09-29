CREATE CONSTRAINT PublicationReporting_uid_unique IF NOT EXISTS
FOR (r:PublicationReporting)
REQUIRE r.uid IS UNIQUE;
CREATE CONSTRAINT PublicationReportingAuthor_uid_unique IF NOT EXISTS
FOR (a:PublicationReportingAuthor)
REQUIRE a.uid IS UNIQUE;
CREATE CONSTRAINT JournalMetric_uid_unique IF NOT EXISTS
FOR (m:JournalMetric)
REQUIRE m.uid IS UNIQUE;
CREATE INDEX PublicationReporting_classification_index IF NOT EXISTS
FOR (r:PublicationReporting)
ON (r.classification);
CREATE INDEX PublicationReporting_reviewed_index IF NOT EXISTS
FOR (r:PublicationReporting)
ON (r.reviewed);
CREATE INDEX JournalMetric_year_index IF NOT EXISTS
FOR (m:JournalMetric)
ON (m.year);
