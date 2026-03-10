package models

import "encoding/xml"

// RIV26A XML structure models for Czech IS VaVaI export

// RivDodavka is the root element <dodavka>
type RivDodavka struct {
	XMLName   xml.Name       `xml:"dodavka"`
	Struktura string         `xml:"struktura,attr"`
	Zahlavi   RivZahlavi     `xml:"zahlavi"`
	Obsah     RivObsah       `xml:"obsah"`
}

// RivObsah wraps the list of results
type RivObsah struct {
	Vysledky []RivVysledek `xml:"vysledek"`
}

// RivZahlavi is the <zahlavi> header element
type RivZahlavi struct {
	Rozsah    RivRozsah    `xml:"rozsah"`
	Dodavatel RivDodavatel `xml:"dodavatel"`
	Verze     string       `xml:"verze"`
	Pruvodka  RivPruvodka  `xml:"pruvodka"`
}

type RivRozsah struct {
	InformacniOblast string          `xml:"informacni-oblast"`
	ObdobiSberu      string          `xml:"obdobi-sberu"`
	Predkladatel     RivPredkladatel `xml:"predkladatel"`
}

type RivPredkladatel struct {
	Subjekt RivSubjektPredkladatel `xml:"subjekt"`
}

type RivSubjektPredkladatel struct {
	Druh  string     `xml:"druh"`
	ICO   string     `xml:"ICO"`
	Nazvy []RivNazev `xml:"nazev"`
}

type RivNazev struct {
	Jazyk string `xml:"jazyk,attr"`
	Value string `xml:",chardata"`
}

type RivDodavatel struct {
	Subjekt    RivSubjektDodavatel `xml:"subjekt"`
	Pracovnik  RivPracovnik        `xml:"pracovnik-povereny-pripravou-dodavky"`
}

type RivSubjektDodavatel struct {
	Kod string `xml:"kod"`
}

type RivPracovnik struct {
	Osoba RivOsoba `xml:"osoba"`
}

type RivOsoba struct {
	CeleJmeno string     `xml:"cele-jmeno"`
	Kontakt   RivKontakt `xml:"kontakt"`
}

type RivKontakt struct {
	Telefon RivTelefon `xml:"telefonni-cislo"`
	Email   string     `xml:"emailova-adresa"`
}

type RivTelefon struct {
	Druh  string `xml:"druh,attr"`
	Value string `xml:",chardata"`
}

type RivPruvodka struct {
	CisloJednaci string `xml:"cislo-jednaci,attr"`
}

// RivVysledek is the <vysledek> element — one per publication
type RivVysledek struct {
	IdentifikacniKod string `xml:"identifikacni-kod,attr"`
	DuvernostUdaju   string `xml:"duvernost-udaju,attr"`
	RokUplatneni     string `xml:"rok-uplatneni,attr"`
	KontrolniKod     string `xml:"kontrolni-kod,attr"`
	Druh             string `xml:"druh,attr"`

	// Children — order matters for RIV schema
	Jazyk      string          `xml:"jazyk"`
	Nazvy      []RivNazev      `xml:"nazev"`
	Anotace    []RivAnotace    `xml:"anotace"`
	Odkaz      string          `xml:"odkaz,omitempty"`
	Doi        string          `xml:"doi,omitempty"`
	Autori     RivAutori       `xml:"autori"`
	Klasifikace RivKlasifikace `xml:"klasifikace"`
	Navaznosti  RivNavaznosti  `xml:"navaznosti"`

	// Type-specific elements (only one set per vysledek)
	// Type J
	Poddruh    string         `xml:"poddruh,omitempty"`
	Periodikum *RivPeriodikum `xml:"periodikum,omitempty"`
	Rocnik     *RivRocnik     `xml:"rocnik,omitempty"`
	Cislo      *RivCislo      `xml:"cislo,omitempty"`

	// Type C
	Kniha *RivKniha `xml:"kniha,omitempty"`

	// Type D
	Sbornik *RivSbornik `xml:"sbornik,omitempty"`
	Akce    *RivAkce    `xml:"akce,omitempty"`

	// Shared trailing elements
	Strany   *RivStrany `xml:"strany,omitempty"`
	KodUtIsi string     `xml:"kod-ut-isi,omitempty"`
	EID      string     `xml:"EID,omitempty"`

	// Type J only trailing
	ZpusobPublikovani string `xml:"zpusob-publikovani,omitempty"`
}

type RivAnotace struct {
	Jazyk string `xml:"jazyk,attr"`
	Value string `xml:",chardata"`
}

// RivAutori represents <autori> with counts
type RivAutori struct {
	PocetCelkem   int        `xml:"pocet-celkem,attr"`
	PocetDomacich int        `xml:"pocet-domacich,attr"`
	Autori        []RivAutor `xml:"autor"`
}

// RivAutor represents a single <autor> element
type RivAutor struct {
	JeDomaci string `xml:"je-domaci,attr"`

	Jmeno                string          `xml:"jmeno"`
	Prijmeni             string          `xml:"prijmeni"`
	CiziStatniPrislusnik *RivEmpty       `xml:"cizi-statni-prislusnik,omitempty"`
	StatniPrislusnost    string          `xml:"statni-prislusnost,omitempty"`
	RodneCislo           string          `xml:"rodne-cislo,omitempty"`
	IdentifikacniCislo   string          `xml:"identifikacni-cislo,omitempty"`
	ORCID                string          `xml:"orcid,omitempty"`
	ScopusID             string          `xml:"scopusid,omitempty"`
	ResearcherID         string          `xml:"researcherid,omitempty"`
}

// RivEmpty is used for self-closing empty elements like <cizi-statni-prislusnik/>
type RivEmpty struct{}

// RivKlasifikace holds OECD field + keywords
type RivKlasifikace struct {
	Obory         []RivObor         `xml:"obor"`
	KlicoveSlova  []RivKlicoveSlovo `xml:"klicove-slovo"`
}

type RivObor struct {
	Postaveni string `xml:"postaveni,attr"`
	Ciselnik  string `xml:"ciselnik,attr"`
	Value     string `xml:",chardata"`
}

type RivKlicoveSlovo struct {
	Jazyk string `xml:"jazyk,attr"`
	Value string `xml:",chardata"`
}

// RivNavaznosti represents funding linkages
type RivNavaznosti struct {
	Navaznost RivNavaznost `xml:"navaznost"`
}

type RivNavaznost struct {
	DruhVztahu                  string   `xml:"druh-vztahu,attr"`
	InstitucionalniPodpora      RivEmpty `xml:"institucionalni-podpora-na-rozvoj-VO"`
}

// Type J specific
type RivPeriodikum struct {
	ISSN      string          `xml:"ISSN,omitempty"`
	EISSN     string          `xml:"eISSN,omitempty"`
	Nazev     string          `xml:"nazev"`
	Vydavatel *RivVydavatel   `xml:"vydavatel,omitempty"`
}

type RivVydavatel struct {
	Stat string `xml:"stat"`
}

type RivRocnik struct {
	StatusUdaje string `xml:"status-udaje,attr,omitempty"`
	Value       string `xml:",chardata"`
}

type RivCislo struct {
	StatusUdaje string `xml:"status-udaje,attr,omitempty"`
	Value       string `xml:",chardata"`
}

type RivStrany struct {
	Pocet  int    `xml:"pocet,attr"`
	Rozsah string `xml:"rozsah,omitempty"`
}

// Type C specific
type RivKniha struct {
	Nazev              string            `xml:"nazev"`
	ISBN               string            `xml:"ISBN,omitempty"`
	FormaVydani        string            `xml:"forma-vydani,omitempty"`
	EdiceCisloSvazku   *RivOptionalField `xml:"edice-cislo-svazku"`
	MistoVydani        *RivOptionalField `xml:"misto-vydani"`
	Nakladatel         *RivNakladatel    `xml:"nakladatel"`
	Strany             *RivStrany        `xml:"strany,omitempty"`
}

// RivOptionalField represents an element that can have a value or status-udaje="neuvedeno"
type RivOptionalField struct {
	StatusUdaje string `xml:"status-udaje,attr,omitempty"`
	Value       string `xml:",chardata"`
}

type RivNakladatel struct {
	Nazev *RivOptionalField `xml:"nazev"`
}

// Type D specific
type RivSbornik struct {
	Nazev       string            `xml:"nazev"`
	ISBN        string            `xml:"ISBN,omitempty"`
	ISSN        string            `xml:"ISSN,omitempty"`
	EISSN       string            `xml:"eISSN,omitempty"`
	FormaVydani string            `xml:"forma-vydani,omitempty"`
	MistoVydani *RivOptionalField `xml:"misto-vydani"`
	Nakladatel  *RivNakladatel    `xml:"nakladatel"`
}

type RivAkce struct {
	Konani    RivKonani    `xml:"konani"`
	Ucastnici RivUcastnici `xml:"ucastnici"`
}

type RivKonani struct {
	Zahajeni string `xml:"zahajeni,omitempty"` // YYYY-MM-DD format
	VRoce    string `xml:"v-roce,omitempty"`   // YYYY format
	Misto    string `xml:"misto,omitempty"`
}

type RivUcastnici struct {
	Klasifikace string `xml:"klasifikace-podle-statni-prislusnosti"`
}

// RivProvider represents a grant group provider available for RIV export
type RivProvider struct {
	Code             string `json:"code"`
	Name             string `json:"name"`
	PublicationCount int    `json:"publicationCount"`
}

// RivValidationWarning represents a single validation warning
type RivValidationWarning struct {
	PublicationCode string `json:"publicationCode"`
	Message         string `json:"message"`
}

// RivValidationResult is the response for the validate endpoint
type RivValidationResult struct {
	TotalPublications int                    `json:"totalPublications"`
	ValidPublications int                    `json:"validPublications"`
	Warnings          []RivValidationWarning `json:"warnings"`
}
