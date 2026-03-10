package publicationsservice

import (
	"encoding/xml"
	"fmt"
	"panda/apigateway/helpers"
	"panda/apigateway/services/publications-service/models"
	"strings"

	"github.com/rs/zerolog/log"
)

// rivPublicationRow represents a single flat row returned by the RIV export Cypher query.
// Each row is one researcher per publication — must be aggregated by publication UID.
type rivPublicationRow struct {
	// Publication fields
	Uid               string  `json:"uid"`
	Code              string  `json:"code"`
	Title             string  `json:"title"`
	Abstract          string  `json:"abstract"`
	Language          string  `json:"language"`
	YearOfPublication string  `json:"yearOfPublication"`
	Doi               string  `json:"doi"`
	WebLink           string  `json:"webLink"`
	Keywords          string  `json:"keywords"`
	OecdFord          *string `json:"oecdFord"`
	Volume            int     `json:"volume"`
	Issue             *int    `json:"issue"`
	Pages             string  `json:"pages"`
	PagesCount        int     `json:"pagesCount"`
	LongJournalTitle  string  `json:"longJournalTitle"`
	AllAuthorsCount   int     `json:"allAuthorsCount"`
	EliAuthorsCount   int     `json:"eliAuthorsCount"`
	WosNumber         *string `json:"wosNumber"`
	Issn              *string `json:"issn"`
	EIssn             *string `json:"eissn"`
	EidScopus         *string `json:"eidScopus"`
	MediaTypeCode     *string `json:"mediaTypeCode"`

	// Related node codes
	PublishingCountryCode *string `json:"publishingCountryCode"`
	OpenAccessCode        *string `json:"openAccessCode"`
	PublishFormatCode     *string `json:"publishFormatCode"`
	ConferenceScopeCode   *string `json:"conferenceScopeCode"`

	// Type C/D fields
	Publisher       *string `json:"publisher"`
	PublishPlace    *string `json:"publishPlace"`
	Isbn            *string `json:"isbn"`
	BookTitle       *string `json:"bookTitle"`
	BookPagesCount  *int    `json:"bookPagesCount"`
	EditionVolume   *string `json:"editionVolume"`
	ProceedingsIsbn *string `json:"proceedingsIsbn"`
	ConferenceDate  *string `json:"conferenceDate"`
	ConferencePlace *string `json:"conferencePlace"`

	// Researcher fields (one per row)
	ResearcherUid      *string `json:"researcherUid"`
	ResearcherFirst    *string `json:"researcherFirst"`
	ResearcherLast     *string `json:"researcherLast"`
	ResearcherIdNumber *string `json:"researcherIdNumber"`
	ResearcherOrcid    *string `json:"researcherOrcid"`
	ResearcherScopus   *string `json:"researcherScopus"`
	ResearcherRID      *string `json:"researcherRID"`
	CitizenshipCode    *string `json:"citizenshipCode"`
}

// rivAggregatedPublication holds one publication with all its researchers
type rivAggregatedPublication struct {
	row         rivPublicationRow
	researchers []rivResearcherData
}

type rivResearcherData struct {
	FirstName            string
	LastName             string
	IdentificationNumber string
	ORCID                string
	ScopusID             string
	ResearcherID         string
	CitizenshipCode      string
}

func (svc *PublicationsService) ExportRiv(year string, provider string) ([]byte, string, error) {
	pubs, warnings, err := svc.buildRivData(year, provider)
	if err != nil {
		return nil, "", err
	}

	if len(warnings) > 0 {
		log.Warn().Int("warnings", len(warnings)).Msg("RIV export has validation warnings")
	}

	dodavka := buildRivDodavka(year, provider, pubs)

	xmlBytes, err := xml.MarshalIndent(dodavka, "", "  ")
	if err != nil {
		return nil, "", fmt.Errorf("XML marshal error: %w", err)
	}

	xmlOutput := []byte(xml.Header + string(xmlBytes))

	yy := year
	if len(year) >= 4 {
		yy = year[2:4]
	}
	filename := fmt.Sprintf("RIV%s-%s-%s,R%s.xml",
		yy, provider, models.RivInstitutionICO, models.RivDeliveryVersion)

	return xmlOutput, filename, nil
}

func (svc *PublicationsService) ValidateRiv(year string, provider string) (models.RivValidationResult, error) {
	pubs, warnings, err := svc.buildRivData(year, provider)
	if err != nil {
		return models.RivValidationResult{}, err
	}

	validCount := len(pubs)
	// Count publications that have warnings
	pubsWithWarnings := make(map[string]bool)
	for _, w := range warnings {
		pubsWithWarnings[w.PublicationCode] = true
	}

	return models.RivValidationResult{
		TotalPublications: len(pubs),
		ValidPublications: validCount - len(pubsWithWarnings),
		Warnings:          warnings,
	}, nil
}

// buildRivData fetches publications for the year and provider, aggregates flat rows, and validates.
func (svc *PublicationsService) buildRivData(year string, provider string) ([]rivAggregatedPublication, []models.RivValidationWarning, error) {
	session, _ := helpers.NewNeo4jSession(*svc.neo4jDriver)

	query := helpers.DatabaseQuery{
		Query: `
			MATCH (p:Publication)-[:HAS_GRANT]->(g:Grant)-[:BELONGS_TO_GROUP]->(gg:GrantGroup)
			WHERE p.yearOfPublication = $year
			  AND (p.deleted IS NULL OR p.deleted = false)
			  AND gg.code = $provider
			OPTIONAL MATCH (p)-[:HAS_MEDIA_TYPE]->(mt:MediaType)
			OPTIONAL MATCH (p)-[:HAS_PUBLISHING_COUNTRY]->(pc:Country)
			OPTIONAL MATCH (p)-[:HAS_OPEN_ACCESS_TYPE]->(oat:OpenAccessType)
			OPTIONAL MATCH (p)-[:HAS_PUBLISH_FORMAT]->(pf:PublishFormat)
			OPTIONAL MATCH (p)-[:HAS_CONFERENCE_SCOPE]->(cs:ConferenceScope)
			OPTIONAL MATCH (p)-[:HAS_RESEARCHER]->(res:Researcher)
			  WHERE (res.deleted IS NULL OR res.deleted = false)
			OPTIONAL MATCH (res)-[:HAS_CITIZENSHIP]->(rc:Country)
			RETURN {
				uid: p.uid, code: p.code, title: p.title, abstract: p.abstract,
				language: p.language, yearOfPublication: p.yearOfPublication,
				doi: p.doi, webLink: p.webLink, keywords: p.keywords,
				oecdFord: p.oecdFord, volume: p.volume, issue: p.issue,
				pages: p.pages, pagesCount: p.pagesCount,
				longJournalTitle: p.longJournalTitle,
				allAuthorsCount: p.allAuthorsCount, eliAuthorsCount: p.eliAuthorsCount,
				wosNumber: p.wosNumber, issn: p.issn, eissn: p.eissn,
				eidScopus: p.eidScopus,
				mediaTypeCode: mt.code,
				publishingCountryCode: pc.code,
				openAccessCode: oat.code,
				publishFormatCode: pf.code,
				conferenceScopeCode: cs.code,
				publisher: p.publisher, publishPlace: p.publishPlace,
				isbn: p.isbn, bookTitle: p.bookTitle,
				bookPagesCount: p.bookPagesCount, editionVolume: p.editionVolume,
				proceedingsIsbn: p.proceedingsIsbn,
				conferenceDate: p.conferenceDate, conferencePlace: p.conferencePlace,
				researcherUid: res.uid, researcherFirst: res.firstName,
				researcherLast: res.lastName,
				researcherIdNumber: res.identificationNumber,
				researcherOrcid: res.orcid, researcherScopus: res.scopusId,
				researcherRID: res.researcherId,
				citizenshipCode: rc.code
			} as row
			ORDER BY p.uid, res.lastName, res.firstName
		`,
		ReturnAlias: "row",
		Parameters:  map[string]interface{}{"year": year, "provider": provider},
	}

	rows, err := helpers.GetNeo4jArrayOfNodes[rivPublicationRow](session, query)
	if err != nil {
		return nil, nil, fmt.Errorf("RIV export query error: %w", err)
	}

	// Aggregate flat rows by publication UID
	pubMap := make(map[string]*rivAggregatedPublication)
	pubOrder := make([]string, 0)

	for _, row := range rows {
		agg, exists := pubMap[row.Uid]
		if !exists {
			agg = &rivAggregatedPublication{row: row, researchers: make([]rivResearcherData, 0)}
			pubMap[row.Uid] = agg
			pubOrder = append(pubOrder, row.Uid)
		}
		if row.ResearcherUid != nil && *row.ResearcherUid != "" {
			// Deduplicate researchers (same researcher can appear multiple times)
			isDupe := false
			for _, existing := range agg.researchers {
				if existing.FirstName == derefStr(row.ResearcherFirst) && existing.LastName == derefStr(row.ResearcherLast) {
					isDupe = true
					break
				}
			}
			if !isDupe {
				agg.researchers = append(agg.researchers, rivResearcherData{
					FirstName:            derefStr(row.ResearcherFirst),
					LastName:             derefStr(row.ResearcherLast),
					IdentificationNumber: derefStr(row.ResearcherIdNumber),
					ORCID:                derefStr(row.ResearcherOrcid),
					ScopusID:             derefStr(row.ResearcherScopus),
					ResearcherID:         derefStr(row.ResearcherRID),
					CitizenshipCode:      derefStr(row.CitizenshipCode),
				})
			}
		}
	}

	// Build ordered result + validate
	pubs := make([]rivAggregatedPublication, 0, len(pubOrder))
	warnings := make([]models.RivValidationWarning, 0)

	idCodes := make(map[string]bool)
	for _, uid := range pubOrder {
		agg := pubMap[uid]
		pubs = append(pubs, *agg)

		// Validation
		code := agg.row.Code
		mediaType := derefStr(agg.row.MediaTypeCode)

		// Check identifikacni-kod uniqueness
		idCode := buildIdentifikacniKod(agg.row.Code, agg.row.YearOfPublication)
		if idCodes[idCode] {
			warnings = append(warnings, models.RivValidationWarning{PublicationCode: code, Message: "duplicate identifikacni-kod: " + idCode})
		}
		idCodes[idCode] = true

		if len(agg.row.Abstract) < 64 {
			warnings = append(warnings, models.RivValidationWarning{PublicationCode: code, Message: "abstract shorter than 64 chars"})
		}
		if strings.TrimSpace(agg.row.Keywords) == "" {
			warnings = append(warnings, models.RivValidationWarning{PublicationCode: code, Message: "no keywords"})
		}
		if agg.row.OecdFord == nil || *agg.row.OecdFord == "" {
			warnings = append(warnings, models.RivValidationWarning{PublicationCode: code, Message: "no oecdFord"})
		}
		if mediaType == "" {
			warnings = append(warnings, models.RivValidationWarning{PublicationCode: code, Message: "no mediaType code set"})
		}

		// Type-specific validation
		switch mediaType {
		case "J":
			if (agg.row.Issn == nil || *agg.row.Issn == "") && (agg.row.EIssn == nil || *agg.row.EIssn == "") {
				warnings = append(warnings, models.RivValidationWarning{PublicationCode: code, Message: "type J: no ISSN or eISSN"})
			}
			if agg.row.WebLink == "" {
				warnings = append(warnings, models.RivValidationWarning{PublicationCode: code, Message: "type J: no webLink"})
			}
		case "D":
			hasIsbn := agg.row.ProceedingsIsbn != nil && *agg.row.ProceedingsIsbn != ""
			hasIssn := agg.row.Issn != nil && *agg.row.Issn != ""
			hasEissn := agg.row.EIssn != nil && *agg.row.EIssn != ""
			if !hasIsbn && !hasIssn && !hasEissn {
				warnings = append(warnings, models.RivValidationWarning{PublicationCode: code, Message: "type D: no proceedingsIsbn, issn, or eissn"})
			}
		}

		// Researcher validation
		for _, res := range agg.researchers {
			if res.IdentificationNumber == "" {
				warnings = append(warnings, models.RivValidationWarning{
					PublicationCode: code,
					Message:         fmt.Sprintf("researcher %s %s: no identification number", res.FirstName, res.LastName),
				})
			}
		}
	}

	return pubs, warnings, nil
}

func buildRivDodavka(year string, provider string, pubs []rivAggregatedPublication) models.RivDodavka {
	vysledky := make([]models.RivVysledek, 0, len(pubs))
	for _, pub := range pubs {
		v := buildRivVysledek(pub)
		vysledky = append(vysledky, v)
	}

	return models.RivDodavka{
		Struktura: "RIV26A",
		Zahlavi: models.RivZahlavi{
			Rozsah: models.RivRozsah{
				InformacniOblast: "RIV",
				ObdobiSberu:      year,
				Predkladatel: models.RivPredkladatel{
					Subjekt: models.RivSubjektPredkladatel{
						Druh: models.RivLegalType,
						ICO:  models.RivInstitutionICO,
						Nazvy: []models.RivNazev{
							{Jazyk: "#ORIG", Value: models.RivInstitutionNameCZ},
							{Jazyk: "eng", Value: models.RivInstitutionNameEN},
						},
					},
				},
			},
			Dodavatel: models.RivDodavatel{
				Subjekt: models.RivSubjektDodavatel{Kod: provider},
				Pracovnik: models.RivPracovnik{
					Osoba: models.RivOsoba{
						CeleJmeno: models.RivContactName,
						Kontakt: models.RivKontakt{
							Telefon: models.RivTelefon{Druh: "telefon", Value: models.RivContactPhone},
							Email:   models.RivContactEmail,
						},
					},
				},
			},
			Verze:    models.RivDeliveryVersion,
			Pruvodka: models.RivPruvodka{CisloJednaci: models.RivDeliveryRef},
		},
		Obsah: models.RivObsah{Vysledky: vysledky},
	}
}

func buildRivVysledek(pub rivAggregatedPublication) models.RivVysledek {
	row := pub.row
	mediaType := derefStr(row.MediaTypeCode)
	druh := models.MediaTypeDruhMap[mediaType]

	v := models.RivVysledek{
		IdentifikacniKod: buildIdentifikacniKod(row.Code, row.YearOfPublication),
		DuvernostUdaju:   "S",
		RokUplatneni:     row.YearOfPublication,
		KontrolniKod:     "0",
		Druh:             druh,
		Jazyk:            row.Language,
	}

	// Titles
	v.Nazvy = []models.RivNazev{{Jazyk: "eng", Value: row.Title}}
	if row.Language != "" && row.Language != "eng" {
		v.Nazvy = append(v.Nazvy, models.RivNazev{Jazyk: "#ORIG", Value: row.Title})
	}

	// Annotations
	v.Anotace = []models.RivAnotace{{Jazyk: "eng", Value: row.Abstract}}
	if row.Language != "" && row.Language != "eng" {
		v.Anotace = append(v.Anotace, models.RivAnotace{Jazyk: "#ORIG", Value: row.Abstract})
	}

	// Link and DOI
	v.Odkaz = row.WebLink
	v.Doi = stripDoiPrefix(row.Doi)

	// Authors
	v.Autori = buildRivAutori(pub)

	// Klasifikace
	v.Klasifikace = buildRivKlasifikace(row)

	// Navaznosti
	v.Navaznosti = models.RivNavaznosti{
		Navaznost: models.RivNavaznost{
			DruhVztahu:             "byl-dosazen-pri-reseni",
			InstitucionalniPodpora: models.RivEmpty{},
		},
	}

	// Type-specific elements
	switch mediaType {
	case "J":
		buildTypeJ(&v, row)
	case "C":
		buildTypeC(&v, row)
	case "D":
		buildTypeD(&v, row)
	}

	return v
}

func buildRivAutori(pub rivAggregatedPublication) models.RivAutori {
	autori := models.RivAutori{
		PocetCelkem:   pub.row.AllAuthorsCount,
		PocetDomacich: pub.row.EliAuthorsCount,
	}

	for _, res := range pub.researchers {
		autor := models.RivAutor{
			JeDomaci: "true",
			Jmeno:    res.FirstName,
			Prijmeni: res.LastName,
		}

		isCzech := res.CitizenshipCode == "CZ"

		if !isCzech && res.CitizenshipCode != "" {
			autor.CiziStatniPrislusnik = &models.RivEmpty{}
			autor.StatniPrislusnost = res.CitizenshipCode
		}

		if isCzech {
			autor.RodneCislo = res.IdentificationNumber
		} else {
			autor.IdentifikacniCislo = res.IdentificationNumber
		}

		autor.ORCID = res.ORCID
		autor.ScopusID = res.ScopusID
		autor.ResearcherID = res.ResearcherID

		autori.Autori = append(autori.Autori, autor)
	}

	return autori
}

func buildRivKlasifikace(row rivPublicationRow) models.RivKlasifikace {
	klas := models.RivKlasifikace{}

	if row.OecdFord != nil && *row.OecdFord != "" {
		klas.Obory = []models.RivObor{
			{Postaveni: "hlavni", Ciselnik: "OblastiOECD", Value: *row.OecdFord},
		}
	}

	if row.Keywords != "" {
		parts := strings.Split(row.Keywords, ";")
		for _, kw := range parts {
			trimmed := strings.TrimSpace(kw)
			if trimmed != "" {
				klas.KlicoveSlova = append(klas.KlicoveSlova, models.RivKlicoveSlovo{
					Jazyk: "eng",
					Value: trimmed,
				})
			}
		}
	}

	return klas
}

func buildTypeJ(v *models.RivVysledek, row rivPublicationRow) {
	// Subtype
	wos := derefStr(row.WosNumber)
	eid := derefStr(row.EidScopus)

	if wos != "" {
		v.Poddruh = "clanek-wos"
	} else if eid != "" {
		v.Poddruh = "clanek-scopus"
	} else {
		v.Poddruh = "clanek-ostatni"
	}

	// Periodikum
	v.Periodikum = &models.RivPeriodikum{
		ISSN:  derefStr(row.Issn),
		EISSN: derefStr(row.EIssn),
		Nazev: row.LongJournalTitle,
	}
	if row.PublishingCountryCode != nil && *row.PublishingCountryCode != "" {
		v.Periodikum.Vydavatel = &models.RivVydavatel{Stat: *row.PublishingCountryCode}
	}

	// Rocnik
	v.Rocnik = &models.RivRocnik{}
	if row.Volume > 0 {
		v.Rocnik.Value = fmt.Sprintf("%d", row.Volume)
	} else {
		v.Rocnik.StatusUdaje = "neuvedeno"
	}

	// Cislo
	v.Cislo = &models.RivCislo{}
	if row.Issue != nil && *row.Issue > 0 {
		v.Cislo.Value = fmt.Sprintf("%d", *row.Issue)
	} else {
		v.Cislo.StatusUdaje = "neuvedeno"
	}

	// Pages
	if row.PagesCount > 0 || row.Pages != "" {
		v.Strany = &models.RivStrany{Pocet: row.PagesCount, Rozsah: row.Pages}
	}

	// WOS/Scopus identifiers
	v.KodUtIsi = stripWosPrefix(wos)
	v.EID = stripEidPrefix(eid)

	// Open access
	v.ZpusobPublikovani = derefStr(row.OpenAccessCode)
}

func buildTypeC(v *models.RivVysledek, row rivPublicationRow) {
	kniha := &models.RivKniha{
		Nazev:       derefStr(row.BookTitle),
		ISBN:        derefStr(row.Isbn),
		FormaVydani: derefStr(row.PublishFormatCode),
	}

	// EditionVolume
	ev := derefStr(row.EditionVolume)
	if ev != "" {
		kniha.EdiceCisloSvazku = &models.RivOptionalField{Value: ev}
	} else {
		kniha.EdiceCisloSvazku = &models.RivOptionalField{StatusUdaje: "neuvedeno"}
	}

	// PublishPlace
	pp := derefStr(row.PublishPlace)
	if pp != "" {
		kniha.MistoVydani = &models.RivOptionalField{Value: pp}
	} else {
		kniha.MistoVydani = &models.RivOptionalField{StatusUdaje: "neuvedeno"}
	}

	// Publisher
	pub := derefStr(row.Publisher)
	if pub != "" {
		kniha.Nakladatel = &models.RivNakladatel{Nazev: &models.RivOptionalField{Value: pub}}
	} else {
		kniha.Nakladatel = &models.RivNakladatel{Nazev: &models.RivOptionalField{StatusUdaje: "neuvedeno"}}
	}

	// Book total pages
	if row.BookPagesCount != nil && *row.BookPagesCount > 0 {
		kniha.Strany = &models.RivStrany{Pocet: *row.BookPagesCount}
	}

	v.Kniha = kniha

	// Chapter pages
	if row.PagesCount > 0 || row.Pages != "" {
		v.Strany = &models.RivStrany{Pocet: row.PagesCount, Rozsah: row.Pages}
	}

	// Identifiers
	v.EID = stripEidPrefix(derefStr(row.EidScopus))
	v.KodUtIsi = stripWosPrefix(derefStr(row.WosNumber))
}

func buildTypeD(v *models.RivVysledek, row rivPublicationRow) {
	sbornik := &models.RivSbornik{
		Nazev:       row.LongJournalTitle,
		ISBN:        derefStr(row.ProceedingsIsbn),
		ISSN:        derefStr(row.Issn),
		EISSN:       derefStr(row.EIssn),
		FormaVydani: derefStr(row.PublishFormatCode),
	}

	pp := derefStr(row.PublishPlace)
	if pp != "" {
		sbornik.MistoVydani = &models.RivOptionalField{Value: pp}
	} else {
		sbornik.MistoVydani = &models.RivOptionalField{StatusUdaje: "neuvedeno"}
	}

	pub := derefStr(row.Publisher)
	if pub != "" {
		sbornik.Nakladatel = &models.RivNakladatel{Nazev: &models.RivOptionalField{Value: pub}}
	} else {
		sbornik.Nakladatel = &models.RivNakladatel{Nazev: &models.RivOptionalField{StatusUdaje: "neuvedeno"}}
	}

	v.Sbornik = sbornik

	// Conference info
	confDate := derefStr(row.ConferenceDate)
	confPlace := derefStr(row.ConferencePlace)
	confScope := derefStr(row.ConferenceScopeCode)

	akce := &models.RivAkce{}
	if len(confDate) == 4 {
		akce.Konani.VRoce = confDate
	} else if len(confDate) >= 10 {
		akce.Konani.Zahajeni = confDate
	}
	akce.Konani.Misto = confPlace
	akce.Ucastnici.Klasifikace = confScope
	v.Akce = akce

	// Pages
	if row.PagesCount > 0 || row.Pages != "" {
		v.Strany = &models.RivStrany{Pocet: row.PagesCount, Rozsah: row.Pages}
	}

	// Identifiers
	v.EID = stripEidPrefix(derefStr(row.EidScopus))
	v.KodUtIsi = stripWosPrefix(derefStr(row.WosNumber))
}

// GetRivProviders returns grant group providers that have publications for the given year
func (svc *PublicationsService) GetRivProviders(year string) ([]models.RivProvider, error) {
	session, _ := helpers.NewNeo4jSession(*svc.neo4jDriver)

	query := helpers.DatabaseQuery{
		Query: `
			MATCH (p:Publication)-[:HAS_GRANT]->(g:Grant)-[:BELONGS_TO_GROUP]->(gg:GrantGroup)
			WHERE p.yearOfPublication = $year
			  AND (p.deleted IS NULL OR p.deleted = false)
			  AND gg.code IS NOT NULL
			  AND gg.code <> "OTHER"
			RETURN {code: gg.code, name: gg.name, publicationCount: COUNT(DISTINCT p)} as provider
			ORDER BY provider.code
		`,
		ReturnAlias: "provider",
		Parameters:  map[string]interface{}{"year": year},
	}

	result, err := helpers.GetNeo4jArrayOfNodes[models.RivProvider](session, query)
	if err != nil {
		return nil, fmt.Errorf("RIV providers query error: %w", err)
	}

	if result == nil {
		result = []models.RivProvider{}
	}

	return result, nil
}

// --- Helper functions ---

func buildIdentifikacniKod(pubCode string, year string) string {
	yy := year
	if len(year) >= 4 {
		yy = year[2:4]
	}
	return fmt.Sprintf("RIV/%s:%s/%s:%s", models.RivInstitutionICO, models.RivOrgUnitCode, yy, pubCode)
}

func derefStr(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}

func stripDoiPrefix(doi string) string {
	doi = strings.TrimPrefix(doi, "https://doi.org/")
	doi = strings.TrimPrefix(doi, "http://doi.org/")
	return doi
}

func stripWosPrefix(wos string) string {
	return strings.TrimPrefix(wos, "WOS:")
}

func stripEidPrefix(eid string) string {
	return strings.TrimPrefix(eid, "EID=")
}
