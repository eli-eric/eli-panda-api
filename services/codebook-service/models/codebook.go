package models

import (
	"panda/apigateway/shared"
)

type Codebook struct {
	UID            string `json:"uid"`
	Name           string `json:"name"`
	Code           string `json:"code,omitempty"`
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

var ZONE_CODEBOOK CodebookType = CodebookType{Code: "ZONE", Type: "AUTOCOMPLETE", NodeLabel: "Zone"}
var UNIT_CODEBOOK CodebookType = CodebookType{Code: "UNIT", Type: "SIMPLE", NodeLabel: "Unit", RoleEdit: shared.ROLE_CODEBOOKS_ADMIN}
var CATALOGUE_PROPERTY_TYPE_CODEBOOK CodebookType = CodebookType{Code: "CATALOGUE_PROPERTY_TYPE", Type: "SIMPLE", NodeLabel: "CataloguePropertyType"}
var SYSTEM_TYPE_CODEBOOK CodebookType = CodebookType{Code: "SYSTEM_TYPE", Type: "SIMPLE"}
var SYSTEM_IMPORTANCE_CODEBOOK CodebookType = CodebookType{Code: "SYSTEM_IMPORTANCE", Type: "SIMPLE"}
var SYSTEM_CRITICALITY_CLASS_CODEBOOK CodebookType = CodebookType{Code: "SYSTEM_CRITICALITY_CLASS", Type: "SIMPLE"}
var ITEM_USAGE_CODEBOOK CodebookType = CodebookType{Code: "ITEM_USAGE", Type: "SIMPLE", NodeLabel: "ItemUsage"}
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
var CONTACT_PERSON_ROLE_CODEBOOK CodebookType = CodebookType{Code: "CONTACT_PERSON_ROLE", Type: "AUTOCOMPLETE", NodeLabel: "ContactPersonRole", RoleEdit: shared.ROLE_ROOM_CARDS_EDIT, FacilityRelation: "BELONGS_TO_FACILITY"}
var SYSTEM_ATTRIBUTE_CODEBOOK CodebookType = CodebookType{Code: "SYSTEM_ATTRIBUTE", Type: "SIMPLE", NodeLabel: "SystemAttribute", RoleEdit: shared.SYSTEM_ATTRIBUTE_EDIT, FacilityRelation: "BELONGS_TO_FACILITY"}

// publications related codebooks
// UserCall, UserExperiment, PublicationCategory, OpenAccessType, Language, PublicationSupport, State
var USER_CALL_CODEBOOK CodebookType = CodebookType{Code: "USER_CALL", Type: "SIMPLE", RoleEdit: shared.ROLE_CODEBOOKS_ADMIN, NodeLabel: "UserCall"}
var USER_EXPERIMENT_CODEBOOK CodebookType = CodebookType{Code: "USER_EXPERIMENT", Type: "SIMPLE", RoleEdit: shared.ROLE_CODEBOOKS_ADMIN, NodeLabel: "UserExperiment"}
var PUBLICATION_CATEGORY_CODEBOOK CodebookType = CodebookType{Code: "PUBLICATION_CATEGORY", Type: "SIMPLE", RoleEdit: shared.ROLE_CODEBOOKS_ADMIN, NodeLabel: "PublicationCategory"}
var OPEN_ACCESS_TYPE_CODEBOOK CodebookType = CodebookType{Code: "OPEN_ACCESS_TYPE", Type: "SIMPLE", RoleEdit: shared.ROLE_CODEBOOKS_ADMIN, NodeLabel: "OpenAccessType"}
var LANGUAGE_CODEBOOK CodebookType = CodebookType{Code: "LANGUAGE", Type: "SIMPLE", RoleEdit: shared.ROLE_CODEBOOKS_ADMIN, NodeLabel: "Language"}
var PUBLICATION_SUPPORT_CODEBOOK CodebookType = CodebookType{Code: "PUBLICATION_SUPPORT", Type: "SIMPLE", RoleEdit: shared.ROLE_CODEBOOKS_ADMIN, NodeLabel: "PublicationSupport"}
var STATE_CODEBOOK CodebookType = CodebookType{Code: "STATE", Type: "SIMPLE", RoleEdit: shared.ROLE_CODEBOOKS_ADMIN, NodeLabel: "State"}
