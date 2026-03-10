package models

// RIV export configuration — institution-specific constants for IS VaVaI (Czech R&D Information System)
// TODO: replace placeholder values with real institution data

const (
	// Institution identification
	RivInstitutionICO   = "67985831"                            // 8-digit IČO
	RivOrgUnitCode      = "_____"                               // 5-char org unit code (or _____ if not applicable)
	RivProviderCode     = "MSM"                                 // 3-char provider code
	RivInstitutionNameCZ = "Fyzikální ústav AV ČR, v. v. i."   // institution name in Czech
	RivInstitutionNameEN = "Institute of Physics of the CAS"    // institution name in English
	RivLegalType         = "verejna-vyzkumna-instituce"         // legal type of institution

	// Contact person for the delivery
	RivContactName  = "TODO: Contact Name"  // TODO: set real contact name
	RivContactEmail = "TODO@eli-beams.eu"   // TODO: set real contact email
	RivContactPhone = "+420000000000"       // TODO: set real contact phone

	// Delivery metadata
	RivDeliveryVersion = "01"               // delivery version
	RivDeliveryRef     = "TODO-REF"         // cover document reference number (cislo-jednaci)
)

// MediaTypeDruhMap maps MediaType codes to RIV <druh> attribute values
var MediaTypeDruhMap = map[string]string{
	"J": "clanek-v-periodiku",
	"C": "kapitola-v-knize",
	"D": "clanek-ve-sborniku",
}
