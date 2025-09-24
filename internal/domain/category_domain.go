package domain

import (
	"context"
	"fmt"
	"github/Doris-Mwito5/savannah-pos/internal/apperr"
	"github/Doris-Mwito5/savannah-pos/internal/db"
	"github/Doris-Mwito5/savannah-pos/internal/models"
	"github/Doris-Mwito5/savannah-pos/internal/null"
	"github/Doris-Mwito5/savannah-pos/internal/utils"
	"strings"
)

const (
	createCategorySQL  = "INSERT INTO categories (name, parent_id, shop_id, created_at, updated_at) VALUES($1, $2, $3, $4, $5) RETURNING(id)"
	getCategoriesSQL   = "SELECT id, name, parent_id, shop_id, created_at, updated_at FROM categories"
	getCategoryByIDSQL = getCategoriesSQL + " WHERE id = $1"
	updateCategorySQL  = "UPDATE categories SET name = $1, parent_id = $2, shop_id = $3, updated_at = $4 WHERE id = $5"
	deleteCategorySQL  = "DELETE FROM categories where id = $1"
	getCategoriesCountSQL = "SELECT COUNT(id) FROM categories"
)

type (
	CategoryDomain interface {
		CreateCategory(ctx context.Context, operations db.SQLOperations, category *models.Category) error
		CategoryByID(ctx context.Context, operations db.SQLOperations, categoryID int64) (*models.Category, error)
		ListShopCateories(ctx context.Context, opearations db.SQLOperations, shopID string, filter *models.Filter) ([]*models.Category, error)
		ShopCategoriesCount(ctx context.Context, opearations db.SQLOperations, shopID string, filter *models.Filter) (int, error)
		DeleteCategory(ctx context.Context, operations db.SQLOperations, categoryID int64) error
	}

	categoryDomain struct{}
)

func NewCategoryDomain() CategoryDomain {
	return &categoryDomain{}
}

func (d *categoryDomain) CreateCategory(
	ctx context.Context,
	operations db.SQLOperations,
	category *models.Category,
) error {

	category.Touch()
	if category.IsNew() {
		err := operations.QueryRowContext(
			ctx,
			createCategorySQL,
			category.Name,
			category.ParentID,
			category.ShopID,
			category.CreatedAt,
			category.UpdatedAt,
		).Scan(&category.ID)
		if err != nil {
			return apperr.NewDatabaseError(
				err,
			).LogErrorMessage("save category query err: %v", err)
		}
		return nil
	}

	_, err := operations.ExecContext(
		ctx,
		updateCategorySQL,
		category.Name,
		category.ParentID,
		category.ShopID,
		category.UpdatedAt,
		category.ID,
	)
	if err != nil {
		return apperr.NewDatabaseError(
			err,
		).LogErrorMessage("update category query err: %v", err)
	}
	return nil
}
func (d *categoryDomain) CategoryByID(
	ctx context.Context,
	operations db.SQLOperations,
	categoryID int64,
) (*models.Category, error) {

	row := operations.QueryRowContext(
		ctx,
		getCategoryByIDSQL,
		categoryID,
	)

	return d.scanRow(row)
}

func (d *categoryDomain) ListShopCateories(
	ctx context.Context,
	opearations db.SQLOperations,
	shopID string,
	filter *models.Filter,
) ([]*models.Category, error) {

	filter.ShopID = null.NullValue(shopID)
	query, args := d.buildQuery(getCategoriesSQL, filter)
	rows, err := opearations.QueryContext(
		ctx,
		query,
		args...,
	)
	if err != nil {
		return []*models.Category{}, apperr.NewDatabaseError(
			err,
		).LogErrorMessage("category query row err: %v", err)
	}

	defer rows.Close()

	categories := make([]*models.Category, 0)

	for rows.Next() {
		category, err := d.scanRow(rows)
		if err != nil {
			return []*models.Category{}, err
		}
		categories = append(categories, category)
	}
	if rows.Err() != nil {
		return []*models.Category{}, apperr.NewDatabaseError(
			rows.Err(),
		).LogErrorMessage(
			"list categories err: %v",
			err,
		)
	}

	return categories, nil
}

func (d *categoryDomain) ShopCategoriesCount(
	ctx context.Context,
	operations db.SQLOperations,
	shopID string, 
	filter *models.Filter,
) (int, error) {
	filter.ShopID = null.NullValue(shopID)
	query, args := d.buildQuery(getCategoriesCountSQL, filter.NoPagination())

	rows := operations.QueryRowContext(
		ctx,
		query,
		args...,
	)

	var count int

	err := rows.Scan(&count)
	if err != nil {
		return 0, apperr.NewDatabaseError(
			err,
		).LogErrorMessage(
			"category count by query row contex err: %v",
			err,
		)
	}
	return count, nil

}

func (d *categoryDomain) DeleteCategory(
	ctx context.Context,
	operations db.SQLOperations,
	categoryID int64,
) error {

	_, err := operations.ExecContext(
		ctx,
		deleteCategorySQL,
		categoryID,
	)
	if err != nil {
		return apperr.NewDatabaseError(
			err,
		).LogErrorMessage("delete category err")
	}

	return nil
}

func (d *categoryDomain) buildQuery(
	query string,
	filter *models.Filter,
) (string, []interface{}) {

	args := make([]interface{}, 0)
	conditions := make([]string, 0)
	counter := utils.NewPlaceholder()

	if filter.ShopID != nil {

		condition := fmt.Sprintf("shop_id = $%d", counter.Touch())
		args = append(args, null.ValueFromNull(filter.ShopID))

		conditions = append(conditions, condition)
	}

	if len(conditions) > 0 {
		query += " WHERE " + strings.Join(conditions, " AND ")
	}

	if filter.Page > 0 && filter.Per > 0 {
		query += fmt.Sprintf(" LIMIT $%d OFFSET $%d", counter.Touch(), counter.Touch())
		args = append(args, filter.Per, (filter.Page-1)*filter.Per)
	}

	return query, args
}

func (d *categoryDomain) scanRow(
	row db.RowScanner,
) (*models.Category, error) {

	var category models.Category

	err := row.Scan(
		&category.ID,
		&category.Name,
		&category.ParentID,
		&category.ShopID,
		&category.CreatedAt,
		&category.UpdatedAt,
	)
	if err != nil {
		return &models.Category{}, apperr.NewDatabaseError(
			err,
		).LogErrorMessage("scan row err: %v", err)
	}

	return &category, nil
}
