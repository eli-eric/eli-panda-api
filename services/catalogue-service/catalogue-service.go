package catalogueService

import (
	"panda/apigateway/config"
	"panda/apigateway/helpers"
	"panda/apigateway/services/catalogue-service/models"

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
