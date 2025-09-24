package models

type CategoryList struct {
	Categories []*Category `json:"categories"`
	Pagination *Pagination `json:"pagination"`
}
