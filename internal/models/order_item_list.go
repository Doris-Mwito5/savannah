package models

type OrderItemList struct {
	OrderItems []*OrderItem `json:"order_item"`
	Pagination *Pagination  `json:"pagination"`
}
