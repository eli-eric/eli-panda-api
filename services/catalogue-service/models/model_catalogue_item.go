package models

import "panda/apigateway/services/codebook-service/models"

type CatalogueItem struct {
	Uid string `json:"uid,omitempty"`

	Name string `json:"name,omitempty"`

	CatalogueNumber string `json:"catalogueNumber,omitempty"`

	Description *string `json:"description,omitempty"`

	CategoryPath string `json:"categoryPath,omitempty"`

	Category models.Codebook `json:"category,omitempty"`

	Manufacturer *models.Codebook `json:"manufacturer,omitempty"`

	ManufacturerNumber *string `json:"manufacturerNumber,omitempty"`

	ManufacturerUrl *string `json:"manufacturerUrl,omitempty"`

	Details []CatalogueItemDetail `json:"details,omitempty"`
}

type CatalogueItemDetail struct {
	Property CatalogueCategoryProperty `json:"property,omitempty"`

	PropertyGroup string `json:"propertyGroup,omitempty"`

	Value *string `json:"value,omitempty"`
}
