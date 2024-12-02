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
	CreatePublication(publication *models.Publication) (result models.Publication, err error)
	UpdatePublication(newPublication *models.Publication, userUID string) (result models.Publication, err error)
	DeletePublication(uid string) (err error)
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

func (svc *PublicationsService) CreatePublication(publication *models.Publication) (result models.Publication, err error) {

	session, _ := helpers.NewNeo4jSession(*svc.neo4jDriver)

	query := CreatePublicationQuery(publication)
	_, err = helpers.WriteNeo4jAndReturnSingleValue[string](session, query)

	return *publication, err
}

func (svc *PublicationsService) UpdatePublication(newPublication *models.Publication, userUID string) (result models.Publication, err error) {

	session, _ := helpers.NewNeo4jSession(*svc.neo4jDriver)

	oldOne, err := svc.GetPublicationByUid(newPublication.Uid)

	if err != nil {
		return result, err
	}

	query := UpdatePublicationQuery(newPublication, &oldOne, userUID)

	_, err = helpers.WriteNeo4jAndReturnSingleValue[string](session, query)

	return result, err
}

func (svc *PublicationsService) DeletePublication(uid string) (err error) {

	session, _ := helpers.NewNeo4jSession(*svc.neo4jDriver)

	query := DeletePublicationQuery(uid)
	err = helpers.WriteNeo4jAndReturnNothing(session, query)

	return err

}

func (svc *PublicationsService) GetPublications() (result []models.Publication, err error) {

	session, _ := helpers.NewNeo4jSession(*svc.neo4jDriver)

	query := GetAllPublicationsQuery()
	result, err = helpers.GetNeo4jArrayOfNodes[models.Publication](session, query)

	helpers.ProcessArrayResult(&result, err)

	return result, err
}
