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
	DeliveryStatus int       `json:"deliveryStatus"`
	OrderDate      time.Time `json:"orderDate"`
	LastUpdateTime time.Time `json:"lastUpdateTime"`
	LastUpdateBy   string    `json:"lastUpdateBy"`
}

type OrderDetail struct {
	UID            string                   `json:"uid" neo4j:"ignore"`
	Name           string                   `json:"name" neo4j:"prop,name"`
	RequestNumber  *string                  `json:"requestNumber" neo4j:"prop,requestNumber"`
	OrderNumber    *string                  `json:"orderNumber" neo4j:"prop,orderNumber"`
	ContractNumber *string                  `json:"contractNumber" neo4j:"prop,contractNumber"`
	Notes          *string                  `json:"notes" neo4j:"prop,notes"`
	Supplier       *codebookModels.Codebook `json:"supplier" neo4j:"rel,Supplier,HAS_SUPPLIER,supplier"`
	OrderStatus    *codebookModels.Codebook `json:"orderStatus" neo4j:"rel,OrderStatus,HAS_ORDER_STATUS,orderStatus"`
	Requester      *codebookModels.Codebook `json:"requester" neo4j:"rel,Employee,HAS_REQUESTER,requester"`
	Procurementer  *codebookModels.Codebook `json:"procurementer" neo4j:"rel,Employee,HAS_PROCUREMENTER,procurementer"`
	OrderDate      time.Time                `json:"orderDate" neo4j:"prop,orderDate"`
	OrderLines     []OrderLine              `json:"orderLines"`
}

type OrderLine struct {
	Name            string                   `json:"name"`
	UID             string                   `json:"uid"`
	CatalogueNumber string                   `json:"catalogueNumber"`
	CatalogueUID    string                   `json:"catalogueUid"`
	System          *codebookModels.Codebook `json:"system"`
	ItemUsgae       *codebookModels.Codebook `json:"itemUsage"`
	Price           float64                  `json:"price"`
	Currency        string                   `json:"currency"`
	EUN             string                   `json:"eun"`
	IsDelivered     bool                     `json:"isDelivered"`
	DeliveredTime   *time.Time               `json:"deliveredTime"`
}
