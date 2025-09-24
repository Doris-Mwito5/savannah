package models

import "github/Doris-Mwito5/savannah-pos/internal/custom_types"

type Order struct {
	custom_types.SequentialIdentifier
	ReferenceNumber string                     `json:"reference_number"`
	PhoneNumber     string                     `json:"phone_number"`
	OrderStatus     custom_types.OrderStatus   `json:"order_status"`
	OrderMedium     custom_types.OrderMedium   `json:"order_medium"`
	PaymentMethod   custom_types.PaymentMethod `json:"payment_method"`
	CustomerID      *int64                     `json:"customer_id"`
	ShopID          string                     `json:"shop_id"`
	TotalItems      int                        `json:"total_items"`
	TotalAmount     float64                    `json:"total_amount"`
	Discount        *float64                   `json:"discount"`
	custom_types.Timestamps
}
