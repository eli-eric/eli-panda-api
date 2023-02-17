package models

type CatalogueCategory struct {
	Uid string `json:"uid,omitempty"`

	Name string `json:"name,omitempty"`

	Code string `json:"code,omitempty"`

	ParentPath string `json:"parentPath,omitempty"`
}
