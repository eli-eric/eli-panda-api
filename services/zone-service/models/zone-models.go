package models

import (
	codebookModels "panda/apigateway/services/codebook-service/models"
)

type Zone struct {
	UID        string                   `json:"uid" neo4j:"key,uid"`
	Name       string                   `json:"name" neo4j:"prop,name"`
	Code       string                   `json:"code" neo4j:"prop,code"`
	Notes      string                   `json:"notes" neo4j:"prop,notes"`
	ParentZone *codebookModels.Codebook `json:"parentZone,omitempty"`
}

type ZoneCreateRequest struct {
	Name      string  `json:"name"`
	Code      string  `json:"code"`
	Notes     string  `json:"notes"`
	ParentUID *string `json:"parentUid,omitempty"`
}

type ZoneUpdateRequest struct {
	Name      string  `json:"name"`
	Code      string  `json:"code"`
	Notes     *string `json:"notes,omitempty"`
	ParentUID *string `json:"parentUid,omitempty"`
}

type ZoneImportResult struct {
	Created int      `json:"created"`
	Skipped int      `json:"skipped"`
	Errors  []string `json:"errors"`
}
