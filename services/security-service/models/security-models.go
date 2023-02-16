package models

import (
	"github.com/golang-jwt/jwt"
)

// jwtCustomClaims are custom claims extending default ones.
// See https://github.com/golang-jwt/jwt for more examples
type JwtCustomClaims struct {
	Roles        []string `json:"roles"`
	FacilityName string   `json:"facilityName"`
	FacilityCode string   `json:"facilityCode"`
	jwt.StandardClaims
}
