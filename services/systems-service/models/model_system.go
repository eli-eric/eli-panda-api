package models

import (
	catalogueModels "panda/apigateway/services/catalogue-service/models"
	"panda/apigateway/services/codebook-service/models"
	"time"
)

type System struct {
	UID           string                 `json:"uid" neo4j:"ignore"`
	Name          string                 `json:"name" neo4j:"prop,name"`
	ParentPath    []SystemSimpleResponse `json:"parentPath" neo4j:"ignore"`
	ParentUID     *string                `json:"parentUid,omitempty" neo4j:"ignore"`
	Description   *string                `json:"description,omitempty" neo4j:"prop,description"`
	SystemType    *models.Codebook       `json:"systemType,omitempty" neo4j:"rel,SystemType,HAS_SYSTEM_TYPE,uid,systemType"`
	SystemCode    *string                `json:"systemCode,omitempty" neo4j:"prop,systemCode"`
	SystemAlias   *string                `json:"systemAlias,omitempty" neo4j:"prop,systemAlias"`
	SystemLevel   *string                `json:"systemLevel,omitempty" neo4j:"prop,systemLevel"`
	Owner         *models.Codebook       `json:"owner,omitempty" neo4j:"rel,Employee,HAS_OWNER,uid,owner"`
	Responsible   *models.Codebook       `json:"responsible,omitempty" neo4j:"rel,Employee,HAS_RESPONSIBLE,uid,responsible"`
	Importance    *models.Codebook       `json:"importance,omitempty" neo4j:"rel,SystemImportance,HAS_IMPORTANCE,uid,importance"`
	Zone          *models.Codebook       `json:"zone,omitempty" neo4j:"rel,Zone,HAS_ZONE,uid,zone"`
	Location      *models.Codebook       `json:"location,omitempty" neo4j:"ignore"`
	PhysicalItem  *PhysicalItem          `json:"physicalItem,omitempty" neo4j:"ignore"`
	HasSubsystems bool                   `json:"hasSubsystems" neo4j:"ignore"`
	Statistics    *SystemStatistics      `json:"statistics,omitempty" neo4j:"ignore"`
	SparesOut     int                    `json:"sparesOut" neo4j:"ignore"`
	SparesIn      int                    `json:"sparesIn" neo4j:"ignore"`
	History       *[]SystemHistory       `json:"history,omitempty" neo4j:"ignore"`
	MiniImageUrl  *[]string              `json:"miniImageUrl" neo4j:"ignore"`
	SubSystems    *[]System              `json:"subSystems,omitempty" neo4j:"ignore"`
}

type SystemTreeUid struct {
	UID      string           `json:"uid"`
	Children *[]SystemTreeUid `json:"children"`
}

type PhysicalItem struct {
	UID           string                        `json:"uid"`
	ItemUsage     *models.Codebook              `json:"itemUsage,omitempty"`
	Price         *any                          `json:"price,omitempty"`
	Currency      *string                       `json:"currency,omitempty"`
	EUN           *string                       `json:"eun,omitempty"`
	SerialNumber  *string                       `json:"serialNumber,omitempty"`
	CatalogueItem catalogueModels.CatalogueItem `json:"catalogueItem"`
	OrderNumber   *string                       `json:"orderNumber,omitempty" neo4j:"prop,orderNumber"`
	OrderUid      *string                       `json:"orderUid,omitempty" neo4j:"prop,orderUid"`
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
	SubsystemsCount        int      `json:"subsystemsCount"`
	SparePartsCount        int      `json:"sparePartsCount"`
	MinimalSpareParstCount *float32 `json:"minimalSpareParstCount" neo4j:"ignore"`
	SparePartsCoverageSum  *float32 `json:"sparePartsCoverageSum" neo4j:"ignore"`
	Sp_coverage            *float32 `json:"sp_coverage" neo4j:"ignore"`
}

// SystemCodesResult is a simplified view of systems for Control Systems users.
// It includes only system name/code plus minimal navigation and audit info.
type SystemCodesResult struct {
	UID          string                 `json:"uid"`
	Name         string                 `json:"name"`
	Code         string                 `json:"code"`
	ParentPath   []SystemSimpleResponse `json:"parentPath"`
	CreatedBy    string                 `json:"createdBy"`
	LastUpdateBy string                 `json:"lastUpdateBy"`
	Zone         *models.Codebook       `json:"zone,omitempty"`
	Location     *models.Codebook       `json:"location,omitempty"`
}

type SystemCodesRequest struct {
	SystemType *models.Codebook `json:"systemType"`
	Zone       *models.Codebook `json:"zone"`
	Batch      int              `json:"batch"`
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
	UID             string           `json:"uid"`
	Name            string           `json:"name"`
	Code            string           `json:"code"`
	Mask            string           `json:"mask"`
	SystemAttribute *models.Codebook `json:"systemAttribute,omitempty" neo4j:"rel,SystemAttribute,HAS_SYSTEM_ATTRIBUTE,uid,systemAttribute"`
}

// SystemTypeTreeItem represents a system type in the tree response
type SystemTypeTreeItem struct {
	UID  string `json:"uid"`
	Name string `json:"name"`
	Code string `json:"code"`
}

// SystemTypeGroupTreeItem represents a system type group with its children in tree response
type SystemTypeGroupTreeItem struct {
	UID      string               `json:"uid"`
	Name     string               `json:"name"`
	Code     string               `json:"code"`
	Children []SystemTypeTreeItem `json:"children"`
}

type SystemWithAllDetails struct {
	System
	ParentSystem SystemSimpleResponse `json:"parentSystem"`
}

type EUN struct {
	EUN string `json:"eun"`
}

type EunLocation struct {
	EUN         string `json:"eun"`
	LocationUID string `json:"location_uid"`
}

type SystemCodeRequest struct {
	ParentUID     string `json:"parentUid"`
	SystemTypeUID string `json:"systemTypeUid"`
	ZoneUID       string `json:"zoneUid"`
}

type PhysicalItemMovement struct {
	SourceSystemUID      string           `json:"sourceSystemUid"`
	ParentSystemUID      string           `json:"parentSystemUid"`
	DestinationSystemUID string           `json:"destinationSystemUid"`
	SystemName           string           `json:"systemName"`
	Location             *models.Codebook `json:"location"`
	ItemUsage            *models.Codebook `json:"itemUsage"`
	Condition            *models.Codebook `json:"condition"`
	DeleteSourceSystem   bool             `json:"deleteSourceSystem"`
}

type SystemsMovement struct {
	TargetParentSystemUid string   `json:"targetParentSystemUid"`
	SystemsToMoveUids     []string `json:"systemsToMoveUids"`
}

type SystemCopyRequest struct {
	SourceSystemUID              string `json:"sourceSystemUid"`
	DestinationSystemUID         string `json:"destinationSystemUid"`
	CopyOnlySourceSystemChildren bool   `json:"copyOnlySourceSystemChildren"`
	CopyRecursive                bool   `json:"copyRecursive"`
}

type JiraSystemImportRequest struct {
	Name            string `json:"name"`
	Code            string `json:"code"`
	ParentSystemUID string `json:"parentSystemUid"`
	Description     string `json:"description"`
	LinkUrl         string `json:"linkUrl"`
	LinkName        string `json:"linkName"`
	ZoneUID         string `json:"zoneUid"`
	SystemTypeUID   string `json:"systemTypeUid"`
}

type SystemPhysicalItemInfo struct {
	SystemUid  string `json:"systemUid"`
	SystemName string `json:"systemName"`
	ItemUid    string `json:"itemUid"`
	ItemName   string `json:"itemName"`
}

type AssignSpareRequest struct {
	SpareItemUid       string          `json:"spareItemUid" binding:"required"`
	SystemUid          string          `json:"systemUid" binding:"required"`
	OldItemCondition   models.Codebook `json:"oldItemCondition" binding:"required"`
	NewParentSystemUid string          `json:"newParentSystemUid,omitempty"`
	NewItemLocation    models.Codebook `json:"newItemLocation" binding:"required"`
}

type AssignSpareResponse struct {
	Success            bool   `json:"success"`
	Message            string `json:"message"`
	UpdatedSystemUid   string `json:"updatedSystemUid"`
	RelocatedItemUid   string `json:"relocatedItemUid,omitempty"`
	NewParentSystemUid string `json:"newParentSystemUid,omitempty"`
}

// SystemSparePartsDetail contains comprehensive system and physical item information with all spare relations
type SystemSparePartsDetail struct {
	System              SystemDetailInfo             `json:"system"`
	Location            *CodebookInfo                `json:"location,omitempty"`
	Zone                *CodebookInfo                `json:"zone,omitempty"`
	SystemType          *CodebookInfo                `json:"systemType,omitempty"`
	SystemAttributes    []CodebookInfo               `json:"systemAttributes"`
	ResponsiblePersons  ResponsiblePersonsInfo       `json:"responsiblePersons"`
	Team                *CodebookInfo                `json:"team,omitempty"`
	PhysicalItems       []PhysicalItemDetailExtended `json:"physicalItems"`
	SparePartsRelations SparePartsRelationsInfo      `json:"sparePartsRelations"`
}

// CodebookInfo represents simplified codebook information
type CodebookInfo struct {
	UID  string  `json:"uid"`
	Name string  `json:"name"`
	Code *string `json:"code,omitempty"`
}

// SystemDetailInfo contains detailed system information
type SystemDetailInfo struct {
	UID                    string  `json:"uid"`
	Name                   string  `json:"name"`
	SystemCode             *string `json:"systemCode,omitempty"`
	Description            *string `json:"description,omitempty"`
	Status                 *string `json:"status,omitempty"`
	IsTechnologicalUnit    *bool   `json:"isTechnologicalUnit,omitempty"`
	IsCritical             *bool   `json:"isCritical,omitempty"`
	MinimalSpareParstCount *int64  `json:"minimalSpareParstCount,omitempty"`
	SparePartsCoverageSum  *int64  `json:"sparePartsCoverageSum,omitempty"`
	SystemLevel            *string `json:"systemLevel,omitempty"`
	SystemAlias            *string `json:"systemAlias,omitempty"`
	Image                  *string `json:"image,omitempty"`
	MiniImageUrl           *string `json:"miniImageUrl,omitempty"`
	LastUpdateTime         *string `json:"lastUpdateTime,omitempty"`
	LastUpdateBy           *string `json:"lastUpdateBy,omitempty"`
}

// ResponsiblePersonsInfo contains information about responsible persons
type ResponsiblePersonsInfo struct {
	Responsible *EmployeeInfo `json:"responsible,omitempty"`
	Operator    *EmployeeInfo `json:"operator,omitempty"`
	Maintainer  *EmployeeInfo `json:"maintainer,omitempty"`
}

// EmployeeInfo contains employee information
type EmployeeInfo struct {
	UID       string  `json:"uid"`
	FirstName *string `json:"firstName,omitempty"`
	LastName  *string `json:"lastName,omitempty"`
	Email     *string `json:"email,omitempty"`
	Phone     *string `json:"phone,omitempty"`
}

// PhysicalItemDetailExtended contains extended physical item information
type PhysicalItemDetailExtended struct {
	UID            string             `json:"uid"`
	Name           string             `json:"name"`
	SerialNumber   *string            `json:"serialNumber,omitempty"`
	EUN            *string            `json:"eun,omitempty"`
	Price          *string            `json:"price,omitempty"`
	Currency       *string            `json:"currency,omitempty"`
	Status         *string            `json:"status,omitempty"`
	Notes          *string            `json:"notes,omitempty"`
	PrintEUN       *bool              `json:"printEUN,omitempty"`
	LastUpdateTime *string            `json:"lastUpdateTime,omitempty"`
	Condition      *CodebookInfo      `json:"condition,omitempty"`
	Usage          *CodebookInfo      `json:"usage,omitempty"`
	CatalogueItem  *CatalogueItemInfo `json:"catalogueItem,omitempty"`
}

// CatalogueItemInfo contains catalogue item information
type CatalogueItemInfo struct {
	UID             string        `json:"uid"`
	Name            string        `json:"name"`
	Description     *string       `json:"description,omitempty"`
	CatalogueNumber *string       `json:"catalogueNumber,omitempty"`
	Image           *string       `json:"image,omitempty"`
	MiniImageUrl    *string       `json:"miniImageUrl,omitempty"`
	ManufacturerUrl *string       `json:"manufacturerUrl,omitempty"`
	LastUpdateTime  *string       `json:"lastUpdateTime,omitempty"`
	Category        *CodebookInfo `json:"category,omitempty"`
	Manufacturer    *CodebookInfo `json:"manufacturer,omitempty"`
}

// SparePartsRelationsInfo contains spare parts relationship information
type SparePartsRelationsInfo struct {
	SpareSystems  []SpareSystemInfo  `json:"spareSystems"`
	ParentSystems []SystemSimpleInfo `json:"parentSystems"`
}

// SpareSystemInfo contains spare system information
type SpareSystemInfo struct {
	UID           string                  `json:"uid"`
	Name          string                  `json:"name"`
	SystemCode    *string                 `json:"systemCode,omitempty"`
	Description   *string                 `json:"description,omitempty"`
	Status        *string                 `json:"status,omitempty"`
	PhysicalItems []SparePhysicalItemInfo `json:"physicalItems"`
}

// SparePhysicalItemInfo contains spare physical item information
type SparePhysicalItemInfo struct {
	UID           string                   `json:"uid"`
	Name          string                   `json:"name"`
	SerialNumber  *string                  `json:"serialNumber,omitempty"`
	EUN           *string                  `json:"eun,omitempty"`
	CatalogueItem *SimpleCatalogueItemInfo `json:"catalogueItem,omitempty"`
}

// SimpleCatalogueItemInfo contains simplified catalogue item information
type SimpleCatalogueItemInfo struct {
	UID             string  `json:"uid"`
	Name            string  `json:"name"`
	CatalogueNumber *string `json:"catalogueNumber,omitempty"`
}

// SystemSimpleInfo contains simplified system information
type SystemSimpleInfo struct {
	UID         string  `json:"uid"`
	Name        string  `json:"name"`
	SystemCode  *string `json:"systemCode,omitempty"`
	Description *string `json:"description,omitempty"`
	Status      *string `json:"status,omitempty"`
}
