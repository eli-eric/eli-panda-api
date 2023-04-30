package models

type Codebook struct {
	UID            string `json:"uid"`
	Name           string `json:"name"`
	AdditionalData string `json:"additionalData,omitempty"`
}
