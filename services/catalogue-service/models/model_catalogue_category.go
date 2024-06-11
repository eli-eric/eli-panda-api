package models

import "panda/apigateway/services/codebook-service/models"

type CatalogueCategory struct {
	UID string `json:"uid,omitempty"`

	Name string `json:"name,omitempty"`

	Code string `json:"code,omitempty"`

	Image string `json:"image,omitempty"`

	ParentPath string `json:"parentPath,omitempty"`

	ParentUID string `json:"parentUID,omitempty"`

	Groups []CatalogueCategoryPropertyGroup `json:"groups"`

	PhysicalItemProperties []CatalogueCategoryProperty `json:"physicalItemProperties"`

	SystemType *models.Codebook `json:"systemType,omitempty"`

	MiniImageUrl *[]string `json:"miniImageUrl" neo4j:"ignore"`
}

type CatalogueCategoryPropertyGroup struct {
	UID string `json:"uid,omitempty"`

	Name string `json:"name,omitempty"`

	Properties []CatalogueCategoryProperty `json:"properties"`
}

type CatalogueCategoryProperty struct {
	UID string `json:"uid,omitempty"`

	Name string `json:"name,omitempty"`

	ListOfValues []string `json:"listOfValues,omitempty"`

	DefaultValue string `json:"defaultValue,omitempty"`

	Type CatalogueCategoryPropertyType `json:"type,omitempty"`

	Unit *models.Codebook `json:"unit,omitempty"`
}

type CatalogueCategoryPropertyType struct {
	UID  string `json:"uid,omitempty"`
	Name string `json:"name,omitempty"`
	Code string `json:"code,omitempty"`
}

type CatalogueCategoryTreeItem struct {
	UID             string                      `json:"uid,omitempty"`
	Has_subcategory []CatalogueCategoryTreeItem `json:"has_subcategory,omitempty"`
}
