package securityService

import (
	"errors"
	"log"
	"panda/apigateway/config"
	"panda/apigateway/helpers"
	codebookModels "panda/apigateway/services/codebook-service/models"
	"panda/apigateway/services/security-service/models"
	"time"

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
	GetEmployeesAutocompleteCodebook(searchText string, limit int, facilityCode string) (result []codebookModels.Codebook, err error)
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
		log.Println(err)
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

func (svc *SecurityService) GetEmployeesAutocompleteCodebook(searchText string, limit int, facilityCode string) (result []codebookModels.Codebook, err error) {
	session, _ := helpers.NewNeo4jSession(*svc.neo4jDriver)

	query := GetEmployeesAutocompleteCodebookQuery(searchText, limit, facilityCode)
	result, err = helpers.GetNeo4jArrayOfNodes[codebookModels.Codebook](session, query)

	return result, err
}
