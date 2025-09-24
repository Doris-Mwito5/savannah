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
	CategoryService interface {
		CreateCategory(ctx context.Context, dB db.DB, form *dtos.CreateCategoryForm) (*models.Category, error)
		CategoryByID(ctx context.Context, dB db.DB, categoryID int64) (*models.Category, error)
		DeleteCategory(ctx context.Context, dB db.DB, categoryID int64) (*models.Category, error)
		ListCategories(ctx context.Context, dB db.DB, shopID string, filter *models.Filter) (*models.CategoryList, error)
	}

	categoryService struct {
		store *domain.Store
	}
)

func NewCategoryService(store *domain.Store) CategoryService {
	return &categoryService{
		store: store,
	}
}

func (s *categoryService) CreateCategory(
	ctx context.Context,
	dB db.DB,
	form *dtos.CreateCategoryForm,
) (*models.Category, error) {

	category := &models.Category{
		Name:     form.Name,
		ParentID: form.ParentID,
		ShopID:   null.NullValue(form.ShopID),
	}

	err := s.store.CategoryDomain.CreateCategory(ctx, dB, category)
	if err != nil {
		return nil, err
	}

	return category, nil
}

func (s *categoryService) CategoryByID(
	ctx context.Context,
	dB db.DB,
	categoryID int64,
) (*models.Category, error) {
	return s.store.CategoryDomain.CategoryByID(ctx, dB, categoryID)
}

func (s *categoryService) DeleteCategory(
	ctx context.Context,
	dB db.DB,
	categoryID int64,
) (*models.Category, error) {

	category, err := s.store.CategoryDomain.CategoryByID(ctx, dB, categoryID)
	if err != nil {
		return &models.Category{}, err
	}

	err = s.store.CategoryDomain.DeleteCategory(ctx, dB, categoryID)
	if err != nil {
		return &models.Category{}, err
	}

	return category, nil
}

func (s *categoryService) ListCategories(
	ctx context.Context,
	dB db.DB,
	settingID string,
	filter *models.Filter,
) (*models.CategoryList, error) {

	categories, err := s.store.CategoryDomain.ListShopCateories(ctx, dB, settingID, filter)
	if err != nil {
		return &models.CategoryList{}, err
	}

	count, err := s.store.CategoryDomain.ShopCategoriesCount(ctx, dB, settingID, filter)
	if err != nil {
		return &models.CategoryList{}, err
	}

	categoryList := &models.CategoryList{
		Categories: categories,
		Pagination: models.NewPagination(count, filter.Page, filter.Per),
	}
	return categoryList, nil
}
