package dtos

import "github/Doris-Mwito5/savannah-pos/internal/custom_types"

type CreateCustomerForm struct {
	Name         string                    `json:"name"`
	Email        string                    `json:"email"`
	PhoneNumber  string                    `json:"phone_number"`
	CustomerType custom_types.CustomerType `json:"customer_type"`
	ShopID       string                    `json:"shop_id"`
}

type UpdateCustomerForm struct {
	Name         *string `json:"name"`
	Email        *string `json:"email"`
	PhoneNumber  *string `json:"phone_number"`
	CustomerType *string `json:"customer_type"`
	ShopID       *string `json:"shop_id"`
}
