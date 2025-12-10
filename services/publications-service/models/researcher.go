package models

import (
	codebookModels "panda/apigateway/services/codebook-service/models"
	"time"
)

type Researcher struct {
	Uid          string                   `json:"uid" neo4j:"key,uid"`                                                  // uid is the unique identifier of the researcher
	FirstName    string                   `json:"firstName" neo4j:"prop,firstName"`                                     // firstName is the first name of the researcher
	LastName     string                   `json:"lastName" neo4j:"prop,lastName"`                                       // lastName is the last name of the researcher
	BirthNumber  *string                  `json:"birthNumber" neo4j:"prop,birthNumber"`                                 // birthNumber is the birth number of the researcher
	ORCID        *string                  `json:"orcid" neo4j:"prop,orcid"`                                             // orcid is the ORCID identifier of the researcher
	ScopusId     *string                  `json:"scopusId" neo4j:"prop,scopusId"`                                       // scopusId is the Scopus identifier of the researcher
	ResearcherID *string                  `json:"researcherId" neo4j:"prop,researcherId"`                               // researcherId is the ResearcherID of the researcher
	Citizenship  *codebookModels.Codebook `json:"citizenship" neo4j:"rel,Country,HAS_CITIZENSHIP,uid,citizenship"`      // citizenship is the country of citizenship of the researcher
	UpdatedAt    *time.Time               `json:"updatedAt" neo4j:"prop,updatedAt"`                                     // updatedAt is the time when the researcher was last updated
}
