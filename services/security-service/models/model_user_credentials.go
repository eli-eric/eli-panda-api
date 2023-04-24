package models

type UserCredentials struct {
	Username string `json:"username,omitempty"`

	Password string `json:"password,omitempty"`
}

type ChangePasswordRequest struct {
	OldPassword string `json:"oldPassword"`
	NewPassword string `json:"newPassword"`
}
