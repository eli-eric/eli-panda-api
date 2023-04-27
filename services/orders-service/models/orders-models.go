package models

import "time"

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
	UID            string      `json:"uid"`
	Name           string      `json:"name"`
	RequestNumber  string      `json:"requestNumber"`
	OrderNumber    string      `json:"orderNumber"`
	ContractNumber string      `json:"contractNumber"`
	Notes          string      `json:"notes"`
	Supplier       Supplier    `json:"supplier"`
	OrderStatus    OrderStatus `json:"orderStatus"`
}

type Supplier struct {
	UID  string `json:"uid"`
	Name string `json:"name"`
}

type OrderStatus struct {
	UID  string `json:"uid"`
	Name string `json:"name"`
}
