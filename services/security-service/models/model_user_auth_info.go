package models

type UserAuthInfo struct {
	Uid string `json:"uid,omitempty"`

	Username string `json:"username,omitempty"`

	Email string `json:"email,omitempty"`

	LastName string `json:"lastName,omitempty"`

	FirstName string `json:"firstName,omitempty"`

	Facility string `json:"facility,omitempty"`

	FacilityCode string `json:"facilityCode,omitempty"`

	IsEnabled bool `json:"isEnabled,omitempty"`

	AccessToken string `json:"accessToken,omitempty"`

	Roles []string `json:"roles,omitempty"`

	PasswordHash string `json:"passwordHash,omitempty"`
}

type Facility struct {
	Code string `json:"code,omitempty"`

	Name string `json:"name,omitempty"`
}
