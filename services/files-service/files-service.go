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
	CreateFileLink(parentUid string, fileLink models.FileLink) (result models.FileLink, err error)
	UpdateFileLink(fileLink models.FileLink) (result models.FileLink, err error)
	DeleteFileLink(uid string) (err error)
	SetMiniImageUrlToNode(uid string, url *string, forceAll bool) (err error)
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

func (svc *FilesService) CreateFileLink(parentUid string, fileLink models.FileLink) (result models.FileLink, err error) {

	session, _ := helpers.NewNeo4jSession(*svc.neo4jDriver)

	query := CreateFileLinkQuery(parentUid, fileLink)
	result, err = helpers.WriteNeo4jReturnSingleRecordAndMapToStruct[models.FileLink](session, query)

	return result, err
}

func (svc *FilesService) UpdateFileLink(fileLink models.FileLink) (result models.FileLink, err error) {

	session, _ := helpers.NewNeo4jSession(*svc.neo4jDriver)

	query := UpdateFileLinkQuery(fileLink)
	result, err = helpers.WriteNeo4jReturnSingleRecordAndMapToStruct[models.FileLink](session, query)

	return result, err
}

func (svc *FilesService) DeleteFileLink(uid string) (err error) {

	session, _ := helpers.NewNeo4jSession(*svc.neo4jDriver)

	query := DeleteFileLinkQuery(uid)
	err = helpers.WriteNeo4jAndReturnNothing(session, query)

	return err
}

func (svc *FilesService) SetMiniImageUrlToNode(uid string, url *string, forceAll bool) (err error) {

	session, _ := helpers.NewNeo4jSession(*svc.neo4jDriver)

	query := SetMiniImageUrlToNodeQuery(uid, url, forceAll)
	err = helpers.WriteNeo4jAndReturnNothing(session, query)

	return err
}
