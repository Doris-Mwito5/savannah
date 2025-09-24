package domain

type Store struct {
	CustomerDomain  CustomerDomain
	CategoryDomain  CategoryDomain
	ProductDomain   ProductDomain
	OrderDomain     OrderDomain
	OrderItemDomain OrderItemDomain
}

func NewStore() *Store {
	return &Store{
		CustomerDomain:  NewCustomerDomain(),
		CategoryDomain:  NewCategoryDomain(),
		ProductDomain:   NewProductDomain(),
		OrderDomain:     NewOrderDomain(),
		OrderItemDomain: NewOrderItemDomain(),
	}
}
