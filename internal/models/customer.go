package models

import "github/Doris-Mwito5/savannah-pos/internal/custom_types"

type Customer struct {
	custom_types.SequentialIdentifier
	Name         string                    `json:"name"`
	Email        string                    `json:"email"`
	PhoneNumber  string                    `json:"phone_number"`
	CustomerType custom_types.CustomerType `json:"customer_type"`
	ShopID       string                    `json:"shop_id"`
	custom_types.Timestamps
}
