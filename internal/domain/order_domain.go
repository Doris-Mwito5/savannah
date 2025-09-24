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
	createOrderSQL    = "INSERT INTO orders (reference_number, phone_number, order_status, order_source, payment_method, customer_id, shop_id, total_items, total_amount, discount, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12) RETURNING(id)"
	getOrdersSQL      = "SELECT id, reference_number, phone_number, order_status, order_source, payment_method, customer_id, shop_id, total_items, total_amount, discount, created_at, updated_at FROM orders"
	getOrderByIDSQL   = getOrdersSQL + " WHERE id = $1"
	getOrdersCountSQL = "SELECT COUNT(id) FROM orders"
	deleteOrdersSQL   = "DELETE FROM orders WHERE id = $1"
	updateOrderSQL    = "UPDATE orders SET reference_number = $1, phone_number = $2, order_status = $3, order_source = $4, payment_method = $5, customer_id = $6, shop_id = $7, total_items = $8, total_amount = $9, discount = $10, updated_at = $11 WHERE id = $12"
)

type (
	OrderDomain interface {
		CreateOrder(ctx context.Context, operations db.SQLOperations, order *models.Order) error
		OrderByID(ctx context.Context, operations db.SQLOperations, orderID int64) (*models.Order, error)
		LisOrders(ctx context.Context, operations db.SQLOperations, shopID string, filter *models.Filter) ([]*models.Order, error)
		OrderCount(ctx context.Context, operations db.SQLOperations, shopID string, filter *models.Filter) (int, error)
	}

	orderDomain struct{}
)

func NewOrderDomain() OrderDomain {
	return &orderDomain{}
}

func (d *orderDomain) CreateOrder(
	ctx context.Context,
	operations db.SQLOperations,
	order *models.Order,
) error {

	order.Touch()
	if order.IsNew() {
		err := operations.QueryRowContext(
			ctx,
			createOrderSQL,
			order.ReferenceNumber,
			order.PhoneNumber,
			order.OrderStatus,
			order.OrderMedium,
			order.PaymentMethod,
			order.CustomerID,
			order.ShopID,
			order.TotalItems,
			order.TotalAmount,
			order.Discount,
			order.CreatedAt,
			order.UpdatedAt,
		).Scan(&order.ID)
		if err != nil {
			return apperr.NewDatabaseError(
				err,
			).LogErrorMessage("save order query row err: %v", err)
		}
		return nil
	}

	_, err := operations.ExecContext(
		ctx,
		updateOrderSQL,
		order.ReferenceNumber,
		order.PhoneNumber,
		order.OrderStatus,
		order.OrderMedium,
		order.PaymentMethod,
		order.CustomerID,
		order.ShopID,
		order.TotalItems,
		order.TotalAmount,
		order.Discount,
		order.UpdatedAt,
		order.ID,
	)
	if err != nil {
		return apperr.NewDatabaseError(
			err,
		).LogErrorMessage("update order query row err: %v", err)
	}
	return nil
}

func (d *orderDomain) OrderByID(
	ctx context.Context,
	operations db.SQLOperations,
	orderID int64,
) (*models.Order, error) {

	row := operations.QueryRowContext(
		ctx,
		getOrderByIDSQL,
		orderID,
	)

	return d.scanRow(row)
}

func (d *orderDomain) OrderCount(
	ctx context.Context,
	operations db.SQLOperations,
	shopID string,
	filter *models.Filter,
) (int, error) {
	filter.ShopID = null.NullValue(shopID)
	query, args := d.buildQuery(getOrdersCountSQL, filter.NoPagination())

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
			"order count by query row contex err: %v",
			err,
		)
	}
	return count, nil
}

func (d *orderDomain) LisOrders(
	ctx context.Context,
	operations db.SQLOperations,
	shopID string,
	filter *models.Filter,
) ([]*models.Order, error) {

	filter.ShopID = null.NullValue(shopID)
	query, args := d.buildQuery(getOrdersSQL, filter)
	rows, err := operations.QueryContext(
		ctx,
		query,
		args...,
	)
	if err != nil {
		return []*models.Order{}, apperr.NewDatabaseError(
			err,
		).LogErrorMessage("order query row err: %v", err)
	}

	defer rows.Close()

	orders := make([]*models.Order, 0)

	for rows.Next() {
		order, err := d.scanRow(rows)
		if err != nil {
			return []*models.Order{}, err
		}

		orders = append(orders, order)
	}

	return orders, nil
}

func (d *orderDomain) buildQuery(
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

func (d *orderDomain) scanRow(
	row db.RowScanner,
) (*models.Order, error) {
	var order models.Order

	err := row.Scan(
		&order.ID,
		&order.ReferenceNumber,
		&order.PhoneNumber, 
		&order.OrderStatus,
		&order.OrderMedium,
		&order.PaymentMethod,
		&order.CustomerID,
		&order.ShopID,
		&order.TotalItems,
		&order.TotalAmount,
		&order.Discount,
		&order.CreatedAt,
		&order.UpdatedAt,
	)
	if err != nil {
		return nil, apperr.NewDatabaseError(
			err,
		).LogErrorMessage("scan order row err: %v", err)
	}

	return &order, nil
}

