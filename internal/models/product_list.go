package models

type ProductList struct {
	Products   []*Product  `json:"products"`
	Pagination *Pagination `json:"pagination"`
}