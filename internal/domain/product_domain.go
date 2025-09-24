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
	createProductSQL = "INSERT INTO products(name, description, wholesale_price, retail_price, category_id, product_image, product_type, stock, created_at, updated_at) VALUES($1, $2, $3, $4, $5, $6, $7, $8, $9, $10) RETURNING id"
	getProductsSQL   = `SELECT p.id, p.name, p.description, p.wholesale_price, p.retail_price, p.category_id, p.product_image, p.product_type, p.stock, p.created_at, p.updated_at FROM products p
		LEFT JOIN categories c ON c.id = p.category_id`
	getProductByIDSQL    = getProductsSQL + " WHERE p.id = $1"
	getInventoryCountSQL = "SELECT COUNT(p.id) FROM products p LEFT JOIN categories c ON c.id = p.category_id"
	updateProductSQL     = `UPDATE products SET name = $1, description = $2, wholesale_price = $3, retail_price = $4, category_id = $5, product_image = $6, product_type = $7, stock = $8, updated_at = $9 WHERE id = $10`
	deleteProductSQL     = "DELETE FROM products WHERE id = $1"
	
	getCategoryHierarchySQL = `
		WITH RECURSIVE category_tree AS (
			SELECT id, name, parent_id, 0 as level
			FROM categories 
			WHERE id = $1
			
			UNION ALL
			
			SELECT c.id, c.name, c.parent_id, ct.level + 1
			FROM categories c
			INNER JOIN category_tree ct ON c.parent_id = ct.id
		)
		SELECT id FROM category_tree`
		
	getAveragePriceByCategorySQL = `
		WITH RECURSIVE category_tree AS (
			SELECT id, name, parent_id, 0 as level
			FROM categories 
			WHERE id = $1
			
			UNION ALL
			
			SELECT c.id, c.name, c.parent_id, ct.level + 1
			FROM categories c
			INNER JOIN category_tree ct ON c.parent_id = ct.id
		)
		SELECT COALESCE(AVG(p.retail_price), 0) as average_price
		FROM products p
		INNER JOIN category_tree ct ON p.category_id = ct.id`
)

type (
	ProductDomain interface {
		CreateProduct(ctx context.Context, operations db.SQLOperations, product *models.Product) error
		ProductByID(ctx context.Context, operations db.SQLOperations, productID int64) (*models.Product, error)
		ListProducts(ctx context.Context, operations db.SQLOperations, shopID string, filter *models.Filter) ([]*models.Product, error)
		ProductCount(ctx context.Context, operations db.SQLOperations, shopID string, filter *models.Filter) (int, error)
		DeleteProduct(ctx context.Context, operations db.SQLOperations, productID int64) error
		GetAveragePriceByCategory(ctx context.Context, operations db.SQLOperations, categoryID int64) (float64, error)
		GetCategoryHierarchy(ctx context.Context, operations db.SQLOperations, categoryID int64) ([]int64, error)
	}

	productDomain struct{}
)

func NewProductDomain() ProductDomain {
	return &productDomain{}
}

func (d *productDomain) CreateProduct(
	ctx context.Context,
	operations db.SQLOperations,
	product *models.Product,
) error {
	product.Touch()
	if product.IsNew() {
		err := operations.QueryRowContext(
			ctx,
			createProductSQL,
			product.Name,
			product.Description,
			product.WholesalePrice,
			product.RetailPrice,
			product.CategoryID,
			product.ProductImage,
			product.ProductType,
			product.Stock,
			product.CreatedAt,
			product.UpdatedAt,
		).Scan(&product.ID)
		if err != nil {
			return apperr.NewDatabaseError(
				err,
			).LogErrorMessage("save product query row context err")
		}

		return nil
	}
	_, err := operations.ExecContext(
		ctx,
		updateProductSQL,
		product.Name,
		product.Description,
		product.WholesalePrice,
		product.RetailPrice,
		product.CategoryID,
		product.ProductImage,
		product.ProductType,
		product.Stock,
		product.UpdatedAt,
		product.ID,
	)
	if err != nil {
		return apperr.NewDatabaseError(
			err,
		).LogErrorMessage("update product query row context err")
	}
	return nil
}

func (d *productDomain) ProductByID(
	ctx context.Context,
	operations db.SQLOperations,
	productID int64,
) (*models.Product, error) {

	row := operations.QueryRowContext(
		ctx,
		getProductByIDSQL,
		productID,
	)

	return d.scanRow(row)
}

func (d *productDomain) DeleteProduct(
	ctx context.Context,
	operations db.SQLOperations,
	productID int64,
) error {

	_, err := operations.ExecContext(ctx, deleteProductSQL, productID)
	if err != nil {
		return apperr.NewDatabaseError(err).LogErrorMessage("delete product error: %v", err)
	}

	return nil
}

func (d *productDomain) ListProducts(
	ctx context.Context,
	operations db.SQLOperations,
	shopID string,
	filter *models.Filter,
) ([]*models.Product, error) {

	filter.ShopID = null.NullValue(shopID)

	query, args := d.buildQuery(getProductsSQL, filter)

	rows, err := operations.QueryContext(
		ctx,
		query,
		args...,
	)
	if err != nil {
		return []*models.Product{}, apperr.NewDatabaseError(
			err,
		).LogErrorMessage(
			"list products query context err: %v",
			err,
		)
	}

	defer rows.Close()

	products := make([]*models.Product, 0)

	for rows.Next() {

		inventory, err := d.scanRow(rows)
		if err != nil {
			return []*models.Product{}, err
		}

		products = append(products, inventory)
	}

	if rows.Err() != nil {
		return []*models.Product{}, apperr.NewDatabaseError(
			rows.Err(),
		).LogErrorMessage(
			"list products rows err: %v",
			err,
		)
	}

	return products, nil
}

func (d *productDomain) ProductCount(
	ctx context.Context,
	operations db.SQLOperations,
	shopID string,
	filter *models.Filter,
) (int, error) {

	filter.ShopID = null.NullValue(shopID)

	query, args := d.buildQuery(getInventoryCountSQL, filter.NoPagination())

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
			"inventory count by query row context err: %v",
			err,
		)
	}

	return count, nil
}

func (d *productDomain) GetAveragePriceByCategory(
	ctx context.Context,
	operations db.SQLOperations,
	categoryID int64,
) (float64, error) {
	
	row := operations.QueryRowContext(
		ctx,
		getAveragePriceByCategorySQL,
		categoryID,
	)

	var averagePrice float64
	err := row.Scan(&averagePrice)
	if err != nil {
		return 0, apperr.NewDatabaseError(
			err,
		).LogErrorMessage("get average price by category err: %v", err)
	}

	return averagePrice, nil
}

func (d *productDomain) GetCategoryHierarchy(
	ctx context.Context,
	operations db.SQLOperations,
	categoryID int64,
) ([]int64, error) {
	
	rows, err := operations.QueryContext(
		ctx,
		getCategoryHierarchySQL,
		categoryID,
	)
	if err != nil {
		return []int64{}, apperr.NewDatabaseError(
			err,
		).LogErrorMessage("get category hierarchy err: %v", err)
	}
	defer rows.Close()

	var categoryIDs []int64
	for rows.Next() {
		var id int64
		err := rows.Scan(&id)
		if err != nil {
			return []int64{}, apperr.NewDatabaseError(
				err,
			).LogErrorMessage("scan category hierarchy row err: %v", err)
		}
		categoryIDs = append(categoryIDs, id)
	}

	if rows.Err() != nil {
		return []int64{}, apperr.NewDatabaseError(
			rows.Err(),
		).LogErrorMessage("category hierarchy rows err: %v", rows.Err())
	}

	return categoryIDs, nil
}

func (d *productDomain) buildQuery(
	query string,
	filter *models.Filter,
) (string, []interface{}) {

	args := make([]interface{}, 0)
	conditions := make([]string, 0)
	counter := utils.NewPlaceholder()

	if filter.CategoryID != nil {
		condition := fmt.Sprintf("category_id = $%d", counter.Touch())
		args = append(args, null.ValueFromNull(filter.CategoryID))
		conditions = append(conditions, condition)
	}

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

func (d *productDomain) scanRow(
	row db.RowScanner,
) (*models.Product, error) {

	var product models.Product

	err := row.Scan(
		&product.ID,
		&product.Name,
		&product.Description,
		&product.WholesalePrice,
		&product.RetailPrice,
		&product.CategoryID,
		&product.ProductImage,
		&product.ProductType,
		&product.Stock,
		&product.CreatedAt,
		&product.UpdatedAt,
	)
	if err != nil {
		return &models.Product{}, apperr.NewDatabaseError(
			err,
		).LogErrorMessage("scan product row err")
	}

	return &product, nil
}