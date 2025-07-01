package models

import (
	"panda/apigateway/services/codebook-service/models"
)

type LayoutRoomCard struct {
	UID         string               `json:"uid" neo4j:"key,uid"`
	Name        string               `json:"name" neo4j:"prop,name"`
	Location    *models.Codebook     `json:"location"`
	PurityClass string               `json:"purityClass" neo4j:"prop,purityClass"`
	Status      string               `json:"status" neo4j:"prop,status"`
	StatusColor string               `json:"statusColor"`
}