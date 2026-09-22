MATCH (r:PublicationReporting)
OPTIONAL MATCH (r)-[:HAS_REPORTING_AUTHOR]->(a:PublicationReportingAuthor)
OPTIONAL MATCH (r)-[:HAS_JOURNAL_METRIC]->(m:JournalMetric)
DETACH DELETE r, a, m;
DROP INDEX JournalMetric_year_index IF EXISTS;
DROP INDEX PublicationReporting_reviewed_index IF EXISTS;
DROP INDEX PublicationReporting_classification_index IF EXISTS;
DROP CONSTRAINT JournalMetric_uid_unique IF EXISTS;
DROP CONSTRAINT PublicationReportingAuthor_uid_unique IF EXISTS;
DROP CONSTRAINT PublicationReporting_uid_unique IF EXISTS;
