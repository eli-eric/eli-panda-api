package publicationsservice

import (
	"errors"

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
				researcherIds: [id IN
					CASE WHEN researcher.researcherId IS NULL THEN [] ELSE [researcher.researcherId] END
					WHERE trim(id) <> "" | toUpper(trim(id))]
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
