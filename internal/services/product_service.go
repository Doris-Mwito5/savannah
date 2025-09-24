package services

import (
	"context"
	"github/Doris-Mwito5/savannah-pos/internal/db"
	"github/Doris-Mwito5/savannah-pos/internal/domain"
	"github/Doris-Mwito5/savannah-pos/internal/dtos"
	"github/Doris-Mwito5/savannah-pos/internal/models"
	"github/Doris-Mwito5/savannah-pos/internal/null"
)

type (
	ProductService interface {
		CreateProduct(ctx context.Context, dB db.DB, form *dtos.CreateProductForm) (*models.Product, error)
		ProductByID(ctx context.Context, dB db.DB, productID int64) (*models.Product, error)
		UpdateProduct(ctx context.Context, dB db.DB, productID int64, form *dtos.UpdateProductForm) (*models.Product, error)
		DeleteProduct(ctx context.Context, dB db.DB, productID int64) (*models.Product, error)
		ListProducts(ctx context.Context, dB db.DB, shopID string, filter *models.Filter) (*models.ProductList, error)
		GetAveragePriceByCategory(	ctx context.Context, dB db.DB, categoryID int64) (float64, error)
	}

	productService struct {
		store *domain.Store
	}
)

func NewProductService(store *domain.Store) ProductService {
	return &productService{
		store: store,
	}
}

func (s *productService) CreateProduct(
	ctx context.Context,
	dB db.DB,
	form *dtos.CreateProductForm,
) (*models.Product, error) {
	//fetch the category
	category, err := s.store.CategoryDomain.CategoryByID(ctx, dB, form.CategoryID)
	if err != nil {
		return &models.Product{}, err
	}
	product := &models.Product{
		Name:           form.Name,
		Description:    form.Description,
		WholesalePrice: form.WholesalePrice,
		RetailPrice:    form.RetailPrice,
		CategoryID:     category.ID,
		ProductImage:   form.ProductImage,
		Stock:          form.Stock,
		ProductType:    form.ProductType,
	}

	err = s.store.ProductDomain.CreateProduct(ctx, dB, product)
	if err != nil {
		return &models.Product{}, err
	}

	return product, nil
}

func (s *productService) ProductByID(
	ctx context.Context,
	dB db.DB,
	productID int64,
) (*models.Product, error) {

	return s.store.ProductDomain.ProductByID(ctx, dB, productID)
}

func (s *productService) UpdateProduct(
	ctx context.Context,
	dB db.DB,
	productID int64,
	form *dtos.UpdateProductForm,
) (*models.Product, error) {

	product, err := s.store.ProductDomain.ProductByID(ctx, dB, productID)
	if err != nil {
		return &models.Product{}, err
	}

	if form.RetailPrice != nil {
		product.RetailPrice = null.ValueFromNull(form.RetailPrice)
	}

	if form.WholesalePrice != nil {
		product.WholesalePrice = null.ValueFromNull(form.WholesalePrice)
	}

	if form.Stock != nil {
		product.Stock = null.ValueFromNull(form.Stock)
	}

	err = s.store.ProductDomain.CreateProduct(ctx, dB, product)
	if err != nil {
		return &models.Product{}, err
	}
	return product, nil
}

func (s *productService) ListProducts(
	ctx context.Context,
	dB db.DB,
	shopID string,
	filter *models.Filter,
) (*models.ProductList, error) {

	products, err := s.store.ProductDomain.ListProducts(ctx, dB, shopID, filter)
	if err != nil {
		return &models.ProductList{}, err
	}

	count, err := s.store.ProductDomain.ProductCount(ctx, dB, shopID, filter)
	if err != nil {
		return &models.ProductList{}, err
	}

	productList := &models.ProductList{
		Products: products,
		Pagination: models.NewPagination(
			count,
			filter.Page,
			filter.Per,
		),
	}
	return productList, nil
}

func (s *productService) DeleteProduct(
	ctx context.Context,
	dB db.DB,
	productID int64,
) (*models.Product, error) {

	product, err := s.store.ProductDomain.ProductByID(ctx, dB, productID)
	if err != nil {
		return &models.Product{}, err
	}

	err = s.store.ProductDomain.DeleteProduct(ctx, dB, productID)
	if err != nil {
		return &models.Product{}, err
	}

	return product, nil
}

func (s *productService) GetAveragePriceByCategory(
	ctx context.Context,
	dB db.DB,
	categoryID int64,
) (float64, error) {
	
	_, err := s.store.CategoryDomain.CategoryByID(ctx, dB, categoryID)
	if err != nil {
		return 0, err
	}
	
	return s.store.ProductDomain.GetAveragePriceByCategory(ctx, dB, categoryID)
}