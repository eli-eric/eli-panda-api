package models

type CatalogueCategory struct {
	UID string `json:"uid,omitempty"`

	Name string `json:"name,omitempty"`

	Code string `json:"code,omitempty"`

	ParentPath string `json:"parentPath,omitempty"`

	Groups []CatalogueCategoryPropertyGroup `json:"groups"`
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

	TypeUID string `json:"typeUID,omitempty"`

	UnitUID string `json:"unitUID,omitempty"`
}
