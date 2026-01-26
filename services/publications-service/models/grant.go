package models

import (
	codebookModels "panda/apigateway/services/codebook-service/models"
	"time"
)

type Grant struct {
	Uid        string                   `json:"uid" neo4j:"key,uid"`
	Code       string                   `json:"code" neo4j:"prop,code"`
	Name       string                   `json:"name" neo4j:"prop,name"`
	GrantGroup *codebookModels.Codebook `json:"grantGroup" neo4j:"rel,GrantGroup,BELONGS_TO_GROUP,uid,grantGroup"`
	UpdatedAt  *time.Time               `json:"updatedAt" neo4j:"prop,updatedAt"`
}
