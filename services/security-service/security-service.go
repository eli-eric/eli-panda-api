package securityService

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"panda/apigateway/config"
	"panda/apigateway/helpers"
	codebookModels "panda/apigateway/services/codebook-service/models"
	"panda/apigateway/services/security-service/models"
	"time"

	"github.com/rs/zerolog/log"

	"github.com/golang-jwt/jwt"
	"github.com/neo4j/neo4j-go-driver/v4/neo4j"
	"golang.org/x/crypto/bcrypt"
)

type SecurityService struct {
	neo4jDriver *neo4j.Driver
	jwtSecret   string
}

type ISecurityService interface {
	AuthenticateByUsernameAndPassword(username string, password string) (authUser models.UserAuthInfo, err error)
	GetUsersCodebook(facilityCode string) (result []codebookModels.Codebook, err error)
	GetUsersAutocompleteCodebook(searchText string, limit int, facilityCode string) (result []codebookModels.Codebook, err error)
	ChangeUserPassword(userName string, userUID string, passwords *models.ChangePasswordRequest) (err error)
	GetEmployeesAutocompleteCodebook(searchText string, limit int, facilityCode string, filter *[]helpers.Filter, isAdmin bool) (result []codebookModels.Codebook, err error)
	GetProcurementersCodebook(facilityCode string) (result []codebookModels.Codebook, err error)
	GetTeamsAutocompleteCodebook(searchText string, limit int, facilityCode string) (result []codebookModels.Codebook, err error)
	GetContactPersonRolesAutocompleteCodebook(searchText string, limit int, facilityCode string) (result []codebookModels.Codebook, err error)
	GetUserByAzureIdToken(azureIdToken string) (authUser models.UserAuthInfo, err error)
}

// Create new security service instance
func NewSecurityService(settings *config.Config, driver *neo4j.Driver) ISecurityService {

	return &SecurityService{neo4jDriver: driver, jwtSecret: settings.JwtSecret}
}

func (svc *SecurityService) AuthenticateByUsernameAndPassword(username string, password string) (authUser models.UserAuthInfo, err error) {

	// Open a new Session

	session, _ := helpers.NewNeo4jSession(*svc.neo4jDriver)

	//the user has to be enabled and has at least one role
	authUser, err = helpers.GetNeo4jSingleRecordAndMapToStruct[models.UserAuthInfo](session, UserWithRolesAndFailityQuery(username))

	//if there is a user in DB lets check the password
	if err == nil {

		if !authUser.IsEnabled {
			return authUser, errors.New("Unauthorized")
		}

		verifErr := bcrypt.CompareHashAndPassword([]byte(authUser.PasswordHash), []byte(password))
		//empty passwordHash -> omitempty json -> not sent to client
		authUser.PasswordHash = ""
		// Throws unauthorized error if there is verifErr
		if verifErr == nil {
			// Set custom claims
			claims := &models.JwtCustomClaims{
				Roles:        authUser.Roles,
				FacilityName: authUser.Facility,
				FacilityCode: authUser.FacilityCode,
				StandardClaims: jwt.StandardClaims{
					ExpiresAt: time.Now().Add(time.Hour * 876000).Unix(),
					Subject:   authUser.Uid,
					Id:        username,
				},
			}

			// Create token with claims
			newToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

			// Generate encoded token and send it as response.
			token, err := newToken.SignedString([]byte(svc.jwtSecret))
			if err == nil {
				authUser.AccessToken = token
			}

			return authUser, err
		}
	} else {
		log.Error().Msg(err.Error())
	}

	return authUser, errors.New("Unauthorized")
}

func (svc *SecurityService) GetUsersCodebook(facilityCode string) (result []codebookModels.Codebook, err error) {
	session, _ := helpers.NewNeo4jSession(*svc.neo4jDriver)

	query := GetUsersCodebookQuery(facilityCode)
	result, err = helpers.GetNeo4jArrayOfNodes[codebookModels.Codebook](session, query)

	return result, err
}

func (svc *SecurityService) GetUsersAutocompleteCodebook(searchText string, limit int, facilityCode string) (result []codebookModels.Codebook, err error) {
	session, _ := helpers.NewNeo4jSession(*svc.neo4jDriver)

	query := GetUsersAutocompleteCodebookQuery(searchText, limit, facilityCode)
	result, err = helpers.GetNeo4jArrayOfNodes[codebookModels.Codebook](session, query)

	return result, err
}

func (svc *SecurityService) GetTeamsAutocompleteCodebook(searchText string, limit int, facilityCode string) (result []codebookModels.Codebook, err error) {
	session, _ := helpers.NewNeo4jSession(*svc.neo4jDriver)

	query := GetTeamsAutocompleteCodebookQuery(searchText, limit, facilityCode)
	result, err = helpers.GetNeo4jArrayOfNodes[codebookModels.Codebook](session, query)

	return result, err
}

func (svc *SecurityService) GetContactPersonRolesAutocompleteCodebook(searchText string, limit int, facilityCode string) (result []codebookModels.Codebook, err error) {
	session, _ := helpers.NewNeo4jSession(*svc.neo4jDriver)

	query := GetContactPersonRolesAutocompleteCodebookQuery(searchText, limit, facilityCode)
	result, err = helpers.GetNeo4jArrayOfNodes[codebookModels.Codebook](session, query)

	return result, err
}

func (svc *SecurityService) ChangeUserPassword(userName string, userUID string, passwords *models.ChangePasswordRequest) (err error) {

	//check the old password first
	_, err = svc.AuthenticateByUsernameAndPassword(userName, passwords.OldPassword)

	//if it is ok we can change set the new one
	if err == nil {
		//create hash from the new password string
		newHashBytes, err := bcrypt.GenerateFromPassword([]byte(passwords.NewPassword), 12)
		if err == nil {
			newPasswordHash := string(newHashBytes)
			session, _ := helpers.NewNeo4jSession(*svc.neo4jDriver)

			_, err := helpers.WriteNeo4jAndReturnSingleValue[string](session, ChangeUserPasswordQuery(userUID, newPasswordHash))
			return err
		}
	}

	return err
}

func (svc *SecurityService) GetEmployeesAutocompleteCodebook(searchText string, limit int, facilityCode string, filter *[]helpers.Filter, isAdmin bool) (result []codebookModels.Codebook, err error) {
	session, _ := helpers.NewNeo4jSession(*svc.neo4jDriver)

	getAllEmployees := false

	if isAdmin {
		if filter != nil {
			for _, f := range *filter {
				if f.Key == "all" && f.Value == "true" {
					getAllEmployees = true
					break
				}
			}
		}
	}

	query := GetEmployeesAutocompleteCodebookQuery(searchText, limit, facilityCode, getAllEmployees)
	result, err = helpers.GetNeo4jArrayOfNodes[codebookModels.Codebook](session, query)

	return result, err
}

// get all employess with flag isProcurementer = true
func (svc *SecurityService) GetProcurementersCodebook(facilityCode string) (result []codebookModels.Codebook, err error) {
	session, _ := helpers.NewNeo4jSession(*svc.neo4jDriver)

	query := GetEmployeesAutocompleteCodebookQuery("", 100, facilityCode, false, EMPLOYEE_FLAG_PROCUREMENTER)
	result, err = helpers.GetNeo4jArrayOfNodes[codebookModels.Codebook](session, query)

	return result, err
}

func (svc *SecurityService) GetUserByAzureIdToken(azureIdToken string) (authUser models.UserAuthInfo, err error) {

	// Open a new Session
	session, _ := helpers.NewNeo4jSession(*svc.neo4jDriver)

	jwksURL := "https://login.microsoftonline.com/80f0138f-bb28-48e8-b6e8-2b1118fca6d8/discovery/v2.0/keys"

	jwks, err := getJWKs(jwksURL)
	if err != nil {
		fmt.Println("Error fetching JWKs:", err)
		return
	}

	token, err := parseAzureADJWT(azureIdToken, jwks)

	if err != nil {
		fmt.Println("Error parsing token:", err)
		return authUser, err
	}

	claims, ok := token.Claims.(*models.JWTAzureAdClaims)
	if !ok {
		fmt.Println("Error parsing claims")
		return authUser, errors.New("Unauthorized")
	}

	//the user has to be enabled and has at least one role
	authUser, err = helpers.GetNeo4jSingleRecordAndMapToStruct[models.UserAuthInfo](session, UserWithRolesAndFailityQuery(claims.Email))

	//if there is a user in DB lets check the password
	if err == nil {

		if !authUser.IsEnabled {
			return authUser, errors.New("Unauthorized")
		}

		authUser.PasswordHash = ""

		// Set custom claims
		claims := &models.JwtCustomClaims{
			Roles:        authUser.Roles,
			FacilityName: authUser.Facility,
			FacilityCode: authUser.FacilityCode,
			StandardClaims: jwt.StandardClaims{
				ExpiresAt: time.Now().Add(time.Hour * 876000).Unix(),
				Subject:   authUser.Uid,
				Id:        authUser.Username,
			},
		}

		// Create token with claims
		newToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

		// Generate encoded token and send it as response.
		token, err := newToken.SignedString([]byte(svc.jwtSecret))
		if err == nil {
			authUser.AccessToken = token
		}

		return authUser, err
	}

	return authUser, err
}

// JWK represents a JSON Web Key
type JWK struct {
	Keys []struct {
		Kty string   `json:"kty"`
		Use string   `json:"use"`
		Kid string   `json:"kid"`
		X5t string   `json:"x5t"`
		N   string   `json:"n"`
		E   string   `json:"e"`
		X5c []string `json:"x5c"`
	} `json:"keys"`
}

// Fetches the JWKs (JSON Web Keys) from Azure AD
func getJWKs(url string) (*JWK, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var jwks JWK
	err = json.Unmarshal(body, &jwks)
	if err != nil {
		return nil, err
	}

	return &jwks, nil
}

// Parses and validates the JWT token
func parseAzureADJWT(tokenString string, jwks *JWK) (*jwt.Token, error) {

	token, err := jwt.ParseWithClaims(tokenString, &models.JWTAzureAdClaims{}, func(token *jwt.Token) (interface{}, error) {
		// Validate the algorithm
		if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		// Get the kid from the token header
		kid := token.Header["kid"].(string)

		// Find the corresponding public key
		var cert string
		for _, key := range jwks.Keys {
			if key.Kid == kid {
				cert = "-----BEGIN CERTIFICATE-----\n" + key.X5c[0] + "\n-----END CERTIFICATE-----"
				break
			}
		}

		// Parse the public key
		key, err := jwt.ParseRSAPublicKeyFromPEM([]byte(cert))
		if err != nil {
			return nil, err
		}

		return key, nil
	})

	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, fmt.Errorf("invalid token")
	}

	return token, nil
}

const EMPLOYEE_FLAG_PROCUREMENTER string = "isProcurementer" //it means that this employee can be selected as procurementer
