// The up migration creates only schema objects. Keep reporting data intact:
// Neo4j 4.4 cannot mix data writes and schema changes in this transaction.
DROP INDEX JournalMetric_year_index IF EXISTS;
DROP INDEX PublicationReporting_reviewed_index IF EXISTS;
DROP INDEX PublicationReporting_classification_index IF EXISTS;
DROP CONSTRAINT JournalMetric_uid_unique IF EXISTS;
DROP CONSTRAINT PublicationReportingAuthor_uid_unique IF EXISTS;
DROP CONSTRAINT PublicationReporting_uid_unique IF EXISTS;
