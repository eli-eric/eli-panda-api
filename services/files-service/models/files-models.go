package models

type FileLink struct {
	UID  string   `json:"uid,omitempty"`
	Name string   `json:"name,omitempty"`
	Url  *string  `json:"url,omitempty"`
	Tags []string `json:"tags,omitempty"`
}

type MiniImageLinks struct {
	UID  string    `json:"uid,omitempty"`
	Name string    `json:"name,omitempty"`
	Url  *[]string `json:"url,omitempty"`
	Tags []string  `json:"tags,omitempty"`
}
