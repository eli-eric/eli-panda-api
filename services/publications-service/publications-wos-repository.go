package publicationsservice

import (
	"errors"
	"fmt"
	"strings"

	"panda/apigateway/helpers"
	codebookmodels "panda/apigateway/services/codebook-service/models"
	"panda/apigateway/services/publications-service/models"

	"github.com/neo4j/neo4j-go-driver/v4/neo4j"
)

type neo4jWosImportRepository struct {
	driver *neo4j.Driver
}

func (repo *neo4jWosImportRepository) findExistingPublication(
	doi,
	currentPublicationUID string,
) (*models.WosExistingPublication, error) {
	session, err := helpers.NewNeo4jSession(*repo.driver)
	if err != nil {
		return nil, err
	}
	defer session.Close()

	query := helpers.DatabaseQuery{
		Query: `
			MATCH (publication:Publication)
			WHERE (publication.deleted IS NULL OR publication.deleted = false)
			  AND ($currentPublicationUID = "" OR publication.uid <> $currentPublicationUID)
			WITH publication, toLower(trim(publication.doi)) AS storedDOI
			WITH publication,
				CASE
					WHEN storedDOI STARTS WITH "https://dx.doi.org/" THEN substring(storedDOI, size("https://dx.doi.org/"))
					WHEN storedDOI STARTS WITH "http://dx.doi.org/" THEN substring(storedDOI, size("http://dx.doi.org/"))
					WHEN storedDOI STARTS WITH "https://doi.org/" THEN substring(storedDOI, size("https://doi.org/"))
					WHEN storedDOI STARTS WITH "http://doi.org/" THEN substring(storedDOI, size("http://doi.org/"))
					WHEN storedDOI STARTS WITH "doi:" THEN trim(substring(storedDOI, size("doi:")))
					ELSE storedDOI
				END AS normalizedStoredDOI
			WHERE normalizedStoredDOI = $doi
			WITH publication
			ORDER BY publication.updatedAt DESC
			RETURN {
				uid: publication.uid,
				code: coalesce(publication.code, ""),
				title: coalesce(publication.title, ""),
				doi: publication.doi
			} AS publication
			LIMIT 1`,
		ReturnAlias: "publication",
		Parameters: map[string]interface{}{
			"doi":                   doi,
			"currentPublicationUID": currentPublicationUID,
		},
	}

	publication, err := helpers.GetNeo4jSingleRecordAndMapToStruct[models.WosExistingPublication](session, query)
	if errors.Is(err, helpers.ERR_NO_ROWS) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &publication, nil
}

func (repo *neo4jWosImportRepository) listResearchersForWos() ([]wosResearcherRecord, error) {
	session, err := helpers.NewNeo4jSession(*repo.driver)
	if err != nil {
		return nil, err
	}
	defer session.Close()

	query := helpers.DatabaseQuery{
		Query: `
			MATCH (researcher:Researcher)
			WHERE researcher.deleted IS NULL OR researcher.deleted = false
			RETURN {
				uid: researcher.uid,
				firstName: coalesce(researcher.firstName, ""),
				lastName: coalesce(researcher.lastName, ""),
				currentResearcherId: coalesce(toUpper(trim(researcher.researcherId)), ""),
				researcherIds: [id IN
					coalesce(researcher.researcherIds, []) +
					CASE WHEN researcher.researcherId IS NULL THEN [] ELSE [researcher.researcherId] END
					WHERE id IS NOT NULL AND trim(id) <> "" |
					toUpper(trim(id))]
			} AS researcher`,
		ReturnAlias: "researcher",
		Parameters:  map[string]interface{}{},
	}

	researchers, err := helpers.GetNeo4jArrayOfNodes[wosResearcherRecord](session, query)
	if researchers == nil {
		researchers = make([]wosResearcherRecord, 0)
	}
	return researchers, err
}

func (repo *neo4jWosImportRepository) findMediaType(
	code,
	facilityCode string,
) (*codebookmodels.Codebook, error) {
	session, err := helpers.NewNeo4jSession(*repo.driver)
	if err != nil {
		return nil, err
	}
	defer session.Close()

	query := helpers.DatabaseQuery{
		Query: `
			MATCH (mediaType:MediaType {code: $code})-[:BELONGS_TO_FACILITY]->(facility:Facility {code: $facilityCode})
			WHERE mediaType.deleted IS NULL OR mediaType.deleted = false
			RETURN {
				uid: mediaType.uid,
				name: mediaType.name,
				code: mediaType.code
			} AS mediaType
			LIMIT 1`,
		ReturnAlias: "mediaType",
		Parameters: map[string]interface{}{
			"code":         code,
			"facilityCode": facilityCode,
		},
	}

	mediaType, err := helpers.GetNeo4jSingleRecordAndMapToStruct[codebookmodels.Codebook](session, query)
	if errors.Is(err, helpers.ERR_NO_ROWS) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &mediaType, nil
}

// rememberResearcherIDResult carries the researcher's IDs after the write plus
// whichever one is now current — the single value RIV export sends onward.
type rememberResearcherIDResult struct {
	researcherIDs []string
	primaryID     string
}

func (repo *neo4jWosImportRepository) rememberResearcherID(
	researcherUID,
	researcherID,
	userUID string,
	makePrimary bool,
) (rememberResearcherIDResult, error) {
	session, err := helpers.NewNeo4jSession(*repo.driver)
	if err != nil {
		return rememberResearcherIDResult{}, err
	}
	defer session.Close()

	result, err := session.WriteTransaction(func(tx neo4j.Transaction) (interface{}, error) {
		lookupResult, err := tx.Run(`
			MATCH (target:Researcher {uid: $researcherUID})
			WHERE target.deleted IS NULL OR target.deleted = false
			OPTIONAL MATCH (other:Researcher)
			WHERE other.uid <> target.uid
			  AND (other.deleted IS NULL OR other.deleted = false)
			  AND any(id IN
				coalesce(other.researcherIds, []) +
				CASE WHEN other.researcherId IS NULL THEN [] ELSE [other.researcherId] END
				WHERE toUpper(trim(id)) = $researcherID)
			RETURN
				[id IN coalesce(target.researcherIds, []) +
					CASE WHEN target.researcherId IS NULL THEN [] ELSE [target.researcherId] END
					WHERE id IS NOT NULL AND trim(id) <> "" | toUpper(trim(id))] AS currentIDs,
				coalesce(toUpper(trim(target.researcherId)), "") AS currentPrimary,
				[uid IN collect(other.uid) WHERE uid IS NOT NULL] AS conflictingUIDs`, map[string]interface{}{
			"researcherUID": researcherUID,
			"researcherID":  researcherID,
		})
		if err != nil {
			return nil, err
		}
		record, err := lookupResult.Single()
		if err != nil {
			if neo4j.IsUsageError(err) && strings.Contains(err.Error(), "no more records") {
				return nil, errResearcherMissing
			}
			return nil, err
		}

		conflictingValue, _ := record.Get("conflictingUIDs")
		if len(stringSlice(conflictingValue)) > 0 {
			return nil, errResearcherIDOwned
		}
		currentValue, _ := record.Get("currentIDs")
		researcherIDs := uniqueNormalizedResearcherIDs(stringSlice(currentValue), researcherID)

		primaryValue, _ := record.Get("currentPrimary")
		primaryID, _ := primaryValue.(string)
		primaryID = resolvePrimaryResearcherID(primaryID, researcherID, researcherIDs, makePrimary)

		updateResult, err := tx.Run(`
			MATCH (target:Researcher {uid: $researcherUID})
			WHERE target.deleted IS NULL OR target.deleted = false
			MATCH (user:User {uid: $userUID})
			SET target.researcherIds = $researcherIDs,
				target.researcherId = $primaryID,
				target.updatedAt = datetime()
			CREATE (target)-[:WAS_UPDATED_BY {at: datetime(), action: "UPDATE"}]->(user)
			RETURN target.researcherIds AS researcherIDs,
				coalesce(target.researcherId, "") AS primaryID`, map[string]interface{}{
			"researcherUID": researcherUID,
			"researcherIDs": researcherIDs,
			"primaryID":     primaryID,
			"userUID":       userUID,
		})
		if err != nil {
			return nil, err
		}
		updatedRecord, err := updateResult.Single()
		if err != nil {
			return nil, fmt.Errorf("researcher or audit user not found: %w", err)
		}
		updatedValue, _ := updatedRecord.Get("researcherIDs")
		updatedPrimary, _ := updatedRecord.Get("primaryID")
		primary, _ := updatedPrimary.(string)
		return rememberResearcherIDResult{
			researcherIDs: stringSlice(updatedValue),
			primaryID:     primary,
		}, nil
	})
	if err != nil {
		return rememberResearcherIDResult{}, err
	}

	remembered, ok := result.(rememberResearcherIDResult)
	if !ok {
		return rememberResearcherIDResult{}, fmt.Errorf("unexpected researcher ID result %T", result)
	}
	return remembered, nil
}

// resolvePrimaryResearcherID decides which ID becomes current after a confirmed
// match. An empty field is filled outright — nothing can be lost — preferring
// the newest known ID and falling back to the one just confirmed when the
// newest is ambiguous. Otherwise the field only moves when the caller asked for
// it, so importing an old paper can never demote a current ID behind the
// librarian's back.
func resolvePrimaryResearcherID(
	currentPrimary, confirmedID string,
	researcherIDs []string,
	makePrimary bool,
) string {
	if strings.TrimSpace(currentPrimary) == "" {
		if newest, ok := newestResearcherID(researcherIDs); ok {
			return newest
		}
		return confirmedID
	}

	if makePrimary {
		return confirmedID
	}

	return currentPrimary
}

func stringSlice(value interface{}) []string {
	switch values := value.(type) {
	case []string:
		return append([]string(nil), values...)
	case []interface{}:
		result := make([]string, 0, len(values))
		for _, value := range values {
			if stringValue, ok := value.(string); ok {
				result = append(result, stringValue)
			}
		}
		return result
	default:
		return []string{}
	}
}

func uniqueNormalizedResearcherIDs(values []string, additions ...string) []string {
	result := make([]string, 0, len(values)+len(additions))
	seen := make(map[string]struct{})
	for _, value := range append(append([]string(nil), values...), additions...) {
		normalized := strings.ToUpper(strings.TrimSpace(value))
		if normalized == "" {
			continue
		}
		if _, exists := seen[normalized]; exists {
			continue
		}
		seen[normalized] = struct{}{}
		result = append(result, normalized)
	}
	return result
}
