package models

import (
	"panda/apigateway/services/codebook-service/models"
	"time"
)

type CatalogueItem struct {
	Uid string `json:"uid,omitempty" neo4j:"ignore"`

	Name string `json:"name,omitempty" neo4j:"prop,name"`

	CatalogueNumber string `json:"catalogueNumber,omitempty" neo4j:"prop,catalogueNumber"`

	Description *string `json:"description,omitempty" neo4j:"prop,description"`

	CategoryUid string `json:"categoryUID,omitempty" neo4j:"ignore"`

	Category models.Codebook `json:"category,omitempty" neo4j:"rel,CatalogueCategory,BELONGS_TO_CATEGORY,uid,cc"`

	Supplier *models.Codebook `json:"supplier,omitempty" neo4j:"rel,Supplier,HAS_SUPPLIER,uid,supp"`

	ManufacturerNumber *string `json:"manufacturerNumber,omitempty" neo4j:"prop,manufacturerNumber"`

	ManufacturerUrl *string `json:"manufacturerUrl,omitempty" neo4j:"prop,manufacturerUrl"`

	Details []CatalogueItemDetail `json:"details,omitempty" neo4j:"ignore"`

	MiniImageUrl *[]string `json:"miniImageUrl" neo4j:"ignore"`

	LastUpdateTime time.Time `json:"lastUpdateTime"`
}

type CatalogueItemDetail struct {
	Property CatalogueCategoryProperty `json:"property,omitempty"`

	PropertyGroup string `json:"propertyGroup,omitempty"`

	Value any `json:"value,omitempty"`
}

type CatalogueItemSimple struct {
	Uid string `json:"uid,omitempty"`

	Name string `json:"name,omitempty"`

	CatalogueNumber string `json:"catalogueNumber,omitempty"`

	Description string `json:"description,omitempty"`

	Category models.Codebook `json:"category,omitempty"`

	Supplier *models.Codebook `json:"supplier,omitempty"`

	ManufacturerNumber string `json:"manufacturerNumber,omitempty"`

	ManufacturerUrl string `json:"manufacturerUrl,omitempty"`

	Details []CatalogueItemDetail `json:"details,omitempty"`

	MiniImageUrl *[]string `json:"miniImageUrl" neo4j:"ignore"`

	LastUpdateTime *time.Time `json:"lastUpdateTime"`

	LastUpdateBy *string `json:"lastUpdateBy"`
}

type CatalogueItemDetailSimple struct {
	PropertyName string `json:"propertyName,omitempty"`

	PropertyUnit string `json:"propertyUnit"`

	PropertyType string `json:"propertyType,omitempty"`

	PropertyGroup string `json:"propertyGroup,omitempty"`

	Value string `json:"value,omitempty"`
}
