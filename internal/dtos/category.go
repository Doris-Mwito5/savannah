package dtos

type CreateCategoryForm struct {
	Name     string `json:"name"`
	ParentID *int64 `json:"parent_id"`
	ShopID   string `json:"shop_id"`
}

type UpdateCategoryForm struct {
	Name     *string `json:"name"`
	ParentID *int64 `json:"parent_id"`
	ShopID   *string `json:"shop_id"`
}
