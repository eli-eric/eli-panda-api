package models

// RIV export configuration — institution-specific constants for IS VaVaI (Czech R&D Information System)

const (
	// Institution identification
	RivInstitutionICO    = "10974938"
	RivOrgUnitCode       = "_____"
	RivInstitutionNameCZ = "Extreme Light Infrastructure ERIC (ELI ERIC)"
	RivInstitutionNameEN = "Extreme Light Infrastructure ERIC (ELI ERIC)"
	RivLegalType         = "konsorcium-evropske-vyzkumne-infrastruktury"

	// Contact person for the delivery
	RivContactName  = "Ladislav Půst"
	RivContactEmail = "ladislav.pust@eli-beams.eu"
	RivContactPhone = "+420775620803"

	// Delivery metadata
	RivDeliveryVersion = "01"
	RivDeliveryRef     = "TODO-REF" // set per delivery
)

// MediaTypeDruhMap maps MediaType codes to RIV <druh> attribute values
var MediaTypeDruhMap = map[string]string{
	"J": "clanek-v-periodiku",
	"C": "kapitola-v-knize",
	"D": "clanek-ve-sborniku",
}
