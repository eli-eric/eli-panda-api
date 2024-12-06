package publicationsservice

import (
	"panda/apigateway/helpers"
	"panda/apigateway/services/publications-service/models"

	"github.com/neo4j/neo4j-go-driver/v4/neo4j"
)

type PublicationsService struct {
	neo4jDriver *neo4j.Driver
}

type IPublicationsService interface {
	GetPublicationByUid(uid string) (models.Publication, error)
	CreatePublication(publication *models.Publication, userUID string) (result models.Publication, err error)
	UpdatePublication(newPublication *models.Publication, userUID string) (result models.Publication, err error)
	DeletePublication(uid string, userUID string) (err error)
	GetPublications(searchText string, page, pageSize int) (result []models.Publication, totalCount int64, err error)
}

func NewPublicationsService(driver *neo4j.Driver) IPublicationsService {
	return &PublicationsService{neo4jDriver: driver}
}

func (svc *PublicationsService) GetPublicationByUid(uid string) (result models.Publication, err error) {

	session, _ := helpers.NewNeo4jSession(*svc.neo4jDriver)

	result.Uid = uid

	err = helpers.GetSingleNode(session, &result)

	return result, err
}

func (svc *PublicationsService) CreatePublication(publication *models.Publication, userUID string) (result models.Publication, err error) {

	session, _ := helpers.NewNeo4jSession(*svc.neo4jDriver)

	updateQuery := helpers.DatabaseQuery{}
	updateQuery.Parameters = make(map[string]interface{})
	updateQuery.Query = `MERGE (n:Publication {uid: $uid}) `
	updateQuery.Parameters["uid"] = publication.Uid

	helpers.AutoResolveObjectToUpdateQuery(&updateQuery, *publication, models.Publication{}, "n")

	updateQuery.Query += ` RETURN n.uid as uid `
	updateQuery.ReturnAlias = "uid"

	historyLog := helpers.HistoryLogQuery(publication.Uid, "CREATE", userUID)

	err = helpers.WriteNeo4jAndReturnNothingMultipleQueries(session,
		updateQuery,
		historyLog)

	return *publication, err
}

func (svc *PublicationsService) UpdatePublication(publication *models.Publication, userUID string) (result models.Publication, err error) {

	session, _ := helpers.NewNeo4jSession(*svc.neo4jDriver)

	oldPublication, err := svc.GetPublicationByUid(publication.Uid)

	if err != nil {
		return result, err
	}

	updateQuery := helpers.DatabaseQuery{}
	updateQuery.Parameters = make(map[string]interface{})
	updateQuery.Query = `MATCH (n:Publication {uid: $uid}) `
	updateQuery.Parameters["uid"] = publication.Uid

	helpers.AutoResolveObjectToUpdateQuery(&updateQuery, *publication, oldPublication, "n")

	updateQuery.Query += ` RETURN n.uid as uid `
	updateQuery.ReturnAlias = "uid"

	historyLog := helpers.HistoryLogQuery(publication.Uid, "UPDATE", userUID)

	err = helpers.WriteNeo4jAndReturnNothingMultipleQueries(session,
		updateQuery,
		historyLog)

	return result, err
}

func (svc *PublicationsService) DeletePublication(uid string, userUID string) (err error) {

	session, _ := helpers.NewNeo4jSession(*svc.neo4jDriver)

	err = helpers.WriteNeo4jAndReturnNothingMultipleQueries(session,
		helpers.SoftDeleteNodeQuery(uid),
		helpers.HistoryLogQuery(uid, "DELETE", userUID))

	return err

}

func (svc *PublicationsService) GetPublications(searchText string, page, pageSize int) (result []models.Publication, totalCount int64, err error) {

	session, _ := helpers.NewNeo4jSession(*svc.neo4jDriver)

	skip := pageSize
	limit := pageSize

	if skip > 0 {
		skip = (page - 1) * pageSize
	} else {
		skip = 0
	}

	if limit <= 0 {
		limit = 10
	}

	result, totalCount, err = helpers.GetMultipleNodes[models.Publication](session, skip, limit, searchText)

	helpers.ProcessArrayResult(&result, err)

	return result, totalCount, err
}
