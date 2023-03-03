package models

type SystemResponse struct {
	UID              string                 `json:"uid"`
	Name             string                 `json:"name"`
	ParentPath       []SystemSimpleResponse `json:"parentPath"`
	Description      string                 `json:"description"`
	SystemType       string                 `json:"systemType"`
	SystemCode       string                 `json:"systemCode"`
	SystemAlias      string                 `json:"systemAlias"`
	Location         string                 `json:"location"`
	ItemUID          string                 `json:"itemUID"`
	Owner            string                 `json:"owner"`
	Importance       string                 `json:"importance"`
	Zone             string                 `json:"zone"`
	SubZone          string                 `json:"subZone"`
	CriticalityClass string                 `json:"criticalityClass"`
}

type SystemForm struct {
	UID                 string `json:"uid"`
	ParentUID           string `json:"parentUID"`
	Name                string `json:"name"`
	Description         string `json:"description"`
	SystemTypeUID       string `json:"systemTypeUID"`
	SystemCode          string `json:"systemCode"`
	SystemAlias         string `json:"systemAlias"`
	LocationUID         string `json:"locationUID"`
	ItemUID             string `json:"itemUID"`
	OwnerUID            string `json:"ownerUID"`
	ImportanceUID       string `json:"importanceUID"`
	ZoneUID             string `json:"zoneUID"`
	SubZoneUID          string `json:"subZoneUID"`
	CriticalityClassUID string `json:"criticalityClassUID"`
	Image               string `json:"image"`
}

type SystemSimpleResponse struct {
	UID  string `json:"uid"`
	Name string `json:"name"`
}
