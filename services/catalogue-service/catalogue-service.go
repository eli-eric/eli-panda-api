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
	query := CatalogueItemsPaginationFiltersQuery(search, categoryPath, pageSize*(page-1), pageSize)
	items, err := helpers.GetNeo4jArrayOfNodes[models.CatalogueItem](session, query)

	result = helpers.GetPaginationResult(items, err)

	return result, err
}
