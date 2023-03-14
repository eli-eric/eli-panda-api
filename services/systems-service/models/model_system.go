package models

type SystemResponse struct {
	UID         string                 `json:"uid"`
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
	UID           *string `json:"uid"`
	ParentUID     *string `json:"parentUID"`
	Name          string  `json:"name"`
	Description   *string `json:"description"`
	SystemTypeUID *string `json:"systemTypeUID"`
	SystemCode    *string `json:"systemCode"`
	SystemAlias   *string `json:"systemAlias"`
	LocationUID   *string `json:"locationUID"`
	ItemUID       *string `json:"itemUID"`
	OwnerUID      *string `json:"ownerUID"`
	ImportanceUID *string `json:"importanceUID"`
	ZoneUID       *string `json:"zoneUID"`
	Image         *string `json:"image"`
}

type SystemSimpleResponse struct {
	UID  string `json:"uid"`
	Name string `json:"name"`
}
