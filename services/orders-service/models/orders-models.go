package models

import (
	codebookModels "panda/apigateway/services/codebook-service/models"
	"time"
)

type OrderListItem struct {
	UID            string    `json:"uid"`
	Name           string    `json:"name"`
	RequestNumber  string    `json:"requestNumber"`
	OrderNumber    string    `json:"orderNumber"`
	ContractNumber string    `json:"contractNumber"`
	Notes          string    `json:"notes"`
	Supplier       string    `json:"supplier"`
	OrderStatus    string    `json:"orderStatus"`
	OrderDate      time.Time `json:"orderDate"`
	LastUpdateTime time.Time `json:"lastUpdateTime"`
	LastUpdateBy   string    `json:"lastUpdateBy"`
}

type OrderDetail struct {
	UID            string                  `json:"uid"`
	Name           string                  `json:"name"`
	RequestNumber  string                  `json:"requestNumber"`
	OrderNumber    string                  `json:"orderNumber"`
	ContractNumber string                  `json:"contractNumber"`
	Notes          string                  `json:"notes"`
	Supplier       codebookModels.Codebook `json:"supplier"`
	OrderStatus    codebookModels.Codebook `json:"orderStatus"`
	OrderDate      time.Time               `json:"orderDate"`
	OrderLines     []OrderLine             `json:"orderLines"`
}

type OrderLine struct {
	Name            string                  `json:"name"`
	UID             string                  `json:"uid"`
	CatalogueNumber string                  `json:"catalogueNumber"`
	CatalogueUID    string                  `json:"catalogueUid"`
	System          codebookModels.Codebook `json:"system"`
	PriceEUR        float64                 `json:"priceEur"`
}
