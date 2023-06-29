package models

import (
	catalogueModels "panda/apigateway/services/catalogue-service/models"
	"panda/apigateway/services/codebook-service/models"
)

type System struct {
	UID           string                 `json:"uid" neo4j:"uid"`
	Name          string                 `json:"name"`
	ParentPath    []SystemSimpleResponse `json:"parentPath"`
	ParentUID     *string                `json:"parentUid,omitempty"`
	Description   *string                `json:"description,omitempty"`
	SystemType    *models.Codebook       `json:"systemType,omitempty"`
	SystemCode    *string                `json:"systemCode,omitempty"`
	SystemAlias   *string                `json:"systemAlias,omitempty"`
	Owner         *models.Codebook       `json:"owner,omitempty"`
	Responsible   *models.Codebook       `json:"responsible,omitempty"`
	Importance    *models.Codebook       `json:"importance,omitempty"`
	Zone          *models.Codebook       `json:"zone,omitempty"`
	Location      *models.Codebook       `json:"location,omitempty"`
	PhysicalItem  *PhysicalItem          `json:"physicalItem,omitempty"`
	HasSubsystems bool                   `json:"hasSubsystems"`
	Statistics    *SystemStatistics      `json:"statistics,omitempty"`
}

type PhysicalItem struct {
	UID           string                        `json:"uid"`
	ItemUsage     *models.Codebook              `json:"itemUsage,omitempty"`
	Price         *float64                      `json:"price,omitempty"`
	Currency      *string                       `json:"currency,omitempty"`
	EUN           *string                       `json:"eun,omitempty"`
	SerialNumber  *string                       `json:"serialNumber,omitempty"`
	CatalogueItem catalogueModels.CatalogueItem `json:"catalogueItem"`
}

type SystemSimpleResponse struct {
	UID  string `json:"uid"`
	Name string `json:"name"`
}

type SystemStatistics struct {
	SubsystemsCount int `json:"subsystemsCount"`
	SparePartsCount int `json:"sparePartsCount"`
}
