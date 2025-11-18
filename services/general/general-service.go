package general

import (
	"panda/apigateway/helpers"
	"panda/apigateway/services/general/models"

	"github.com/neo4j/neo4j-go-driver/v4/neo4j"
	"github.com/rs/zerolog/log"
)

type GeneralService struct {
	neo4jDriver *neo4j.Driver
}

type IGeneralService interface {
	GetGraphNodesByUid(uid string) (result []models.GraphNode, err error)
	GetGraphLinksByUid(uid string) (result []models.GraphLink, err error)
	GetUUID() (uuid string, err error)
	GlobalSearch(searchText string, page int, pageSize int) (result models.GlobalSearchResponse, err error)
}

func NewGeneralService(driver *neo4j.Driver) IGeneralService {
	return &GeneralService{neo4jDriver: driver}
}

func (svc *GeneralService) GetGraphNodesByUid(uid string) (result []models.GraphNode, err error) {

	session, _ := helpers.NewNeo4jSession(*svc.neo4jDriver)

	query := GetGraphNodesByUidQuery(uid)
	result, err = helpers.GetNeo4jArrayOfNodes[models.GraphNode](session, query)

	helpers.ProcessArrayResult(&result, err)

	return result, err
}

func (svc *GeneralService) GetGraphLinksByUid(uid string) (result []models.GraphLink, err error) {

	session, _ := helpers.NewNeo4jSession(*svc.neo4jDriver)

	query := GetGraphLinksByUidQuery(uid)
	result, err = helpers.GetNeo4jArrayOfNodes[models.GraphLink](session, query)

	helpers.ProcessArrayResult(&result, err)

	return result, err
}

func (svc *GeneralService) GetUUID() (uuid string, err error) {

	session, _ := helpers.NewNeo4jSession(*svc.neo4jDriver)

	query := GetUUIDQuery()
	uuid, err = helpers.GetNeo4jSingleRecordSingleValue[string](session, query)

	return uuid, err
}

func (svc *GeneralService) GlobalSearch(searchText string, page int, pageSize int) (result models.GlobalSearchResponse, err error) {

	session, _ := helpers.NewNeo4jSession(*svc.neo4jDriver)

	// Calculate skip value for pagination
	skip := (page - 1) * pageSize

	log.Debug().
		Str("searchText", searchText).
		Int("page", page).
		Int("pageSize", pageSize).
		Int("skip", skip).
		Msg("GlobalSearch: preparing query")

	// Get search results
	searchQuery := GetGlobalSearchQuery(searchText, skip, pageSize)

	log.Debug().
		Interface("parameters", searchQuery.Parameters).
		Msg("GlobalSearch: executing Neo4j query")

	searchResults, err := helpers.GetNeo4jArrayOfNodes[models.GlobalSearchResult](session, searchQuery)
	if err != nil {
		log.Error().Err(err).Msg("GlobalSearch: failed to execute query")
		return result, err
	}

	log.Debug().
		Int("resultsCount", len(searchResults)).
		Msg("GlobalSearch: query returned results")

	// Get total count - this now returns multiple records that need to be summed
	countQuery := GetGlobalSearchCountQuery(searchText)
	countResults, err := helpers.GetNeo4jSingleRecordSingleValue[int64](session, countQuery)
	if err != nil {
		log.Error().Err(err).Msg("GlobalSearch: failed to get count")
		return result, err
	}

	log.Debug().
		Int64("totalCount", countResults).
		Msg("GlobalSearch: count query returned")

	result.Data = searchResults
	result.TotalCount = countResults

	helpers.ProcessArrayResult(&result.Data, err)

	return result, err
}
