package models

type CatalogueStatistics struct {
	FacilityName                   *string `json:"facilityName"`
	Total                          int     `json:"total"`
	SparePartsCount                int     `json:"sparePartsCount"`
	InSystemPartsCount             int     `json:"inSystemPartsCount"`
	ExperimentalLoanPoolPartsCount int     `json:"experimentalLoanPoolPartsCount"`
	TestAndMeasurementPartsCount   int     `json:"testAndMeasurementPartsCount"`
	StockItemsCount                int     `json:"stockItemsCount"`
	OthersCount                    int     `json:"othersCount"`
}
