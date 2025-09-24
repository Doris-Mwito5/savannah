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
	createOrderItemSQL       = "INSERT INTO order_items (order_id, product_id, unit_price, quantity, total_amount, created_at, updated_at) VALUES "
	getOrderItemsSQL         = "SELECT oi.id, oi.order_id, p.product_id,  oi.unit_price, oi.quantity, oi.total_amount, oi.created_at, oi.updated_at FROM order_items oi INNER JOIN products p ON oi.product_id = p.id"
	getOrderItemByIDSQL      = getOrderItemsSQL + " WHERE id= $1"
	getOrderItemsCountSQL    = "SELECT COUNT(id) FROM order_items"
	deleteOrderItemSQL       = "DELETE FROM order_items WHERE product_id = $1"
)

type (
	OrderItemDomain interface {
		InsertOrderItems(ctx context.Context, operations db.SQLOperations, orderItems []*models.OrderItem) error
		OrderItemByID(ctx context.Context, operations db.SQLOperations, orderItemID int64) (*models.OrderItem, error)
		OrderItemCount(ctx context.Context, operations db.SQLOperations, orderID int64, filter *models.Filter) (int, error)
		OrderItems(ctx context.Context, operations db.SQLOperations, orderID int64, filter *models.Filter) ([]*models.OrderItem, error)
		DeleteOrderItems(ctx context.Context, operations db.SQLOperations, productID int64) error
	}

	orderItemDomain struct{}
)

func NewOrderItemDomain() OrderItemDomain {
	return &orderItemDomain{}
}

func (d *orderItemDomain) InsertOrderItems(
    ctx context.Context,
    operations db.SQLOperations,
    orderItems []*models.OrderItem,
) error {

    // Base query with all columns, including `updated_at`
    baseSQL := "INSERT INTO order_items (order_id, product_id, unit_price, quantity, total_amount, created_at, updated_at) VALUES "
    
    counter := utils.NewPlaceholder()
    placeholders := make([]string, len(orderItems))
    values := make([]any, 0)

    for index, orderItem := range orderItems {
        orderItem.Touch()

        // Create a placeholder string for each item with 7 values, e.g., "($1, $2, $3, $4, $5, $6, $7)"
        placeholder := make([]string, 7)
        for i := 0; i < 7; i++ {
            placeholder[i] = fmt.Sprintf("$%d", counter.Touch())
        }

        placeholders[index] = "(" + strings.Join(placeholder, ",") + ")"
        
        // Append the 7 values for the current order item
        values = append(values,
            orderItem.OrderID,
            orderItem.ProductID,
            orderItem.UnitPrice,
            orderItem.Quantity,
            orderItem.TotalAmount,
            orderItem.CreatedAt,
			orderItem.UpdatedAt,
        )
    }

    // Join the placeholders and form the final query.
    query := baseSQL + strings.Join(placeholders, ",")

    _, err := operations.ExecContext(ctx, query, values...)
    if err != nil {
        return apperr.NewDatabaseError(err).LogErrorMessage("insert order items exec error")
    }

    return nil
}

func (d *orderItemDomain) OrderItemByID(
	ctx context.Context,
	operations db.SQLOperations,
	orderItemID int64,
) (*models.OrderItem, error) {

	row := operations.QueryRowContext(
		ctx,
		getOrderItemByIDSQL,
		orderItemID,
	)

	return d.scanRow(row)
}

func (d *orderItemDomain) OrderItemCount(
	ctx context.Context,
	operations db.SQLOperations,
	orderID int64,
	filter *models.Filter,
) (int, error) {

	filter.OrderID = null.NullValue(orderID)

	query, args := d.buildQuery(getOrderItemsCountSQL, filter.NoPagination())

	rows := operations.QueryRowContext(
		ctx,
		query,
		args...,
	)

	var count int
	err := rows.Scan(&count)
	if err != nil {
		return 0, apperr.NewDatabaseError(err).LogErrorMessage("order item count query row context err")
	}

	return count, nil
}

func (d *orderItemDomain) OrderItems(
	ctx context.Context,
	operations db.SQLOperations,
	orderID int64,
	filter *models.Filter,
) ([]*models.OrderItem, error) {

	filter.OrderID = null.NullValue(orderID)

	query, args := d.buildQuery(getOrderItemsSQL, filter)

	rows, err := operations.QueryContext(
		ctx,
		query,
		args...,
	)
	if err != nil {
		return []*models.OrderItem{}, apperr.NewDatabaseError(err).LogErrorMessage("order items query context err")
	}
	defer rows.Close()

	orderItems := make([]*models.OrderItem, 0)
	for rows.Next() {
		orderItem, err := d.scanRow(rows)
		if err != nil {
			return []*models.OrderItem{}, err
		}
		orderItems = append(orderItems, orderItem)
	}
	if rows.Err() != nil {
		return []*models.OrderItem{}, apperr.NewDatabaseError(rows.Err()).LogErrorMessage("list order items err")
	}

	return orderItems, nil
}

func (d *orderItemDomain) DeleteOrderItems(
	ctx context.Context,
	operations db.SQLOperations,
	productID int64,
) error {

	_, err := operations.ExecContext(ctx, deleteOrderItemSQL, productID)
	if err != nil {
		return apperr.NewDatabaseError(err).LogErrorMessage("delete order items error: %v", err)
	}

	return nil
}

func (d *orderItemDomain) buildQuery(
	query string,
	filter *models.Filter,
) (string, []interface{}) {
	args := make([]interface{}, 0)
	conditions := make([]string, 0)
	counter := utils.NewPlaceholder()

	if filter.OrderID != nil {
		condition := fmt.Sprintf("oi.order_id = $%d", counter.Touch())
		args = append(args, null.ValueFromNull(filter.OrderID))
		conditions = append(conditions, condition)
	}

	if len(conditions) > 0 {
		query += " WHERE " + strings.Join(conditions, " AND ")
	}

	if filter.Page > 0 && filter.Per > 0 {
		query += fmt.Sprintf(" ORDER BY oi.created_at DESC LIMIT $%d OFFSET $%d", counter.Touch(), counter.Touch())
		args = append(args, filter.Per, (filter.Page-1)*filter.Per)
	}

	return query, args
}

func (d *orderItemDomain) scanRow(
	row db.RowScanner,
) (*models.OrderItem, error) {
	var orderItem models.OrderItem

	err := row.Scan(
		&orderItem.ID,
		&orderItem.OrderID,
		&orderItem.ProductID,
		&orderItem.UnitPrice,
		&orderItem.Quantity,
		&orderItem.TotalAmount,
		&orderItem.CreatedAt,
	)
	if err != nil {
		return &models.OrderItem{}, apperr.NewDatabaseError(err).LogErrorMessage("scan row err")
	}

	return &orderItem, nil
}
