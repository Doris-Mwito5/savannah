package dtos

import "github/Doris-Mwito5/savannah-pos/internal/custom_types"

type CreateProductForm struct {
	Name           string                   `json:"name"`
	Description    *string                  `json:"description,omitempty"`
	WholesalePrice float64                  `json:"wholesale_price"`
	RetailPrice    float64                  `json:"retail_price"`
	CategoryID     int64                    `json:"category_id"`
	ProductImage   *string                  `json:"product_image,omitempty"`
	Stock          int64                    `json:"stock,omitempty"`
	ProductType    custom_types.ProductType `json:"product_type"`
}

type UpdateProductForm struct {
	Name           *string  `json:"name"`
	Description    *string  `json:"description,omitempty"`
	WholesalePrice *float64 `json:"wholesale_price"`
	RetailPrice    *float64 `json:"retail_price"`
	CategoryID     *int64   `json:"category_id"`
	ProductImage   *string  `json:"product_image,omitempty"`
	Stock          *int64   `json:"stock,omitempty"`
	ProductType    *string  `json:"product_type"`
}
