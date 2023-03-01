package catalogueService

import (
	"panda/apigateway/config"
	"panda/apigateway/helpers"
	"panda/apigateway/services/catalogue-service/models"
	codebookModels "panda/apigateway/services/codebook-service/models"

	"github.com/neo4j/neo4j-go-driver/v4/neo4j"
)

type CatalogueService struct {
	neo4jDriver *neo4j.Driver
	jwtSecret   string
}

type ICatalogueService interface {
	GetCatalogueCategoriesByParentPath(parentPath string) (categories []models.CatalogueCategory, err error)
	GetCatalogueItems(search string, categoryPath string, pageSize int, page int) (result helpers.PaginationResult[models.CatalogueItem], err error)
	GetCatalogueItemWithDetailsByUid(uid string) (catalogueItem models.CatalogueItem, err error)
	GetCatalogueCategoryWithDetailsByUid(uid string) (catalogueItem models.CatalogueCategory, err error)
	GetCatalogueCategoryImageByUid(uid string) (imageBase64 string, err error)
	UpdateCatalogueCategory(catalogueCategory *models.CatalogueCategory) (err error)
	CreateCatalogueCategory(catalogueCategory *models.CatalogueCategory) (err error)
	DeleteCatalogueCategory(uid string) (err error)
	GetZonesCodebook() (result []codebookModels.Codebook, err error)
	GetSubZonesCodebook(parentUID string) (result []codebookModels.Codebook, err error)
	GetUnitsCodebook() (result []codebookModels.Codebook, err error)
	GetPropertyTypesCodebook() (result []codebookModels.Codebook, err error)
}

// Create new security service instance
func NewCatalogueService(settings *config.Config, driver *neo4j.Driver) ICatalogueService {

	return &CatalogueService{neo4jDriver: driver, jwtSecret: settings.JwtSecret}
}

func (svc *CatalogueService) GetCatalogueCategoriesByParentPath(parentPath string) (categories []models.CatalogueCategory, err error) {

	// Open a new Session
	session, _ := helpers.NewNeo4jSession(*svc.neo4jDriver)

	//get all categories by parent path
	query := CatalogueCategoriesByParentPathQuery(parentPath)
	categories, err = helpers.GetNeo4jArrayOfNodes[models.CatalogueCategory](session, query)

	helpers.ProcessArrayResult(&categories, err)

	return categories, err
}

func (svc *CatalogueService) GetCatalogueItems(search string, categoryPath string, pageSize int, page int) (result helpers.PaginationResult[models.CatalogueItem], err error) {

	session, _ := helpers.NewNeo4jSession(*svc.neo4jDriver)

	//get all categories by parent path
	query := CatalogueItemsFiltersPaginationQuery(search, categoryPath, pageSize*(page-1), pageSize)
	items, err := helpers.GetNeo4jArrayOfNodes[models.CatalogueItem](session, query)
	totalCount, _ := helpers.GetNeo4jSingleRecordSingleValue[int64](session, CatalogueItemsFiltersTotalCountQuery(search, categoryPath))

	result = helpers.GetPaginationResult(items, totalCount, err)

	return result, err
}

func (svc *CatalogueService) GetCatalogueItemWithDetailsByUid(uid string) (result models.CatalogueItem, err error) {
	session, _ := helpers.NewNeo4jSession(*svc.neo4jDriver)

	query := CatalogueItemWithDetailsByUidQuery(uid)
	result, err = helpers.GetNeo4jSingleRecordAndMapToStruct[models.CatalogueItem](session, query)

	return result, err
}

func (svc *CatalogueService) GetCatalogueCategoryWithDetailsByUid(uid string) (result models.CatalogueCategory, err error) {
	session, _ := helpers.NewNeo4jSession(*svc.neo4jDriver)

	query := CatalogueCategoryWithDetailsQuery(uid)
	result, err = helpers.GetNeo4jSingleRecordAndMapToStruct[models.CatalogueCategory](session, query)

	return result, err
}

func (svc *CatalogueService) UpdateCatalogueCategory(catalogueCategory *models.CatalogueCategory) (err error) {
	session, _ := helpers.NewNeo4jSession(*svc.neo4jDriver)

	//get the original record from db to compare because of the delete
	originalItem, err := svc.GetCatalogueCategoryWithDetailsByUid(catalogueCategory.UID)
	if err == nil {
		//update category query
		query := UpdateCatalogueCategoryQuery(catalogueCategory, &originalItem)
		_, err = helpers.WriteNeo4jAndReturnSingleValue[string](session, query)
	}

	return err
}

func (svc *CatalogueService) CreateCatalogueCategory(catalogueCategory *models.CatalogueCategory) (err error) {
	session, _ := helpers.NewNeo4jSession(*svc.neo4jDriver)

	//use this hybrid update/create method
	query := UpdateCatalogueCategoryQuery(catalogueCategory, nil)
	_, err = helpers.WriteNeo4jAndReturnSingleValue[string](session, query)

	return err
}

func (svc *CatalogueService) GetCatalogueCategoryImageByUid(uid string) (result string, err error) {

	session, _ := helpers.NewNeo4jSession(*svc.neo4jDriver)

	query := CatalogueCategoryImageByUidQuery(uid)
	result, err = helpers.GetNeo4jSingleRecordSingleValue[string](session, query)

	return result, err
}

func (svc *CatalogueService) DeleteCatalogueCategory(uid string) (err error) {

	session, _ := helpers.NewNeo4jSession(*svc.neo4jDriver)

	query := DeleteCatalogueCategoryByUidQuery(uid)
	err = helpers.WriteNeo4jAndReturnNothing(session, query)

	return err
}

func (svc *CatalogueService) GetZonesCodebook() (result []codebookModels.Codebook, err error) {

	session, _ := helpers.NewNeo4jSession(*svc.neo4jDriver)

	query := GetZonesCodebookQuery()
	result, err = helpers.GetNeo4jArrayOfNodes[codebookModels.Codebook](session, query)

	return result, err
}

func (svc *CatalogueService) GetSubZonesCodebook(parentUID string) (result []codebookModels.Codebook, err error) {

	session, _ := helpers.NewNeo4jSession(*svc.neo4jDriver)

	query := GetSubZonesByParentUidCodebookQuery(parentUID)
	result, err = helpers.GetNeo4jArrayOfNodes[codebookModels.Codebook](session, query)

	return result, err
}

func (svc *CatalogueService) GetUnitsCodebook() (result []codebookModels.Codebook, err error) {

	session, _ := helpers.NewNeo4jSession(*svc.neo4jDriver)

	query := GetUnitsCodebookQuery()
	result, err = helpers.GetNeo4jArrayOfNodes[codebookModels.Codebook](session, query)

	return result, err
}

func (svc *CatalogueService) GetPropertyTypesCodebook() (result []codebookModels.Codebook, err error) {

	session, _ := helpers.NewNeo4jSession(*svc.neo4jDriver)

	query := GetPropertyTypesCodebookQuery()
	result, err = helpers.GetNeo4jArrayOfNodes[codebookModels.Codebook](session, query)

	return result, err
}
