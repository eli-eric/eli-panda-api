package models

type UserStatusResponse struct {
	UserUID   string `json:"userUID"`
	IsEnabled bool   `json:"isEnabled"`
}
