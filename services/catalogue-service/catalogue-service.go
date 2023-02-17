package catalogueService

import (
	"errors"
	"fmt"
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
	GetCataloguecategoriesByParentPath(parentPath string) (categories []models.CatalogueCategory, err error)
}

// Create new security service instance
func NewCatalogueService(settings *config.Config, driver *neo4j.Driver) ICatalogueService {

	return &CatalogueService{neo4jDriver: driver, jwtSecret: settings.JwtSecret}
}

func (svc *CatalogueService) GetCataloguecategoriesByParentPath(parentPath string) (categories []models.CatalogueCategory, err error) {

	// Open a new Session

	session, _ := helpers.NewNeo4jSession(*svc.neo4jDriver)

	//the user has to be enabled and has at least one role
	categoryNodes, err := helpers.GetNeo4jArrayOfNodes(session, `
		match(category:CatalogueCategory)
		with category
		optional match(parent)-[:HAS_SUBCATEGORY*1..50]->(category) 
		with category, apoc.text.join(reverse(collect(parent.code)),"/") as parentPath
		where parentPath = $parentPath
		return {uid:category.uid,code:category.code, name:category.name,parentPath: parentPath} as categories
	`, map[string]interface{}{"parentPath": parentPath}, "categories")

	if err == nil {

		fmt.Println(categoryNodes)

		arr := categoryNodes.([]map[string]interface{})

		for i := 0; i < len(arr); i++ {
			catItem, err := helpers.MapStruct[models.CatalogueCategory](arr[i])
			if err == nil {
				categories = append(categories, catItem)
			}
		}

		return categories, err
	}

	fmt.Println(err)

	return categories, errors.New("Unauthorized")
}

// // get items with search and pagination
// `MATCH(itm:CatalogueItem)
// OPTIONAL MATCH(itm)-[:HAS_MANUFACTURER]->(manu)
// OPTIONAL MATCH(itm)-[propVal:HAS_CATALOGUE_PROPERTY]->(prop)
// OPTIONAL MATCH(prop)-[:HAS_UNIT]->(unit)
// OPTIONAL MATCH(prop)-[:IS_PROPERTY_TYPE]->(propType)
// OPTIONAL MATCH(group)-[:CONTAINS_PROPERTY]->(prop)
// WITH itm, propType.code as propTypeCode, manu.name as manufacturerName, prop.name as propName, group.name as groupName, propVal.value as value, unit.name as unit
// ORDER BY itm.name
// RETURN {
//     uid: itm.uid,
//     name: itm.name,
//     description: itm.description,
//     manufacturer: manufacturerName,
//     manufacturerUrl: itm.manufacturerUrl,
//     manufacturerNumber: itm.catalogueNumber,
//     details: collect({ propertyName: propName, propertyType: propTypeCode,propertyUnit: unit, propertyGroup: groupName, value: value})
// } as items
// SKIP 5 * 6
// LIMIT 5`
