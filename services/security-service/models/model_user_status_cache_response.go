package models

import "time"

type UserStatusCacheItem struct {
	UserUID   string    `json:"userUID"`
	IsEnabled bool      `json:"isEnabled"`
	ExpiresAt time.Time `json:"expiresAt"`
}
