package dtos

import "github/Doris-Mwito5/savannah-pos/internal/custom_types"

type CreateOrderForm struct {
	CustomerName    string                     `json:"customer_name"`
	CustomerEmail   string                     `json:"customer_email"`
	ReferenceNumber string                     `json:"reference_number"`
	PhoneNumber     string                     `json:"phone_number"`
	OrderStatus     custom_types.OrderStatus   `json:"order_status"`
	OrderMedium     custom_types.OrderMedium   `json:"order_medium"`
	PaymentMethod   custom_types.PaymentMethod `json:"payment_method"`
	CustomerID      *int64                     `json:"customer_id"`
	ShopID          string                     `json:"shop_id"`
	Items           []OrderItemForm            `json:"items"`
	TotalAmount     float64                    `json:"total_amount"`
	Discount        *float64                    `json:"discount"`
}

type UpdateOrderForm struct {
	PhoneNumber   *string  `json:"phone_number"`
	OrderStatus   *string  `json:"order_status"`
	OrderMedium   *string  `json:"order_medium"`
	PaymentMethod *string  `json:"payment_method"`
	TotalAmount   *float64 `json:"total_amount"`
	Discount      *float64 `json:"discount"`
}

type OrderItemForm struct {
	OrderID     int64   `json:"order_id"`
	ProductID   int64   `json:"product_id"`
	UnitPrice   float64 `json:"unit_price"`
	Quantity    int64   `json:"quantity"`
	TotalAmount float64 `json:"total_amount"`
}
