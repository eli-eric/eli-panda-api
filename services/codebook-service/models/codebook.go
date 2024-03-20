package models

import (
	"panda/apigateway/shared"
)

type Codebook struct {
	UID            string `json:"uid"`
	Name           string `json:"name"`
	AdditionalData string `json:"additionalData,omitempty"`
}

type CodebookTreeItem struct {
	UID      string             `json:"uid"`
	Name     string             `json:"name"`
	Children []CodebookTreeItem `json:"children,omitempty"`
}

type CodebookTreeItemCatalogueCategory struct {
	UID      string                              `json:"uid"`
	Name     string                              `json:"name"`
	Children []CodebookTreeItemCatalogueCategory `json:"has_subcategory,omitempty"`
}

type CodebookType struct {
	Code             string `json:"code"`
	Type             string `json:"type"`
	RoleEdit         string `json:"roleEdit,omitempty"`
	NodeLabel        string `json:"-"`
	FacilityRelation string `json:"-"`
}

type CodebookResponse struct {
	Metadata CodebookType `json:"metadata"`
	Data     []Codebook   `json:"data"`
}

var ZONE_CODEBOOK CodebookType = CodebookType{Code: "ZONE", Type: "SIMPLE"}
var UNIT_CODEBOOK CodebookType = CodebookType{Code: "UNIT", Type: "SIMPLE", NodeLabel: "Unit", RoleEdit: shared.ROLE_CODEBOOKS_ADMIN}
var CATALOGUE_PROPERTY_TYPE_CODEBOOK CodebookType = CodebookType{Code: "CATALOGUE_PROPERTY_TYPE", Type: "SIMPLE"}
var SYSTEM_TYPE_CODEBOOK CodebookType = CodebookType{Code: "SYSTEM_TYPE", Type: "SIMPLE"}
var SYSTEM_IMPORTANCE_CODEBOOK CodebookType = CodebookType{Code: "SYSTEM_IMPORTANCE", Type: "SIMPLE"}
var SYSTEM_CRITICALITY_CLASS_CODEBOOK CodebookType = CodebookType{Code: "SYSTEM_CRITICALITY_CLASS", Type: "SIMPLE"}
var ITEM_USAGE_CODEBOOK CodebookType = CodebookType{Code: "ITEM_USAGE", Type: "SIMPLE"}
var ITEM_CONDITION_STATUS_CODEBOOK CodebookType = CodebookType{Code: "ITEM_CONDITION_STATUS", Type: "SIMPLE"}
var USER_CODEBOOK CodebookType = CodebookType{Code: "USER", Type: "SIMPLE"}
var ORDER_STATUS_CODEBOOK CodebookType = CodebookType{Code: "ORDER_STATUS", Type: "SIMPLE"}
var PROCUREMENTER_CODEBOOK CodebookType = CodebookType{Code: "PROCUREMENTER", Type: "SIMPLE"}
var LOCATION_AUTOCOMPLETE_CODEBOOK CodebookType = CodebookType{Code: "LOCATION", Type: "AUTOCOMPLETE"}
var EMPLOYEE_AUTOCOMPLETE_CODEBOOK CodebookType = CodebookType{Code: "EMPLOYEE", Type: "AUTOCOMPLETE"}
var SYSTEM_AUTOCOMPLETE_CODEBOOK CodebookType = CodebookType{Code: "SYSTEM", Type: "AUTOCOMPLETE"}
var USER_AUTOCOMPLETE_CODEBOOK CodebookType = CodebookType{Code: "USER", Type: "AUTOCOMPLETE"}
var SUPPLIER_AUTOCOMPLETE_CODEBOOK CodebookType = CodebookType{Code: "SUPPLIER", Type: "AUTOCOMPLETE", NodeLabel: "Supplier", RoleEdit: shared.ROLE_CODEBOOKS_ADMIN}
var CATALOGUE_CATEGORY_AUTOCOMPLETE_CODEBOOK CodebookType = CodebookType{Code: "CATALOGUE_CATEGORY", Type: "AUTOCOMPLETE"}
var TEAM_AUTOCOMPLETE_CODEBOOK CodebookType = CodebookType{Code: "TEAM", Type: "AUTOCOMPLETE", NodeLabel: "Team", RoleEdit: shared.ROLE_CODEBOOKS_ADMIN, FacilityRelation: "BELONGS_TO_FACILITY"}
var CONTACT_PERSON_ROLE_CODEBOOK CodebookType = CodebookType{Code: "CONTACT_PERSON_ROLE", Type: "AUTOCOMPLETE", NodeLabel: "ContactPersonRole", RoleEdit: shared.ROLE_CODEBOOKS_ADMIN, FacilityRelation: "BELONGS_TO_FACILITY"}
