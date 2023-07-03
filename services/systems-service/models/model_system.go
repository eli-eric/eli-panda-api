package models

import (
	catalogueModels "panda/apigateway/services/catalogue-service/models"
	"panda/apigateway/services/codebook-service/models"
)

type System struct {
	UID           string                 `json:"uid" neo4j:"ignore"`
	Name          string                 `json:"name" neo4j:"prop,name"`
	ParentPath    []SystemSimpleResponse `json:"parentPath" neo4j:"ignore"`
	ParentUID     *string                `json:"parentUid,omitempty" neo4j:"ignore"`
	Description   *string                `json:"description,omitempty" neo4j:"prop,description"`
	SystemType    *models.Codebook       `json:"systemType,omitempty" neo4j:"rel,SystemType,HAS_SYSTEM_TYPE,systemType"`
	SystemCode    *string                `json:"systemCode,omitempty" neo4j:"prop,systemCode"`
	SystemAlias   *string                `json:"systemAlias,omitempty" neo4j:"prop,systemAlias"`
	Owner         *models.Codebook       `json:"owner,omitempty" neo4j:"rel,Employee,HAS_OWNER,owner"`
	Responsible   *models.Codebook       `json:"responsible,omitempty" neo4j:"rel,Employee,HAS_RESPONSIBLE,responsible"`
	Importance    *models.Codebook       `json:"importance,omitempty" neo4j:"rel,SystemImportance,HAS_IMPORTANCE,importance"`
	Zone          *models.Codebook       `json:"zone,omitempty" neo4j:"rel,Zone,HAS_ZONE,zone"`
	Location      *models.Codebook       `json:"location,omitempty" neo4j:"ignore"`
	PhysicalItem  *PhysicalItem          `json:"physicalItem,omitempty" neo4j:"ignore"`
	HasSubsystems bool                   `json:"hasSubsystems" neo4j:"ignore"`
	Statistics    *SystemStatistics      `json:"statistics,omitempty" neo4j:"ignore"`
}

type PhysicalItem struct {
	UID           string                        `json:"uid"`
	ItemUsage     *models.Codebook              `json:"itemUsage,omitempty"`
	Price         *any                          `json:"price,omitempty"`
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

// {
// 	"0": {
// 		"direction": "to",
// 		"relationTypeCode": "IS_SPARE_FOR",
// 		"foreignSystemName": "Switchable radical protocol",
// 		"relationUid": "bd6f9ad8-28f3-4a2c-a499-9d7519d28e1d"
// 	}
// }

type SystemRelationship struct {
	Direction         string `json:"direction"`
	RelationTypeCode  string `json:"relationTypeCode"`
	ForeignSystemName string `json:"foreignSystemName"`
	RelationUID       string `json:"relationUid"`
}
