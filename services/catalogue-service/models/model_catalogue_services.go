package models

import "panda/apigateway/services/codebook-service/models"

type CatalogueServiceType struct {
	Uid         string           `json:"uid" neo4j:"key,uid"`
	Name        string           `json:"name" neo4j:"prop,name"`
	Description *string          `json:"description" neo4j:"prop,description"`
	Category    *models.Codebook `json:"category" neo4j:"rel,CatalogueCategory,BELONGS_TO_CATEGORY,uid,category"`
	Properties  []string         `json:"properties" neo4j:"prop,properties"`
}
