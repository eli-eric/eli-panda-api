package models

import "panda/apigateway/services/codebook-service/models"

type CatalogueServiceType struct {
	Uid         string          `json:"uid,omitempty"`
	Name        string          `json:"name,omitempty" neo4j:"prop,name"`
	Description string          `json:"description,omitempty" neo4j:"prop,description"`
	Category    models.Codebook `json:"category,omitempty" neo4j:"rel,CatalogueCategory,BELONGS_TO_CATEGORY,uid,cc"`
	Properties  []string        `json:"properties,omitempty" neo4j:"prop,properties"`
}
