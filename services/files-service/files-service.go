package filesservice

import (
	"panda/apigateway/helpers"
	"panda/apigateway/services/files-service/models"

	"github.com/neo4j/neo4j-go-driver/v4/neo4j"
)

type FilesService struct {
	neo4jDriver *neo4j.Driver
}

type IFilesService interface {
	GetFileLinksByParentUid(parentUid string) (result []models.FileLink, err error)
	CreateFileLink(parentUid string, fileLink models.FileLink) (err error, result models.FileLink)
	UpdateFileLink(fileLink models.FileLink) (err error)
	DeleteFileLink(uid string) (err error)
}

func NewFilesService(driver *neo4j.Driver) IFilesService {
	return &FilesService{neo4jDriver: driver}
}

func (svc *FilesService) GetFileLinksByParentUid(parentUid string) (result []models.FileLink, err error) {

	session, _ := helpers.NewNeo4jSession(*svc.neo4jDriver)

	query := GetFileLinksByParentUidQuery(parentUid)
	result, err = helpers.GetNeo4jArrayOfNodes[models.FileLink](session, query)

	helpers.ProcessArrayResult(&result, err)

	return result, err
}

func (svc *FilesService) CreateFileLink(parentUid string, fileLink models.FileLink) (err error, result models.FileLink) {

	session, _ := helpers.NewNeo4jSession(*svc.neo4jDriver)

	query := CreateFileLinkQuery(parentUid, fileLink)
	fileLink, err = helpers.WriteNeo4jReturnSingleRecordAndMapToStruct[models.FileLink](session, query)

	return err, fileLink
}

func (svc *FilesService) UpdateFileLink(fileLink models.FileLink) (err error) {
	return err
}

func (svc *FilesService) DeleteFileLink(uid string) (err error) {
	return err
}
