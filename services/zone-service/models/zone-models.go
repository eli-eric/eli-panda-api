package models

import (
	codebookModels "panda/apigateway/services/codebook-service/models"
)

type Zone struct {
	UID                 string                   `json:"uid" neo4j:"key,uid"`
	Name                string                   `json:"name" neo4j:"prop,name"`
	Code                string                   `json:"code" neo4j:"prop,code"`
	Notes               string                   `json:"notes" neo4j:"prop,notes"`
	ParentZone          *codebookModels.Codebook `json:"parentZone,omitempty"`
	DefaultParentSystem *codebookModels.Codebook `json:"defaultParentSystem,omitempty"`
}

type ZoneCreateRequest struct {
	Name      string  `json:"name"`
	Code      string  `json:"code"`
	Notes     string  `json:"notes"`
	ParentUID *string `json:"parentUid,omitempty"`
	// DefaultParentSystemUID is the uid of the System new system codes are created under.
	// nil = not set, non-empty = set.
	DefaultParentSystemUID *string `json:"defaultParentSystemUid,omitempty"`
}

type ZoneUpdateRequest struct {
	Name      string  `json:"name"`
	Code      string  `json:"code"`
	Notes     *string `json:"notes,omitempty"`
	ParentUID *string `json:"parentUid,omitempty"`
	// DefaultParentSystemUID is tri-state like ParentUID:
	// nil = preserve current value, "" = detach, non-empty = set.
	DefaultParentSystemUID *string `json:"defaultParentSystemUid,omitempty"`
}

type ZoneImportResult struct {
	Created int      `json:"created"`
	Skipped int      `json:"skipped"`
	Errors  []string `json:"errors"`
}
