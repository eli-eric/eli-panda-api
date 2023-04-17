package models

type SystemResponse struct {
	UID         string                 `json:"uid" neo4j:"uid"`
	Name        string                 `json:"name"`
	ParentPath  []SystemSimpleResponse `json:"parentPath"`
	ParentUID   *string                `json:"parentUID,omitempty"`
	Description *string                `json:"description,omitempty"`
	SystemType  *SystemSimpleResponse  `json:"systemType,omitempty"`
	SystemCode  *string                `json:"systemCode,omitempty"`
	SystemAlias *string                `json:"systemAlias,omitempty"`
	Location    *SystemSimpleResponse  `json:"location,omitempty"`
	ItemUID     *string                `json:"itemUID,omitempty"`
	Owner       *SystemSimpleResponse  `json:"owner,omitempty"`
	Importance  *SystemSimpleResponse  `json:"importance,omitempty"`
	Zone        *SystemSimpleResponse  `json:"zone,omitempty"`
}

type SystemForm struct {
	UID           *string `json:"uid" neo4j:"ignore"`
	ParentUID     *string `json:"parentUID" neo4j:"ignore"`
	Name          string  `json:"name" neo4j:"prop,name"`
	Description   *string `json:"description" neo4j:"prop,description"`
	SystemTypeUID *string `json:"systemTypeUID" neo4j:"rel,SystemType,HAS_SYSTEM_TYPE,uid,st"`
	SystemCode    *string `json:"systemCode" neo4j:"prop,systemCode"`
	SystemAlias   *string `json:"systemAlias" neo4j:"prop,systemAlias"`
	LocationUID   *string `json:"locationUID" neo4j:"rel,Location,HAS_LOCATION,code,l"`
	ItemUID       *string `json:"itemUID" neo4j:"ignore"`
	OwnerUID      *string `json:"ownerUID" neo4j:"rel,User,HAS_OWNER,uid,own"`
	ImportanceUID *string `json:"importanceUID" neo4j:"rel,SystemImportance,HAS_IMPORTANCE,uid,imp"`
	ZoneUID       *string `json:"zoneUID" neo4j:"rel,Zone,HAS_ZONE,uid,z"`
	Image         *string `json:"image" neo4j:"ignore"`
}

type SystemSimpleResponse struct {
	UID  string `json:"uid"`
	Name string `json:"name"`
}
