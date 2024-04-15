package models

type FileLink struct {
	UID  string   `json:"uid"`
	Name string   `json:"name"`
	Url  string   `json:"url"`
	Tags []string `json:"tags"`
}
