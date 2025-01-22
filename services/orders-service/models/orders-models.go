package models

import (
	codebookModels "panda/apigateway/services/codebook-service/models"
	"time"
)

type OrderListItem struct {
	UID                    string                   `json:"uid"`
	Name                   string                   `json:"name"`
	RequestNumber          string                   `json:"requestNumber"`
	OrderNumber            string                   `json:"orderNumber"`
	ContractNumber         string                   `json:"contractNumber"`
	Notes                  string                   `json:"notes"`
	Supplier               string                   `json:"supplier"`
	OrderStatus            *codebookModels.Codebook `json:"orderStatus"`
	DeliveryStatus         int                      `json:"deliveryStatus"` //0 - none delivered, 1 - partially delivered, 2 - delivered
	OrderDate              time.Time                `json:"orderDate"`
	LastUpdateTime         time.Time                `json:"lastUpdateTime"`
	LastUpdateBy           string                   `json:"lastUpdateBy"`
	Requestor              string                   `json:"requestor"`
	ProcurementResponsible string                   `json:"procurementResponsible"`
}

type OrderDetail struct {
	UID                    string                   `json:"uid" neo4j:"ignore"`
	Name                   string                   `json:"name" neo4j:"prop,name"`
	RequestNumber          *string                  `json:"requestNumber" neo4j:"prop,requestNumber"`
	OrderNumber            *string                  `json:"orderNumber" neo4j:"prop,orderNumber"`
	ContractNumber         *string                  `json:"contractNumber" neo4j:"prop,contractNumber"`
	Notes                  *string                  `json:"notes" neo4j:"prop,notes"`
	Supplier               *codebookModels.Codebook `json:"supplier" neo4j:"rel,Supplier,HAS_SUPPLIER,uid,supplier"`
	OrderStatus            *codebookModels.Codebook `json:"orderStatus" neo4j:"rel,OrderStatus,HAS_ORDER_STATUS,uid,orderStatus"`
	Requestor              *codebookModels.Codebook `json:"requestor" neo4j:"rel,Employee,HAS_REQUESTOR,uid,requestor"`
	ProcurementResponsible *codebookModels.Codebook `json:"procurementResponsible" neo4j:"rel,Employee,HAS_PROCUREMENT_RESPONSIBLE,uid,procurementResponsible"`
	OrderDate              time.Time                `json:"orderDate" neo4j:"prop,orderDate"`
	OrderLines             []OrderLine              `json:"orderLines"`
	LastUpdateTime         time.Time                `json:"lastUpdateTime"`
}

type OrderLine struct {
	Name            string                   `json:"name"`
	UID             string                   `json:"uid"`
	CatalogueNumber string                   `json:"catalogueNumber"`
	CatalogueUID    string                   `json:"catalogueUid"`
	System          *codebookModels.Codebook `json:"system"`
	ItemUsage       *codebookModels.Codebook `json:"itemUsage"`
	Location        *codebookModels.Codebook `json:"location"`
	Price           *float64                 `json:"price"`
	Currency        *string                  `json:"currency"`
	EUN             *string                  `json:"eun"`
	IsDelivered     bool                     `json:"isDelivered"`
	DeliveredTime   *time.Time               `json:"deliveredTime"`
	SerialNumber    *string                  `json:"serialNumber"`
	Notes           *string                  `json:"notes" neo4j:"prop,notes"`
	LastUpdateTime  *time.Time               `json:"lastUpdateTime"`
}

type ItemForEunPrint struct {
	EUN             string `json:"eun"`
	Name            string `json:"name"`
	Manufacturer    string `json:"manufacturer"`
	CatalogueNumber string `json:"catalogueNumber"`
	SerialNumber    string `json:"serialNumber"`
	Quantity        int    `json:"quantity"`
	Location        string `json:"location"`
}

type OrderLineDelivery struct {
	IsDelivered  bool    `json:"isDelivered"`
	SerialNumber *string `json:"serialNumber"`
	EUN          *string `json:"eun"`
}

type OrderLineMinMax struct {
	Min int `json:"min"`
	Max int `json:"max"`
}
