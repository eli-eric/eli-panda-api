package publicationsservice

import (
	"encoding/json"
	"net/http"
	"panda/apigateway/helpers"
	codebookModels "panda/apigateway/services/codebook-service/models"
	"panda/apigateway/services/publications-service/models"
	"strconv"
	"strings"

	"github.com/neo4j/neo4j-go-driver/v4/neo4j"
)

type PublicationsService struct {
	neo4jDriver      *neo4j.Driver
	wosStarterApiUrl string
	wosStarterApiKey string
}

type IPublicationsService interface {
	GetPublicationByUid(uid string) (models.Publication, error)
	CreatePublication(publication *models.Publication, userUID string) (result models.Publication, err error)
	UpdatePublication(newPublication *models.Publication, userUID string) (result models.Publication, err error)
	DeletePublication(uid string, userUID string) (err error)
	GetPublications(searchText string, page, pageSize int) (result []models.Publication, totalCount int64, err error)
	GetPublicationByDoiFromWOS(doi string) (models.WosAPIResponse, error)
}

func NewPublicationsService(driver *neo4j.Driver, wosSAPIURL, wosSAPIKEY string) IPublicationsService {
	return &PublicationsService{neo4jDriver: driver, wosStarterApiUrl: wosSAPIURL, wosStarterApiKey: wosSAPIKEY}
}

func (svc *PublicationsService) GetPublicationByUid(uid string) (result models.Publication, err error) {

	session, _ := helpers.NewNeo4jSession(*svc.neo4jDriver)

	result.Uid = uid

	err = helpers.GetSingleNode(session, &result)

	if err == nil {
		decodeAuthorsDepartments(&result)
	}

	return result, err
}

func decodeAuthorsDepartments(publication *models.Publication) {

	if len(publication.AuthorsDepartmentsArray) > 0 {
		publication.AuthorsDepartments = make([]models.AuthorsDepartment, 0)
		for i := 0; i < len(publication.AuthorsDepartmentsArray); i++ {
			authorDepartmentString := strings.Split(publication.AuthorsDepartmentsArray[i], "||")
			if len(authorDepartmentString) == 3 {
				uid := authorDepartmentString[0]
				name := authorDepartmentString[1]
				authorsCount, _ := strconv.Atoi(authorDepartmentString[2])

				authorDepartment := models.AuthorsDepartment{Department: codebookModels.Codebook{UID: uid, Name: name}, AuthorsCount: authorsCount}
				publication.AuthorsDepartments = append(publication.AuthorsDepartments, authorDepartment)
			}
		}
	}
}

func encodeAuthorsDepartments(publication *models.Publication) {

	if len(publication.AuthorsDepartments) > 0 {
		publication.AuthorsDepartmentsArray = make([]string, 0)
		for i := 0; i < len(publication.AuthorsDepartments); i++ {
			authorDepartment := publication.AuthorsDepartments[i]
			authorsDepartmentString := authorDepartment.Department.UID + "||" + authorDepartment.Department.Name + "||" + strconv.Itoa(authorDepartment.AuthorsCount)
			publication.AuthorsDepartmentsArray = append(publication.AuthorsDepartmentsArray, authorsDepartmentString)
		}
	}
}

func (svc *PublicationsService) CreatePublication(publication *models.Publication, userUID string) (result models.Publication, err error) {

	session, _ := helpers.NewNeo4jSession(*svc.neo4jDriver)

	encodeAuthorsDepartments(publication)

	updateQuery := helpers.DatabaseQuery{}
	updateQuery.Parameters = make(map[string]interface{})
	updateQuery.Query = `MERGE (n:Publication {uid: $uid}) SET n.updatedAt = datetime() WITH n `
	updateQuery.Parameters["uid"] = publication.Uid

	helpers.AutoResolveObjectToUpdateQuery(&updateQuery, *publication, models.Publication{}, "n")

	updateQuery.Query += ` RETURN n.uid as uid `
	updateQuery.ReturnAlias = "uid"

	historyLog := helpers.HistoryLogQuery(publication.Uid, "CREATE", userUID)

	err = helpers.WriteNeo4jAndReturnNothingMultipleQueries(session,
		updateQuery,
		historyLog)

	return *publication, err
}

func (svc *PublicationsService) UpdatePublication(publication *models.Publication, userUID string) (result models.Publication, err error) {

	session, _ := helpers.NewNeo4jSession(*svc.neo4jDriver)

	encodeAuthorsDepartments(publication)

	oldPublication, err := svc.GetPublicationByUid(publication.Uid)

	if err != nil {
		return result, err
	}

	updateQuery := helpers.DatabaseQuery{}
	updateQuery.Parameters = make(map[string]interface{})
	updateQuery.Query = `MATCH (n:Publication {uid: $uid}) SET n.updatedAt = datetime() WITH n  `
	updateQuery.Parameters["uid"] = publication.Uid

	helpers.AutoResolveObjectToUpdateQuery(&updateQuery, *publication, oldPublication, "n")

	updateQuery.Query += ` RETURN n.uid as uid `
	updateQuery.ReturnAlias = "uid"

	historyLog := helpers.HistoryLogQuery(publication.Uid, "UPDATE", userUID)

	err = helpers.WriteNeo4jAndReturnNothingMultipleQueries(session,
		updateQuery,
		historyLog)

	return result, err
}

func (svc *PublicationsService) DeletePublication(uid string, userUID string) (err error) {

	session, _ := helpers.NewNeo4jSession(*svc.neo4jDriver)

	err = helpers.WriteNeo4jAndReturnNothingMultipleQueries(session,
		helpers.SoftDeleteNodeQuery(uid),
		helpers.HistoryLogQuery(uid, "DELETE", userUID))

	return err

}

func (svc *PublicationsService) GetPublications(searchText string, page, pageSize int) (result []models.Publication, totalCount int64, err error) {

	session, _ := helpers.NewNeo4jSession(*svc.neo4jDriver)

	skip := pageSize
	limit := pageSize

	if skip > 0 {
		skip = (page - 1) * pageSize
	} else {
		skip = 0
	}

	if limit <= 0 {
		limit = 10
	}

	result, totalCount, err = helpers.GetMultipleNodes[models.Publication](session, skip, limit, searchText)

	helpers.ProcessArrayResult(&result, err)

	for i := 0; i < len(result); i++ {
		decodeAuthorsDepartments(&result[i])
	}

	return result, totalCount, err
}

func (svc *PublicationsService) GetPublicationByDoiFromWOS(doi string) (result models.WosAPIResponse, err error) {

	// exmaple get url /documents?db=WOS&q=DO=10.1103/PhysRevResearch.6.013126
	// get from wos rest api

	contentType := "application/json"
	query := "/documents?db=WOS&q=DO=" + doi

	request, err := http.NewRequest("GET", svc.wosStarterApiUrl+query, nil)

	if err != nil {
		return result, err
	}

	// addd header X-ApiKey
	// wos.starter with institution key
	request.Header.Add("X-ApiKey", svc.wosStarterApiKey)
	request.Header.Add("Content-Type", contentType)

	client := &http.Client{}
	response, err := client.Do(request)

	if err != nil {
		return result, err
	} else {
		err := json.NewDecoder(response.Body).Decode(&result)

		if err != nil {
			return result, err
		}

		defer response.Body.Close()

		return result, nil
	}
}
