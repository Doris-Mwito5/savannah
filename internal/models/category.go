package models

import "github/Doris-Mwito5/savannah-pos/internal/custom_types"

type Category struct {
	custom_types.SequentialIdentifier
	Name     string  `json:"name"`
	ParentID *int64 `json:"parent_id"`
	ShopID   *string `json:"shop_id"`
	custom_types.Timestamps
}
