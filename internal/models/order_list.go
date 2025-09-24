package models

type OrderList struct {
	Orders     []*Order    `json:"orders"`
	Pagination *Pagination `json:"pagination"`
}
