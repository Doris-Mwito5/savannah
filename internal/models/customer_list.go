package models

type CustomerList struct {
	Customers  []*Customer `json:"customers"`
	Pagination *Pagination `json:"pagination"`
}
