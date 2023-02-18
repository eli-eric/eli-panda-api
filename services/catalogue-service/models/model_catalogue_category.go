package models

type CatalogueItem struct {
	Uid string `json:"uid,omitempty"`

	Name string `json:"name,omitempty"`

	Description string `json:"description,omitempty"`

	CategoryPath string `json:"categoryPath,omitempty"`

	CategoryName string `json:"categoryName,omitempty"`

	Manufacturer string `json:"manufacturer,omitempty"`

	ManufacturerNumber string `json:"manufacturerNumber,omitempty"`

	ManufacturerUrl string `json:"manufacturerUrl,omitempty"`

	Details []CatalogueItemDetail `json:"details,omitempty"`
}

type CatalogueItemDetail struct {
	PropertyName string `json:"propertyName,omitempty"`

	PropertyUnit string `json:"propertyUnit"`

	PropertyType string `json:"propertyType,omitempty"`

	PropertyGroup string `json:"propertyGroup,omitempty"`

	Value string `json:"value,omitempty"`
}
