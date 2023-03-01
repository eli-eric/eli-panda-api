package codebookService

import (
	"panda/apigateway/config"
	catalogueService "panda/apigateway/services/catalogue-service"
	"panda/apigateway/services/codebook-service/models"
	securityService "panda/apigateway/services/security-service"
)

type CodebookService struct {
	catalogueService catalogueService.ICatalogueService
	securityService  securityService.ISecurityService
}

type ICodebookService interface {
	GetCodebook(codebookCode string, parentUID string) (codebookList []models.Codebook, err error)
}

// Create new security service instance
func NewCodebookService(settings *config.Config, catalogueService catalogueService.ICatalogueService, securityService securityService.ISecurityService) ICodebookService {

	return &CodebookService{catalogueService: catalogueService, securityService: securityService}
}

func (svc *CodebookService) GetCodebook(codebookCode string, parentUID string) (codebookList []models.Codebook, err error) {

	codebookList = make([]models.Codebook, 0)

	switch codebookCode {
	case "ZONE":
		codebookList, err = svc.catalogueService.GetZonesCodebook()
	case "SUB_ZONE":
		codebookList, err = svc.catalogueService.GetSubZonesCodebook(parentUID)
	case "UNIT":
		codebookList, err = svc.catalogueService.GetUnitsCodebook()
	case "CATALOGUE_PROPERTY_TYPE":
		codebookList, err = svc.catalogueService.GetPropertyTypesCodebook()
	}

	return codebookList, err
}
