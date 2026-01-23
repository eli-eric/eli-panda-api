package publicationsservice

import (
	"encoding/json"
	"fmt"
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
	GetPublications(searchText string, page, pageSize int, sorting *[]helpers.Sorting, filtering *[]helpers.ColumnFilter) (result []models.Publication, totalCount int64, err error)
	GetPublicationByDoiFromWOS(doi string) (models.WosAPIResponse, error)
	// Researcher methods
	GetResearchers(searchText string, page, pageSize int, sorting *[]helpers.Sorting) (result []models.Researcher, totalCount int64, err error)
	GetResearcherByUid(uid string) (models.Researcher, error)
	CreateResearcher(researcher *models.Researcher, userUID string) (result models.Researcher, err error)
	CreateResearchers(researchers []models.Researcher, userUID string) (result []models.Researcher, err error)
	UpdateResearcher(researcher *models.Researcher, userUID string) (result models.Researcher, err error)
	DeleteResearcher(uid string, userUID string) (err error)
	// Grant methods
	GetGrants(searchText string, page, pageSize int, facilityCode string) (result []models.Grant, totalCount int64, err error)
	GetGrantByUid(uid string) (models.Grant, error)
	CreateGrant(grant *models.Grant, userUID string, facilityCode string) (result models.Grant, err error)
	UpdateGrant(grant *models.Grant, userUID string) (result models.Grant, err error)
	DeleteGrant(uid string, userUID string) (err error)
	// Codebook autocomplete methods
	GetExperimentalSystemsAutocomplete(searchText string, limit int, facilityCode string) ([]codebookModels.Codebook, error)
	GetUserExperimentsAutocomplete(searchText string, limit int, facilityCode string) ([]codebookModels.Codebook, error)
	GetCountriesAutocomplete(searchText string, limit int) ([]codebookModels.Codebook, error)
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
		// Fetch connected researchers
		result.EliResearchers, _ = svc.getPublicationResearchers(uid)
		// Fetch connected grants
		result.Grants, _ = svc.getPublicationGrants(uid)
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

	if err != nil {
		return result, err
	}

	// Connect researchers to the publication
	if len(publication.EliResearchers) > 0 {
		err = svc.connectPublicationResearchers(publication.Uid, publication.EliResearchers)
		if err != nil {
			return result, err
		}
	}

	// Connect grants to the publication
	if len(publication.Grants) > 0 {
		err = svc.connectPublicationGrants(publication.Uid, publication.Grants)
		if err != nil {
			return result, err
		}
	}

	return *publication, nil
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

	if err != nil {
		return result, err
	}

	// Update researchers (diff-based: disconnect removed, connect new)
	err = svc.updatePublicationResearchers(publication.Uid, publication.EliResearchers)
	if err != nil {
		return result, err
	}

	// Update grants (diff-based: disconnect removed, connect new)
	err = svc.updatePublicationGrants(publication.Uid, publication.Grants)
	if err != nil {
		return result, err
	}

	return *publication, nil
}

func (svc *PublicationsService) DeletePublication(uid string, userUID string) (err error) {

	session, _ := helpers.NewNeo4jSession(*svc.neo4jDriver)

	err = helpers.WriteNeo4jAndReturnNothingMultipleQueries(session,
		helpers.SoftDeleteNodeQuery(uid),
		helpers.HistoryLogQuery(uid, "DELETE", userUID))

	return err

}

func (svc *PublicationsService) GetPublications(searchText string, page, pageSize int, sorting *[]helpers.Sorting, filtering *[]helpers.ColumnFilter) (result []models.Publication, totalCount int64, err error) {

	session, _ := helpers.NewNeo4jSession(*svc.neo4jDriver)

	limit := pageSize
	if limit <= 0 {
		limit = 10
	}

	skip := 0
	if page > 1 {
		skip = (page - 1) * limit
	}

	// Build query with sorting
	query := buildPublicationsQuery(searchText, skip, limit, sorting)
	result, err = helpers.GetNeo4jArrayOfNodes[models.Publication](session, query)

	helpers.ProcessArrayResult(&result, err)

	// Get total count
	countQuery := buildPublicationsCountQuery(searchText)
	totalCount, _ = helpers.GetNeo4jSingleRecordSingleValue[int64](session, countQuery)

	for i := 0; i < len(result); i++ {
		decodeAuthorsDepartments(&result[i])
		// Fetch connected researchers
		result[i].EliResearchers, _ = svc.getPublicationResearchers(result[i].Uid)
		// Fetch connected grants
		result[i].Grants, _ = svc.getPublicationGrants(result[i].Uid)
	}

	return result, totalCount, err
}

func buildPublicationsQuery(searchText string, skip, limit int, sorting *[]helpers.Sorting) helpers.DatabaseQuery {
	query := helpers.DatabaseQuery{}
	query.Parameters = make(map[string]interface{})

	// Base query with optional matches for codebook relationships
	query.Query = `
		MATCH (n:Publication) WHERE (n.deleted IS NULL OR n.deleted = false)
	`

	// Add search condition
	if searchText != "" {
		query.Query += `
			AND (toLower(n.title) CONTAINS toLower($search)
				OR toLower(n.doi) CONTAINS toLower($search)
				OR toLower(n.code) CONTAINS toLower($search)
				OR toLower(n.allAuthors) CONTAINS toLower($search)
				OR toLower(n.eliAuthors) CONTAINS toLower($search)
				OR toLower(n.keywords) CONTAINS toLower($search)
				OR toLower(n.yearOfPublication) CONTAINS toLower($search))
		`
		query.Parameters["search"] = searchText
	}

	// Optional matches for relationships
	query.Query += `
		OPTIONAL MATCH (n)-[:HAS_MEDIA_TYPE]->(mediaTypeCb:MediaType)
		OPTIONAL MATCH (n)-[:HAS_OPEN_ACCESS_TYPE]->(openAccessType:OpenAccessType)
		OPTIONAL MATCH (n)-[:HAS_PUBLISHING_COUNTRY]->(publishingCountry:Country)
		OPTIONAL MATCH (n)-[:HAS_USER_CALL]->(userCall:UserCall)
		OPTIONAL MATCH (n)-[:HAS_USER_EXPERIMENT]->(userExperimentCb:UserExperiment)
		OPTIONAL MATCH (n)-[:HAS_EXPERIMENTAL_SYSTEM]->(experimentalSystemCb:ExperimentalSystem)
		WITH n, mediaTypeCb, openAccessType, publishingCountry, userCall, userExperimentCb, experimentalSystemCb
	`

	// Add sorting
	orderBy := getPublicationsSortingClause(sorting)
	query.Query += orderBy

	// Add pagination
	query.Query += fmt.Sprintf(" SKIP %d LIMIT %d ", skip, limit)

	// Return statement
	query.Query += `
		RETURN {
			uid: n.uid,
			doi: n.doi,
			code: n.code,
			title: n.title,
			abstract: n.abstract,
			mediaType: n.mediaType,
			mediaTypeCb: CASE WHEN mediaTypeCb IS NOT NULL THEN {uid: mediaTypeCb.uid, name: mediaTypeCb.name, code: mediaTypeCb.code} ELSE null END,
			longJournalTitle: n.longJournalTitle,
			shortJournalTitle: n.shortJournalTitle,
			volume: n.volume,
			issue: n.issue,
			pages: n.pages,
			pagesCount: n.pagesCount,
			citeAs: n.citeAs,
			impactFactor: n.impactFactor,
			quartilBasis: n.quartilBasis,
			quartil: n.quartil,
			yearOfPublication: n.yearOfPublication,
			pdfFileName: n.pdfFileName,
			pdfFileUrl: n.pdfFileUrl,
			dateOfPublication: n.dateOfPublication,
			keywords: n.keywords,
			oecdFord: n.oecdFord,
			wosNumber: n.wosNumber,
			issn: n.issn,
			eissn: n.eissn,
			webLink: n.webLink,
			eidScopus: n.eidScopus,
			language: n.language,
			grant: n.grant,
			note: n.note,
			allAuthors: n.allAuthors,
			allAuthorsCount: n.allAuthorsCount,
			eliAuthors: n.eliAuthors,
			eliAuthorsCount: n.eliAuthorsCount,
			authorsDepartmentsArray: n.authorsDepartmentsArray,
			openAccessType: CASE WHEN openAccessType IS NOT NULL THEN {uid: openAccessType.uid, name: openAccessType.name} ELSE null END,
			publishingCountry: CASE WHEN publishingCountry IS NOT NULL THEN {uid: publishingCountry.uid, name: publishingCountry.name} ELSE null END,
			userCall: CASE WHEN userCall IS NOT NULL THEN {uid: userCall.uid, name: userCall.name} ELSE null END,
			userExperiment: n.userExperiment,
			userExperimentCb: CASE WHEN userExperimentCb IS NOT NULL THEN {uid: userExperimentCb.uid, name: userExperimentCb.name, code: userExperimentCb.code} ELSE null END,
			experimentalSystem: n.experimentalSystem,
			experimentalSystemCb: CASE WHEN experimentalSystemCb IS NOT NULL THEN {uid: experimentalSystemCb.uid, name: experimentalSystemCb.name, code: experimentalSystemCb.code} ELSE null END,
			updatedAt: n.updatedAt
		} as n
	`

	query.ReturnAlias = "n"
	return query
}

func buildPublicationsCountQuery(searchText string) helpers.DatabaseQuery {
	query := helpers.DatabaseQuery{}
	query.Parameters = make(map[string]interface{})

	query.Query = `
		MATCH (n:Publication) WHERE (n.deleted IS NULL OR n.deleted = false)
	`

	if searchText != "" {
		query.Query += `
			AND (toLower(n.title) CONTAINS toLower($search)
				OR toLower(n.doi) CONTAINS toLower($search)
				OR toLower(n.code) CONTAINS toLower($search)
				OR toLower(n.allAuthors) CONTAINS toLower($search)
				OR toLower(n.eliAuthors) CONTAINS toLower($search)
				OR toLower(n.keywords) CONTAINS toLower($search)
				OR toLower(n.yearOfPublication) CONTAINS toLower($search))
		`
		query.Parameters["search"] = searchText
	}

	query.Query += " RETURN count(n) as totalCount"
	query.ReturnAlias = "totalCount"
	return query
}

func getPublicationsSortingClause(sorting *[]helpers.Sorting) string {
	if sorting == nil || len(*sorting) == 0 {
		return " ORDER BY n.updatedAt DESC "
	}

	orderBy := " ORDER BY "
	for i, sort := range *sorting {
		sortField := mapPublicationSortField(sort.ID)
		direction := helpers.GetSortingDirectionString(sort.DESC)

		if i > 0 {
			orderBy += ", "
		}
		orderBy += fmt.Sprintf("%s %s", sortField, direction)
	}
	return orderBy
}

func mapPublicationSortField(fieldID string) string {
	// Map frontend field IDs to Neo4j property paths
	fieldMap := map[string]string{
		"title":             "n.title",
		"doi":               "n.doi",
		"code":              "n.code",
		"yearOfPublication": "n.yearOfPublication",
		"allAuthors":        "n.allAuthors",
		"eliAuthors":        "n.eliAuthors",
		"longJournalTitle":  "n.longJournalTitle",
		"impactFactor":      "n.impactFactor",
		"quartil":           "n.quartil",
		"updatedAt":         "n.updatedAt",
		"language":          "n.language",
		"volume":            "n.volume",
		"pagesCount":        "n.pagesCount",
	}

	if mapped, ok := fieldMap[fieldID]; ok {
		return mapped
	}
	// Default to the field as-is with n. prefix
	return "n." + fieldID
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

// Researcher methods

func (svc *PublicationsService) GetResearcherByUid(uid string) (result models.Researcher, err error) {

	session, _ := helpers.NewNeo4jSession(*svc.neo4jDriver)

	result.Uid = uid

	err = helpers.GetSingleNode(session, &result)

	return result, err
}

func (svc *PublicationsService) GetResearchers(searchText string, page, pageSize int, sorting *[]helpers.Sorting) (result []models.Researcher, totalCount int64, err error) {

	session, _ := helpers.NewNeo4jSession(*svc.neo4jDriver)

	limit := pageSize
	if limit <= 0 {
		limit = 10
	}

	skip := 0
	if page > 1 {
		skip = (page - 1) * limit
	}

	// Build query with sorting
	query := buildResearchersQuery(searchText, skip, limit, sorting)
	result, err = helpers.GetNeo4jArrayOfNodes[models.Researcher](session, query)

	helpers.ProcessArrayResult(&result, err)

	// Get total count
	countQuery := buildResearchersCountQuery(searchText)
	totalCount, _ = helpers.GetNeo4jSingleRecordSingleValue[int64](session, countQuery)

	return result, totalCount, err
}

func buildResearchersQuery(searchText string, skip, limit int, sorting *[]helpers.Sorting) helpers.DatabaseQuery {
	query := helpers.DatabaseQuery{}
	query.Parameters = make(map[string]interface{})

	query.Query = `
		MATCH (n:Researcher) WHERE (n.deleted IS NULL OR n.deleted = false)
	`

	if searchText != "" {
		query.Query += `
			AND (toLower(n.firstName) CONTAINS toLower($search)
				OR toLower(n.lastName) CONTAINS toLower($search)
				OR toLower(n.identificationNumber) CONTAINS toLower($search)
				OR toLower(n.orcid) CONTAINS toLower($search)
				OR toLower(n.scopusId) CONTAINS toLower($search)
				OR toLower(n.researcherId) CONTAINS toLower($search))
		`
		query.Parameters["search"] = searchText
	}

	query.Query += `
		OPTIONAL MATCH (n)-[:HAS_CITIZENSHIP]->(citizenship:Country)
		WITH n, citizenship
	`

	// Add sorting
	orderBy := getResearchersSortingClause(sorting)
	query.Query += orderBy

	query.Query += fmt.Sprintf(" SKIP %d LIMIT %d ", skip, limit)

	query.Query += `
		RETURN {
			uid: n.uid,
			firstName: n.firstName,
			lastName: n.lastName,
			identificationNumber: n.identificationNumber,
			orcid: n.orcid,
			scopusId: n.scopusId,
			researcherId: n.researcherId,
			citizenship: CASE WHEN citizenship IS NOT NULL THEN {uid: citizenship.uid, name: citizenship.name, code: citizenship.code} ELSE null END,
			updatedAt: n.updatedAt
		} as n
	`

	query.ReturnAlias = "n"
	return query
}

func buildResearchersCountQuery(searchText string) helpers.DatabaseQuery {
	query := helpers.DatabaseQuery{}
	query.Parameters = make(map[string]interface{})

	query.Query = `
		MATCH (n:Researcher) WHERE (n.deleted IS NULL OR n.deleted = false)
	`

	if searchText != "" {
		query.Query += `
			AND (toLower(n.firstName) CONTAINS toLower($search)
				OR toLower(n.lastName) CONTAINS toLower($search)
				OR toLower(n.identificationNumber) CONTAINS toLower($search)
				OR toLower(n.orcid) CONTAINS toLower($search)
				OR toLower(n.scopusId) CONTAINS toLower($search)
				OR toLower(n.researcherId) CONTAINS toLower($search))
		`
		query.Parameters["search"] = searchText
	}

	query.Query += " RETURN count(n) as totalCount"
	query.ReturnAlias = "totalCount"
	return query
}

func getResearchersSortingClause(sorting *[]helpers.Sorting) string {
	if sorting == nil || len(*sorting) == 0 {
		return " ORDER BY n.lastName ASC, n.firstName ASC "
	}

	orderBy := " ORDER BY "
	for i, sort := range *sorting {
		sortField := mapResearcherSortField(sort.ID)
		direction := helpers.GetSortingDirectionString(sort.DESC)

		if i > 0 {
			orderBy += ", "
		}
		orderBy += fmt.Sprintf("%s %s", sortField, direction)
	}
	return orderBy
}

func mapResearcherSortField(fieldID string) string {
	fieldMap := map[string]string{
		"firstName":            "n.firstName",
		"lastName":             "n.lastName",
		"identificationNumber": "n.identificationNumber",
		"orcid":                "n.orcid",
		"scopusId":             "n.scopusId",
		"researcherId":         "n.researcherId",
		"updatedAt":            "n.updatedAt",
	}

	if mapped, ok := fieldMap[fieldID]; ok {
		return mapped
	}
	return "n." + fieldID
}

func (svc *PublicationsService) CreateResearcher(researcher *models.Researcher, userUID string) (result models.Researcher, err error) {

	session, _ := helpers.NewNeo4jSession(*svc.neo4jDriver)

	updateQuery := helpers.DatabaseQuery{}
	updateQuery.Parameters = make(map[string]interface{})
	updateQuery.Query = `MERGE (n:Researcher {uid: $uid}) SET n.updatedAt = datetime() WITH n `
	updateQuery.Parameters["uid"] = researcher.Uid

	helpers.AutoResolveObjectToUpdateQuery(&updateQuery, *researcher, models.Researcher{}, "n")

	updateQuery.Query += ` RETURN n.uid as uid `
	updateQuery.ReturnAlias = "uid"

	historyLog := helpers.HistoryLogQuery(researcher.Uid, "CREATE", userUID)

	err = helpers.WriteNeo4jAndReturnNothingMultipleQueries(session,
		updateQuery,
		historyLog)

	return *researcher, err
}

func (svc *PublicationsService) CreateResearchers(researchers []models.Researcher, userUID string) (result []models.Researcher, err error) {

	result = make([]models.Researcher, 0)

	for i := range researchers {
		researcher := &researchers[i]
		createdResearcher, err := svc.CreateResearcher(researcher, userUID)
		if err != nil {
			return result, err
		}
		result = append(result, createdResearcher)
	}

	return result, nil
}

func (svc *PublicationsService) UpdateResearcher(researcher *models.Researcher, userUID string) (result models.Researcher, err error) {

	session, _ := helpers.NewNeo4jSession(*svc.neo4jDriver)

	oldResearcher, err := svc.GetResearcherByUid(researcher.Uid)

	if err != nil {
		return result, err
	}

	updateQuery := helpers.DatabaseQuery{}
	updateQuery.Parameters = make(map[string]interface{})
	updateQuery.Query = `MATCH (n:Researcher {uid: $uid}) SET n.updatedAt = datetime() WITH n  `
	updateQuery.Parameters["uid"] = researcher.Uid

	helpers.AutoResolveObjectToUpdateQuery(&updateQuery, *researcher, oldResearcher, "n")

	updateQuery.Query += ` RETURN n.uid as uid `
	updateQuery.ReturnAlias = "uid"

	historyLog := helpers.HistoryLogQuery(researcher.Uid, "UPDATE", userUID)

	err = helpers.WriteNeo4jAndReturnNothingMultipleQueries(session,
		updateQuery,
		historyLog)

	return *researcher, err
}

func (svc *PublicationsService) DeleteResearcher(uid string, userUID string) (err error) {

	session, _ := helpers.NewNeo4jSession(*svc.neo4jDriver)

	err = helpers.WriteNeo4jAndReturnNothingMultipleQueries(session,
		helpers.SoftDeleteNodeQuery(uid),
		helpers.HistoryLogQuery(uid, "DELETE", userUID))

	return err
}

// Publication-Researcher relationship helpers

func (svc *PublicationsService) getPublicationResearchers(uid string) ([]models.ResearcherRef, error) {
	session, _ := helpers.NewNeo4jSession(*svc.neo4jDriver)

	query := helpers.DatabaseQuery{
		Query: `MATCH (p:Publication {uid: $uid})-[:HAS_RESEARCHER]->(r:Researcher)
				WHERE r.deleted IS NULL OR r.deleted = false
				RETURN {uid: r.uid, firstName: r.firstName, lastName: r.lastName} as researcher`,
		ReturnAlias: "researcher",
		Parameters:  map[string]interface{}{"uid": uid},
	}

	result, err := helpers.GetNeo4jArrayOfNodes[models.ResearcherRef](session, query)
	return result, err
}

func (svc *PublicationsService) connectPublicationResearchers(pubUid string, researchers []models.ResearcherRef) error {
	if len(researchers) == 0 {
		return nil
	}

	session, _ := helpers.NewNeo4jSession(*svc.neo4jDriver)

	for _, res := range researchers {
		query := helpers.DatabaseQuery{
			Query: `MATCH (p:Publication {uid: $pubUid})
					MATCH (r:Researcher {uid: $resUid})
					WHERE r.deleted IS NULL OR r.deleted = false
					MERGE (p)-[:HAS_RESEARCHER]->(r)`,
			Parameters: map[string]interface{}{
				"pubUid": pubUid,
				"resUid": res.Uid,
			},
		}
		err := helpers.WriteNeo4jAndReturnNothing(session, query)
		if err != nil {
			return err
		}
	}
	return nil
}

func (svc *PublicationsService) updatePublicationResearchers(pubUid string, newResearchers []models.ResearcherRef) error {
	session, _ := helpers.NewNeo4jSession(*svc.neo4jDriver)

	// 1. Get current researchers
	oldResearchers, err := svc.getPublicationResearchers(pubUid)
	if err != nil {
		return err
	}

	// Handle nil slices
	if newResearchers == nil {
		newResearchers = []models.ResearcherRef{}
	}

	// 2. Find researchers to disconnect (in old, not in new)
	for _, oldRes := range oldResearchers {
		found := false
		for _, newRes := range newResearchers {
			if oldRes.Uid == newRes.Uid {
				found = true
				break
			}
		}
		if !found {
			// Disconnect: DELETE relationship only (not the Researcher node)
			query := helpers.DatabaseQuery{
				Query: `MATCH (p:Publication {uid: $pubUid})-[rel:HAS_RESEARCHER]->(r:Researcher {uid: $resUid})
						DELETE rel`,
				Parameters: map[string]interface{}{
					"pubUid": pubUid,
					"resUid": oldRes.Uid,
				},
			}
			helpers.WriteNeo4jAndReturnNothing(session, query)
		}
	}

	// 3. Find researchers to connect (in new, not in old)
	for _, newRes := range newResearchers {
		found := false
		for _, oldRes := range oldResearchers {
			if newRes.Uid == oldRes.Uid {
				found = true
				break
			}
		}
		if !found {
			// Connect: CREATE relationship (only if researcher is not deleted)
			query := helpers.DatabaseQuery{
				Query: `MATCH (p:Publication {uid: $pubUid})
						MATCH (r:Researcher {uid: $resUid})
						WHERE r.deleted IS NULL OR r.deleted = false
						MERGE (p)-[:HAS_RESEARCHER]->(r)`,
				Parameters: map[string]interface{}{
					"pubUid": pubUid,
					"resUid": newRes.Uid,
				},
			}
			helpers.WriteNeo4jAndReturnNothing(session, query)
		}
	}

	return nil
}

// Publication-Grant relationship helpers

func (svc *PublicationsService) getPublicationGrants(uid string) ([]models.GrantRef, error) {
	session, _ := helpers.NewNeo4jSession(*svc.neo4jDriver)

	query := helpers.DatabaseQuery{
		Query: `MATCH (p:Publication {uid: $uid})-[:HAS_GRANT]->(g:Grant)
				WHERE g.deleted IS NULL OR g.deleted = false
				RETURN {uid: g.uid, code: g.code, name: g.name} as grant`,
		ReturnAlias: "grant",
		Parameters: map[string]interface{}{
			"uid": uid,
		},
	}

	result, err := helpers.GetNeo4jArrayOfNodes[models.GrantRef](session, query)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (svc *PublicationsService) connectPublicationGrants(pubUid string, grants []models.GrantRef) error {
	session, _ := helpers.NewNeo4jSession(*svc.neo4jDriver)

	for _, grant := range grants {
		query := helpers.DatabaseQuery{
			Query: `MATCH (p:Publication {uid: $pubUid})
					MATCH (g:Grant {uid: $grantUid})
					WHERE g.deleted IS NULL OR g.deleted = false
					MERGE (p)-[:HAS_GRANT]->(g)`,
			Parameters: map[string]interface{}{
				"pubUid":   pubUid,
				"grantUid": grant.Uid,
			},
		}
		err := helpers.WriteNeo4jAndReturnNothing(session, query)
		if err != nil {
			return err
		}
	}
	return nil
}

func (svc *PublicationsService) updatePublicationGrants(pubUid string, newGrants []models.GrantRef) error {
	session, _ := helpers.NewNeo4jSession(*svc.neo4jDriver)

	// 1. Get current grants
	oldGrants, err := svc.getPublicationGrants(pubUid)
	if err != nil {
		return err
	}

	// Handle nil slices
	if newGrants == nil {
		newGrants = []models.GrantRef{}
	}

	// 2. Find grants to disconnect (in old, not in new)
	for _, oldGrant := range oldGrants {
		found := false
		for _, newGrant := range newGrants {
			if oldGrant.Uid == newGrant.Uid {
				found = true
				break
			}
		}
		if !found {
			query := helpers.DatabaseQuery{
				Query: `MATCH (p:Publication {uid: $pubUid})-[rel:HAS_GRANT]->(g:Grant {uid: $grantUid})
						DELETE rel`,
				Parameters: map[string]interface{}{
					"pubUid":   pubUid,
					"grantUid": oldGrant.Uid,
				},
			}
			helpers.WriteNeo4jAndReturnNothing(session, query)
		}
	}

	// 3. Find grants to connect (in new, not in old)
	for _, newGrant := range newGrants {
		found := false
		for _, oldGrant := range oldGrants {
			if newGrant.Uid == oldGrant.Uid {
				found = true
				break
			}
		}
		if !found {
			query := helpers.DatabaseQuery{
				Query: `MATCH (p:Publication {uid: $pubUid})
						MATCH (g:Grant {uid: $grantUid})
						WHERE g.deleted IS NULL OR g.deleted = false
						MERGE (p)-[:HAS_GRANT]->(g)`,
				Parameters: map[string]interface{}{
					"pubUid":   pubUid,
					"grantUid": newGrant.Uid,
				},
			}
			helpers.WriteNeo4jAndReturnNothing(session, query)
		}
	}

	return nil
}

// Grant CRUD methods

func (svc *PublicationsService) GetGrantByUid(uid string) (result models.Grant, err error) {
	session, _ := helpers.NewNeo4jSession(*svc.neo4jDriver)
	result.Uid = uid
	err = helpers.GetSingleNode(session, &result)
	return result, err
}

func (svc *PublicationsService) GetGrants(searchText string, page, pageSize int, facilityCode string) (result []models.Grant, totalCount int64, err error) {
	session, _ := helpers.NewNeo4jSession(*svc.neo4jDriver)

	limit := pageSize
	if limit <= 0 {
		limit = 10
	}

	skip := 0
	if page > 1 {
		skip = (page - 1) * limit
	}

	query := helpers.DatabaseQuery{
		Query: `MATCH (f:Facility {code: $facilityCode})
				MATCH (g:Grant)-[:BELONGS_TO_FACILITY]->(f)
				WHERE (g.deleted IS NULL OR g.deleted = false)
				AND (toLower(g.code) CONTAINS toLower($search) OR toLower(g.name) CONTAINS toLower($search))
				OPTIONAL MATCH (g)-[:BELONGS_TO_GROUP]->(gg:GrantGroup)
				RETURN {uid: g.uid, code: g.code, name: g.name, updatedAt: g.updatedAt, grantGroup: CASE WHEN gg IS NOT NULL THEN {uid: gg.uid, name: gg.name, code: gg.code} ELSE null END} as grant
				ORDER BY grant.name
				SKIP $skip LIMIT $limit`,
		ReturnAlias: "grant",
		Parameters: map[string]interface{}{
			"facilityCode": facilityCode,
			"search":       searchText,
			"skip":         skip,
			"limit":        limit,
		},
	}

	result, err = helpers.GetNeo4jArrayOfNodes[models.Grant](session, query)
	helpers.ProcessArrayResult(&result, err)

	if err != nil {
		return result, totalCount, err
	}

	countQuery := helpers.DatabaseQuery{
		Query: `MATCH (f:Facility {code: $facilityCode})
				MATCH (g:Grant)-[:BELONGS_TO_FACILITY]->(f)
				WHERE (g.deleted IS NULL OR g.deleted = false)
				AND (toLower(g.code) CONTAINS toLower($search) OR toLower(g.name) CONTAINS toLower($search))
				RETURN count(g) as totalCount`,
		ReturnAlias: "totalCount",
		Parameters: map[string]interface{}{
			"facilityCode": facilityCode,
			"search":       searchText,
		},
	}
	totalCount, err = helpers.GetNeo4jSingleRecordSingleValue[int64](session, countQuery)

	return result, totalCount, err
}

func (svc *PublicationsService) CreateGrant(grant *models.Grant, userUID string, facilityCode string) (result models.Grant, err error) {
	session, _ := helpers.NewNeo4jSession(*svc.neo4jDriver)

	updateQuery := helpers.DatabaseQuery{}
	updateQuery.Parameters = make(map[string]interface{})
	updateQuery.Query = `MATCH (f:Facility {code: $facilityCode})
						 MERGE (n:Grant {uid: $uid})
						 SET n.updatedAt = datetime(), n.code = $code, n.name = $name
						 MERGE (n)-[:BELONGS_TO_FACILITY]->(f)
						 WITH n `
	updateQuery.Parameters["uid"] = grant.Uid
	updateQuery.Parameters["code"] = grant.Code
	updateQuery.Parameters["name"] = grant.Name
	updateQuery.Parameters["facilityCode"] = facilityCode

	if grant.GrantGroup != nil && grant.GrantGroup.UID != "" {
		updateQuery.Query += `MATCH (gg:GrantGroup {uid: $grantGroupUid})
							  MERGE (n)-[:BELONGS_TO_GROUP]->(gg)
							  WITH n `
		updateQuery.Parameters["grantGroupUid"] = grant.GrantGroup.UID
	}

	updateQuery.Query += ` RETURN n.uid as uid `
	updateQuery.ReturnAlias = "uid"

	historyLog := helpers.HistoryLogQuery(grant.Uid, "CREATE", userUID)

	err = helpers.WriteNeo4jAndReturnNothingMultipleQueries(session, updateQuery, historyLog)

	return *grant, err
}

func (svc *PublicationsService) UpdateGrant(grant *models.Grant, userUID string) (result models.Grant, err error) {
	session, _ := helpers.NewNeo4jSession(*svc.neo4jDriver)

	oldGrant, err := svc.GetGrantByUid(grant.Uid)
	if err != nil {
		return result, err
	}

	updateQuery := helpers.DatabaseQuery{}
	updateQuery.Parameters = make(map[string]interface{})
	updateQuery.Query = `MATCH (n:Grant {uid: $uid}) SET n.updatedAt = datetime() WITH n `
	updateQuery.Parameters["uid"] = grant.Uid

	helpers.AutoResolveObjectToUpdateQuery(&updateQuery, *grant, oldGrant, "n")

	updateQuery.Query += ` RETURN n.uid as uid `
	updateQuery.ReturnAlias = "uid"

	historyLog := helpers.HistoryLogQuery(grant.Uid, "UPDATE", userUID)

	err = helpers.WriteNeo4jAndReturnNothingMultipleQueries(session, updateQuery, historyLog)

	return *grant, err
}

func (svc *PublicationsService) DeleteGrant(uid string, userUID string) (err error) {
	session, _ := helpers.NewNeo4jSession(*svc.neo4jDriver)

	err = helpers.WriteNeo4jAndReturnNothingMultipleQueries(session,
		helpers.SoftDeleteNodeQuery(uid),
		helpers.HistoryLogQuery(uid, "DELETE", userUID))

	return err
}

// Codebook autocomplete methods

func (svc *PublicationsService) GetExperimentalSystemsAutocomplete(searchText string, limit int, facilityCode string) (result []codebookModels.Codebook, err error) {
	session, _ := helpers.NewNeo4jSession(*svc.neo4jDriver)

	query := helpers.DatabaseQuery{
		Query: `MATCH(f:Facility{code:$facilityCode})
				MATCH(r:ExperimentalSystem)-[:BELONGS_TO_FACILITY]->(f)
				WHERE apoc.text.clean(r.name) CONTAINS apoc.text.clean($searchText)
				RETURN {uid: r.uid, name: r.name, code: r.code} as result
				ORDER BY result.name LIMIT $limit`,
		ReturnAlias: "result",
		Parameters: map[string]interface{}{
			"searchText":   searchText,
			"facilityCode": facilityCode,
			"limit":        limit,
		},
	}

	result, err = helpers.GetNeo4jArrayOfNodes[codebookModels.Codebook](session, query)
	return result, err
}

func (svc *PublicationsService) GetUserExperimentsAutocomplete(searchText string, limit int, facilityCode string) (result []codebookModels.Codebook, err error) {
	session, _ := helpers.NewNeo4jSession(*svc.neo4jDriver)

	query := helpers.DatabaseQuery{
		Query: `MATCH(f:Facility{code:$facilityCode})
				MATCH(r:UserExperiment)-[:BELONGS_TO_FACILITY]->(f)
				WHERE apoc.text.clean(r.name) CONTAINS apoc.text.clean($searchText)
				RETURN {uid: r.uid, name: r.name, code: r.code} as result
				ORDER BY result.name LIMIT $limit`,
		ReturnAlias: "result",
		Parameters: map[string]interface{}{
			"searchText":   searchText,
			"facilityCode": facilityCode,
			"limit":        limit,
		},
	}

	result, err = helpers.GetNeo4jArrayOfNodes[codebookModels.Codebook](session, query)
	return result, err
}

func (svc *PublicationsService) GetCountriesAutocomplete(searchText string, limit int) (result []codebookModels.Codebook, err error) {
	session, _ := helpers.NewNeo4jSession(*svc.neo4jDriver)

	query := helpers.DatabaseQuery{
		Query: `MATCH(r:Country)
				WHERE apoc.text.clean(r.name) CONTAINS apoc.text.clean($searchText)
				RETURN {uid: r.uid, name: r.name, code: r.code} as result
				ORDER BY result.name LIMIT $limit`,
		ReturnAlias: "result",
		Parameters: map[string]interface{}{
			"searchText": searchText,
			"limit":      limit,
		},
	}

	result, err = helpers.GetNeo4jArrayOfNodes[codebookModels.Codebook](session, query)
	return result, err
}
