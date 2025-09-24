package models

import "github/Doris-Mwito5/savannah-pos/internal/custom_types"

type OrderItem struct {
	custom_types.SequentialIdentifier
	OrderID     int64   `json:"order_id"`
	ProductID   int64   `json:"product_id"`
	UnitPrice   float64 `json:"unit_price"`
	Quantity    int64   `json:"quantity"`
	TotalAmount float64 `json:"total_amount"`
	custom_types.Timestamps
}
