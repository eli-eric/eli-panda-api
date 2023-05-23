package codebookService

import (
	"panda/apigateway/config"
	"panda/apigateway/helpers"
	catalogueService "panda/apigateway/services/catalogue-service"
	"panda/apigateway/services/codebook-service/models"
	ordersService "panda/apigateway/services/orders-service"
	securityService "panda/apigateway/services/security-service"
	systemsService "panda/apigateway/services/systems-service"

	"github.com/google/uuid"
	"github.com/neo4j/neo4j-go-driver/v4/neo4j"
)

type CodebookService struct {
	neo4jDriver      *neo4j.Driver
	catalogueService catalogueService.ICatalogueService
	securityService  securityService.ISecurityService
	systemsService   systemsService.ISystemsService
	ordersService    ordersService.IOrdersService
}

type ICodebookService interface {
	GetCodebook(codebookCode string, parentUID string, facilityCode string) (codebookResponse models.CodebookResponse, err error)
	GetAutocompleteCodebook(codebookCode string, searchString string, limit int, facilityCode string) (CodebookResponse models.CodebookResponse, err error)
	GetListOfCodebooks() (codebookList []models.CodebookType)
	CreateNewCodebook(codebookCode string, facilityCode string, userUID string, userRoles []string, codebook *models.Codebook) (result models.Codebook, err error)
}

// Create new security service instance
func NewCodebookService(settings *config.Config,
	driver *neo4j.Driver,
	catalogueService catalogueService.ICatalogueService,
	securityService securityService.ISecurityService,
	systemsService systemsService.ISystemsService,
	orderService ordersService.IOrdersService) ICodebookService {

	return &CodebookService{neo4jDriver: driver, catalogueService: catalogueService, securityService: securityService, systemsService: systemsService, ordersService: orderService}
}

func (svc *CodebookService) GetCodebook(codebookCode string, parentUID string, facilityCode string) (codebookResponse models.CodebookResponse, err error) {

	codebookList := make([]models.Codebook, 0)
	codebookMetadata := codebooksMap[codebookCode]

	if codebookMetadata != (models.CodebookType{}) {

		switch codebookCode {
		case "ZONE":
			codebookList, err = svc.systemsService.GetZonesCodebook(facilityCode)
		case "UNIT":
			codebookList, err = svc.catalogueService.GetUnitsCodebook()
		case "CATALOGUE_PROPERTY_TYPE":
			codebookList, err = svc.catalogueService.GetPropertyTypesCodebook()
		case "SYSTEM_TYPE":
			codebookList, err = svc.systemsService.GetSystemTypesCodebook(facilityCode)
		case "SYSTEM_IMPORTANCE":
			codebookList, err = svc.systemsService.GetSystemImportancesCodebook()
		case "SYSTEM_CRITICALITY_CLASS":
			codebookList, err = svc.systemsService.GetSystemCriticalitiesCodebook()
		case "ITEM_USAGE":
			codebookList, err = svc.systemsService.GetItemUsagesCodebook()
		case "ITEM_CONDITION_STATUS":
			codebookList, err = svc.systemsService.GetItemConditionsCodebook()
		case "ORDER_STATUS":
			codebookList, err = svc.ordersService.GetOrderStatusesCodebook()
		case "PROCUREMENTER":
			codebookList, err = svc.securityService.GetProcurementersCodebook(facilityCode)
		}

		if err == nil {
			codebookResponse = models.CodebookResponse{Metadata: codebookMetadata, Data: codebookList}
		}

	} else {
		err = helpers.ERR_NOT_FOUND
	}

	return codebookResponse, err
}

func (svc *CodebookService) GetAutocompleteCodebook(codebookCode string, searchString string, limit int, facilityCode string) (codebookResponse models.CodebookResponse, err error) {

	codebookList := make([]models.Codebook, 0)
	codebookMetadata := codebooksMap[codebookCode]

	if codebookMetadata != (models.CodebookType{}) {

		switch codebookCode {
		case "LOCATION":
			codebookList, err = svc.systemsService.GetLocationAutocompleteCodebook(searchString, limit, facilityCode)
		case "USER":
			codebookList, err = svc.securityService.GetUsersAutocompleteCodebook(searchString, limit, facilityCode)
		case "SUPPLIER":
			codebookList, err = svc.ordersService.GetSuppliersAutocompleteCodebook(searchString, limit)
		case "SYSTEM":
			codebookList, err = svc.systemsService.GetSystemsAutocompleteCodebook(searchString, limit, facilityCode)
		case "EMPLOYEE":
			codebookList, err = svc.securityService.GetEmployeesAutocompleteCodebook(searchString, limit, facilityCode)
		}

		if err == nil {
			codebookResponse = models.CodebookResponse{Metadata: codebookMetadata, Data: codebookList}
		}

	} else {
		err = helpers.ERR_NOT_FOUND
	}

	return codebookResponse, err
}

func (svc *CodebookService) GetListOfCodebooks() (codebookList []models.CodebookType) {

	return []models.CodebookType{
		models.ZONE_CODEBOOK,
		models.UNIT_CODEBOOK,
		models.CATALOGUE_PROPERTY_TYPE_CODEBOOK,
		models.SYSTEM_TYPE_CODEBOOK,
		models.SYSTEM_IMPORTANCE_CODEBOOK,
		models.SYSTEM_CRITICALITY_CLASS_CODEBOOK,
		models.ITEM_USAGE_CODEBOOK,
		models.ITEM_CONDITION_STATUS_CODEBOOK,
		models.USER_CODEBOOK,
		models.ORDER_STATUS_CODEBOOK,
		models.LOCATION_AUTOCOMPLETE_CODEBOOK,
		models.EMPLOYEE_AUTOCOMPLETE_CODEBOOK,
		models.SYSTEM_AUTOCOMPLETE_CODEBOOK,
		models.USER_AUTOCOMPLETE_CODEBOOK,
		models.SUPPLIER_AUTOCOMPLETE_CODEBOOK,
	}
}

func (svc *CodebookService) CreateNewCodebook(codebookCode string, facilityCode string, userUID string, userRoles []string, codebook *models.Codebook) (result models.Codebook, err error) {

	codebookDefinition := codebooksMap[codebookCode]

	if codebookDefinition != (models.CodebookType{}) {
		if checkUserRoles(userRoles, codebookDefinition.RoleEdit) {

			if codebookDefinition.NodeLabel != "" {
				// Open a new Session
				session, _ := helpers.NewNeo4jSession(*svc.neo4jDriver)

				dbquery := helpers.DatabaseQuery{}
				dbquery.Query = `				
				CREATE (n:` + codebookDefinition.NodeLabel + ` {uid: $uid, name: $name }) 
					WITH n
					MATCH (u:User {uid: $userUID})
					CREATE (n)-[:WAS_UPDATED_BY{ at: datetime(), action: "INSERT" }]->(u)
				RETURN { uid: n.uid, name: n.name } as codebook`

				dbquery.Parameters = map[string]interface{}{
					"uid":     uuid.New().String(),
					"name":    codebook.Name,
					"userUID": userUID,
				}
				dbquery.ReturnAlias = "codebook"

				result, err = helpers.WriteNeo4jReturnSingleRecordAndMapToStruct[models.Codebook](session, dbquery)
			}
		} else {
			err = helpers.ERR_UNAUTHORIZED
		}
	}

	return result, err
}

func checkUserRoles(userRoles []string, role string) (result bool) {
	for _, userRole := range userRoles {
		if userRole == role {
			return true
		}
	}
	return false
}

var codebooksMap = map[string]models.CodebookType{
	"ZONE":                       models.ZONE_CODEBOOK,
	"UNIT":                       models.UNIT_CODEBOOK,
	"CATALOGUE_PROPERTY_TYPE":    models.CATALOGUE_PROPERTY_TYPE_CODEBOOK,
	"SYSTEM_TYPE":                models.SYSTEM_TYPE_CODEBOOK,
	"SYSTEM_IMPORTANCE":          models.SYSTEM_IMPORTANCE_CODEBOOK,
	"SYSTEM_CRITICALITY_CLASS":   models.SYSTEM_CRITICALITY_CLASS_CODEBOOK,
	"ITEM_USAGE":                 models.ITEM_USAGE_CODEBOOK,
	"ITEM_CONDITION_STATUS":      models.ITEM_CONDITION_STATUS_CODEBOOK,
	"USER":                       models.USER_CODEBOOK,
	"ORDER_STATUS":               models.ORDER_STATUS_CODEBOOK,
	"LOCATION":                   models.LOCATION_AUTOCOMPLETE_CODEBOOK,
	"EMPLOYEE":                   models.EMPLOYEE_AUTOCOMPLETE_CODEBOOK,
	"SYSTEM":                     models.SYSTEM_AUTOCOMPLETE_CODEBOOK,
	"USER_AUTOCOMPLETE_CODEBOOK": models.USER_AUTOCOMPLETE_CODEBOOK,
	"SUPPLIER":                   models.SUPPLIER_AUTOCOMPLETE_CODEBOOK,
}
