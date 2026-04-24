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

	Order *int `json:"order,omitempty"`

	Properties []CatalogueCategoryProperty `json:"properties"`
}

type CatalogueCategoryProperty struct {
	UID string `json:"uid,omitempty"`

	Name string `json:"name,omitempty"`

	ListOfValues []string `json:"listOfValues,omitempty"`

	DefaultValue string `json:"defaultValue,omitempty"`

	Order *int `json:"order,omitempty"`

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

// PatchCatalogueCategoryFields carries the parsed body of PATCH /v1/catalogue/category/:uid.
// Each pointer is nil when the corresponding JSON key was absent.
type PatchCatalogueCategoryFields struct {
	Name       *string
	Code       *string
	SystemType *Optional[models.Codebook]
}

// PatchCatalogueCategoryGroupFields — PATCH /.../group/:gid payload.
type PatchCatalogueCategoryGroupFields struct {
	Name  *string
	Order *int
}

// CreateCatalogueCategoryGroupFields — POST /.../group payload.
type CreateCatalogueCategoryGroupFields struct {
	Name  string
	Order *int // optional; server auto-assigns if nil
}

// PatchCatalogueCategoryPropertyFields — PATCH /.../property/:pid payload.
// GroupUID is present only when the caller wants to move the property to a different group.
type PatchCatalogueCategoryPropertyFields struct {
	Name         *string
	DefaultValue *Optional[string]
	ListOfValues *[]string
	Order        *int
	Type         *CatalogueCategoryPropertyType
	Unit         *Optional[models.Codebook]
	GroupUID     *string
}

// CreateCatalogueCategoryPropertyFields — POST /.../group/:gid/property payload.
type CreateCatalogueCategoryPropertyFields struct {
	Name         string
	DefaultValue *string
	ListOfValues []string
	Order        *int
	Type         CatalogueCategoryPropertyType
	Unit         *models.Codebook
}

// PatchCatalogueCategoryPhysicalPropertyFields — PATCH /.../physical-property/:pid payload.
// Same shape as property but physical props have no group and aren't referenced by items.
type PatchCatalogueCategoryPhysicalPropertyFields struct {
	Name         *string
	DefaultValue *Optional[string]
	ListOfValues *[]string
	Order        *int
	Type         *CatalogueCategoryPropertyType
	Unit         *Optional[models.Codebook]
}

// CreateCatalogueCategoryPhysicalPropertyFields — POST /.../physical-property payload.
type CreateCatalogueCategoryPhysicalPropertyFields struct {
	Name         string
	DefaultValue *string
	ListOfValues []string
	Order        *int
	Type         CatalogueCategoryPropertyType
	Unit         *models.Codebook
}
