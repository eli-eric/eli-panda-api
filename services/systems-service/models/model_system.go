package models

import (
	catalogueModels "panda/apigateway/services/catalogue-service/models"
	"panda/apigateway/services/codebook-service/models"
	"time"
)

type System struct {
	UID             string                 `json:"uid" neo4j:"ignore"`
	Name            string                 `json:"name" neo4j:"prop,name"`
	ParentPath      []SystemSimpleResponse `json:"parentPath" neo4j:"ignore"`
	ParentUID       *string                `json:"parentUid,omitempty" neo4j:"ignore"`
	Description     *string                `json:"description,omitempty" neo4j:"prop,description"`
	SystemType      *models.Codebook       `json:"systemType,omitempty" neo4j:"rel,SystemType,HAS_SYSTEM_TYPE,uid,systemType"`
	SystemCode      *string                `json:"systemCode,omitempty" neo4j:"prop,systemCode"`
	SystemAlias     *string                `json:"systemAlias,omitempty" neo4j:"prop,systemAlias"`
	SystemLevel     *string                `json:"systemLevel,omitempty" neo4j:"prop,systemLevel"`
	Owner           *models.Codebook       `json:"owner,omitempty" neo4j:"rel,Employee,HAS_OWNER,uid,owner"`
	Responsible     *models.Codebook       `json:"responsible,omitempty" neo4j:"rel,Employee,HAS_RESPONSIBLE,uid,responsible"`
	Importance      *models.Codebook       `json:"importance,omitempty" neo4j:"rel,SystemImportance,HAS_IMPORTANCE,uid,importance"`
	Zone            *models.Codebook       `json:"zone,omitempty" neo4j:"rel,Zone,HAS_ZONE,uid,zone"`
	Location        *models.Codebook       `json:"location,omitempty" neo4j:"ignore"`
	PhysicalItem    *PhysicalItem          `json:"physicalItem,omitempty" neo4j:"ignore"`
	HasSubsystems   bool                   `json:"hasSubsystems" neo4j:"ignore"`
	Statistics      *SystemStatistics      `json:"statistics,omitempty" neo4j:"ignore"`
	SparesOut       int                    `json:"sparesOut" neo4j:"ignore"`
	SparesIn        int                    `json:"sparesIn" neo4j:"ignore"`
	History         *[]SystemHistory       `json:"history,omitempty" neo4j:"ignore"`
	SystemAttribute *models.Codebook       `json:"systemAttribute,omitempty" neo4j:"rel,SystemAttribute,HAS_SYSTEM_ATTRIBUTE,uid,systemAttribute"`
	MiniImageUrl    *[]string              `json:"miniImageUrl" neo4j:"ignore"`
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

type PhysicalItemDetail struct {
	Property catalogueModels.CatalogueCategoryProperty `json:"property,omitempty"`

	Value any `json:"value,omitempty"`
}

type SystemSimpleResponse struct {
	UID  string `json:"uid"`
	Name string `json:"name"`
}

type SystemStatistics struct {
	SubsystemsCount int `json:"subsystemsCount"`
	SparePartsCount int `json:"sparePartsCount"`
}

type SystemRelationship struct {
	Direction         string `json:"direction"`
	RelationTypeCode  string `json:"relationTypeCode"`
	ForeignSystemName string `json:"foreignSystemName"`
	RelationUID       int64  `json:"relationUid"`
}

type SystemRelationshipRequest struct {
	RelationTypeCode string `json:"relationTypeCode"`
	SystemFromUID    string `json:"systemFromUid"`
	SystemToUID      string `json:"systemToUid"`
}

type SystemHistory struct {
	Uid         string               `json:"uid"`
	ChangedAt   time.Time            `json:"changedAt"`
	ChangedBy   string               `json:"changedBy"`
	HistoryType string               `json:"historyType"`
	Action      string               `json:"action"`
	Detail      *SystemHistoryDetail `json:"detail"`
}

type SystemHistoryDetail struct {
	SystemUid  string `json:"systemUid"`
	SystemName string `json:"systemName"`
	Direction  string `json:"direction"`
}

type SystemType struct {
	UID  string `json:"uid"`
	Name string `json:"name"`
	Code string `json:"code"`
	Mask string `json:"mask"`
}

type SystemWithAllDetails struct {
	System
	ParentSystem SystemSimpleResponse `json:"parentSystem"`
}

type EUN struct {
	EUN string `json:"eun"`
}
