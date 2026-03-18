package models

type Zone struct {
	UID        string  `json:"uid" neo4j:"key,uid"`
	Name       string  `json:"name" neo4j:"prop,name"`
	Code       string  `json:"code" neo4j:"prop,code"`
	ParentUID  *string `json:"parentUid,omitempty"`
	ParentName *string `json:"parentName,omitempty"`
	ParentCode *string `json:"parentCode,omitempty"`
}

type ZoneCreateRequest struct {
	Name      string  `json:"name"`
	Code      string  `json:"code"`
	ParentUID *string `json:"parentUid,omitempty"`
}

type ZoneUpdateRequest struct {
	Name      string  `json:"name"`
	Code      string  `json:"code"`
	ParentUID *string `json:"parentUid,omitempty"`
}

type ZoneImportResult struct {
	Created int      `json:"created"`
	Skipped int      `json:"skipped"`
	Errors  []string `json:"errors"`
}
