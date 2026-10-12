package models

// FilterValueRange holds inclusive numeric bounds for a range filter input.
// A nil bound means "open" on that side.
type FilterValueRange struct {
	Min *float64 `json:"min"`
	Max *float64 `json:"max"`
}

// FilterStringRange holds inclusive bounds for a lexicographic (ISO string)
// range such as the free-form publication dates.
type FilterStringRange struct {
	Min *string `json:"min"`
	Max *string `json:"max"`
}

// PublicationFilterOptions carries the distinct values and bounds the
// publications filter sheet needs to populate its listboxes and range inputs
// (ELIPANDA-503 optional endpoint, consumed by ELIPANDA-504).
type PublicationFilterOptions struct {
	Years           []string                     `json:"years"`
	Quartils        []string                     `json:"quartils"`
	QuartilBases    []string                     `json:"quartilBases"`
	Languages       []string                     `json:"languages"`
	EliPublications []string                     `json:"eliPublications"`
	Ranges          map[string]FilterValueRange  `json:"ranges"`
	DateBounds      map[string]FilterStringRange `json:"dateBounds"`
}
