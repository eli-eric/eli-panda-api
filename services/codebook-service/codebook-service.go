package codebookService

import (
	"panda/apigateway/config"
	catalogueService "panda/apigateway/services/catalogue-service"
	"panda/apigateway/services/codebook-service/models"
	ordersService "panda/apigateway/services/orders-service"
	securityService "panda/apigateway/services/security-service"
	systemsService "panda/apigateway/services/systems-service"
)

type CodebookService struct {
	catalogueService catalogueService.ICatalogueService
	securityService  securityService.ISecurityService
	systemsService   systemsService.ISystemsService
	ordersService    ordersService.IOrdersService
}

type ICodebookService interface {
	GetCodebook(codebookCode string, parentUID string, facilityCode string) (codebookList []models.Codebook, err error)
	GetAutocompleteCodebook(codebookCode string, searchString string, limit int, facilityCode string) (codebookList []models.Codebook, err error)
}

// Create new security service instance
func NewCodebookService(settings *config.Config,
	catalogueService catalogueService.ICatalogueService,
	securityService securityService.ISecurityService,
	systemsService systemsService.ISystemsService,
	orderService ordersService.IOrdersService) ICodebookService {

	return &CodebookService{catalogueService: catalogueService, securityService: securityService, systemsService: systemsService, ordersService: orderService}
}

func (svc *CodebookService) GetCodebook(codebookCode string, parentUID string, facilityCode string) (codebookList []models.Codebook, err error) {

	codebookList = make([]models.Codebook, 0)

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
	case "USER":
		codebookList, err = svc.securityService.GetUsersCodebook(facilityCode)
	case "ORDER_STATUS":
		codebookList, err = svc.ordersService.GetOrderStatusesCodebook()
	}

	return codebookList, err
}

func (svc *CodebookService) GetAutocompleteCodebook(codebookCode string, searchString string, limit int, facilityCode string) (codebookList []models.Codebook, err error) {

	codebookList = make([]models.Codebook, 0)

	switch codebookCode {
	case "LOCATION":
		codebookList, err = svc.systemsService.GetLocationAutocompleteCodebook(searchString, limit, facilityCode)
	case "USER":
		codebookList, err = svc.securityService.GetUsersAutocompleteCodebook(searchString, limit, facilityCode)
	case "SUPPLIER":
		codebookList, err = svc.ordersService.GetSuppliersAutocompleteCodebook(searchString, limit)
	case "SYSTEM":
		codebookList, err = svc.systemsService.GetSystemsAutocompleteCodebook(searchString, limit, facilityCode)
	}

	return codebookList, err
}
