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
	GetPublicationByUid(uid string) (result models.Publication, err error)
	CreatePublication(publication *models.Publication, userUID string) (result models.Publication, err error)
	UpdatePublication(newPublication *models.Publication, userUID string) (result models.Publication, err error)
	DeletePublication(uid string, userUID string) (err error)
	GetPublications() (result []models.Publication, err error)
}

func NewPublicationsService(driver *neo4j.Driver) IPublicationsService {
	return &PublicationsService{neo4jDriver: driver}
}

func (svc *PublicationsService) GetPublicationByUid(uid string) (result models.Publication, err error) {

	session, _ := helpers.NewNeo4jSession(*svc.neo4jDriver)

	query := GetPublicationByUidQuery(uid)
	result, err = helpers.GetNeo4jSingleRecordAndMapToStruct[models.Publication](session, query)

	return result, err
}

func (svc *PublicationsService) CreatePublication(publication *models.Publication, userUID string) (result models.Publication, err error) {

	session, _ := helpers.NewNeo4jSession(*svc.neo4jDriver)

	createPublication, err := helpers.CreateOrUpdateNodeQuery(session, publication)

	if err != nil {
		return result, err
	}

	err = helpers.WriteNeo4jAndReturnNothingMultipleQueries(session,
		createPublication,
		helpers.HistoryLogQuery(publication.Uid, "CREATE", userUID))

	return *publication, err
}

func (svc *PublicationsService) UpdatePublication(newPublication *models.Publication, userUID string) (result models.Publication, err error) {

	session, _ := helpers.NewNeo4jSession(*svc.neo4jDriver)

	updatePublication, err := helpers.CreateOrUpdateNodeQuery(session, newPublication)

	if err != nil {
		return result, err
	}

	err = helpers.WriteNeo4jAndReturnNothingMultipleQueries(session,
		updatePublication,
		helpers.HistoryLogQuery(newPublication.Uid, "UPDATE", userUID))

	return result, err
}

func (svc *PublicationsService) DeletePublication(uid string, userUID string) (err error) {

	session, _ := helpers.NewNeo4jSession(*svc.neo4jDriver)

	err = helpers.WriteNeo4jAndReturnNothingMultipleQueries(session,
		helpers.SoftDeleteNodeQuery(uid),
		helpers.HistoryLogQuery(uid, "DELETE", userUID))

	return err

}

func (svc *PublicationsService) GetPublications() (result []models.Publication, err error) {

	session, _ := helpers.NewNeo4jSession(*svc.neo4jDriver)

	query := GetAllPublicationsQuery()
	result, err = helpers.GetNeo4jArrayOfNodes[models.Publication](session, query)

	helpers.ProcessArrayResult(&result, err)

	return result, err
}
