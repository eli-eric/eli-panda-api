package models

// RIV export configuration — institution-specific constants for IS VaVaI (Czech R&D Information System)

const (
	// Institution identification
	RivInstitutionICO    = "10974938"
	RivOrgUnitCode       = "_____"
	RivInstitutionNameCZ = "Extreme Light Infrastructure ERIC (ELI ERIC)"
	RivInstitutionNameEN = "Extreme Light Infrastructure ERIC (ELI ERIC)"
	RivLegalType         = "pravnicka-osoba-jina"

	// Contact person for the delivery
	RivContactName  = "Ladislav Půst"
	RivContactEmail = "ladislav.pust@eli-beams.eu"
	RivContactPhone = "+420775620803"

	// Result metadata
	RivDuvernostUdaju = "verejne-pristupne"

	// Delivery metadata
	RivDeliveryMode    = "R" // R=new, Z=update, V=deletion
	RivDeliveryVersion = "01"
	RivDeliveryRef     = "TODO-REF" // set per delivery
)

// MediaTypeDruhMap maps MediaType codes to RIV <druh> attribute values
var MediaTypeDruhMap = map[string]string{
	"J": "clanek-v-periodiku",
	"C": "kapitola-v-knize",
	"D": "clanek-ve-sborniku",
}

// OpenAccessCodeToRIV maps OpenAccessType.code to RIV zpusob-publikovani long-form values
var OpenAccessCodeToRIV = map[string]string{
	"A": "open-access",
	"B": "embargoed-access",
	"C": "restricted-access",
	"D": "metadata-only",
}

// RivLanguageMap maps English language names to ISO 639-2 three-letter codes
var RivLanguageMap = map[string]string{
	// Germanic
	"English":    "eng",
	"German":     "ger",
	"Dutch":      "dut",
	"Swedish":    "swe",
	"Norwegian":  "nor",
	"Danish":     "dan",
	"Afrikaans":  "afr",
	// Romance
	"French":     "fre",
	"Spanish":    "spa",
	"Portuguese": "por",
	"Italian":    "ita",
	"Romanian":   "rum",
	"Catalan":    "cat",
	// Slavic
	"Czech":      "cze",
	"Slovak":     "slo",
	"Polish":     "pol",
	"Russian":    "rus",
	"Ukrainian":  "ukr",
	"Bulgarian":  "bul",
	"Croatian":   "hrv",
	"Serbian":    "srp",
	"Slovenian":  "slv",
	// Other European
	"Hungarian":  "hun",
	"Finnish":    "fin",
	"Greek":      "gre",
	"Turkish":    "tur",
	"Albanian":   "alb",
	"Lithuanian": "lit",
	"Latvian":    "lav",
	"Estonian":   "est",
	// Asian
	"Chinese":    "chi",
	"Japanese":   "jpn",
	"Korean":     "kor",
	"Arabic":     "ara",
	"Hindi":      "hin",
	"Persian":    "per",
	"Hebrew":     "heb",
	// Other
	"Latin":      "lat",
	"Esperanto":  "epo",
}
