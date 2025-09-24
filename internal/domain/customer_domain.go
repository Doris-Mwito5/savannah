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
	createCustomerSQL                   = "INSERT INTO customers (name, email, phone_number, customer_type, shop_id, created_at, updated_at) VALUES($1, $2, $3, $4, $5, $6, $7) RETURNING(id)"
	getCustomersSQL                     = "SELECT id, name, email, phone_number, customer_type, shop_id, created_at, updated_at FROM customers"
	getCustomerByIDSQL                  = getCustomersSQL + " WHERE id = $1"
	getCustomerByEmailSQL               = getCustomersSQL + " WHERE email = $1"
	getCustomerByEmailAndPhoneNumberSQL = getCustomersSQL + " WHERE email = $1 AND phone_number = $2"
	getCustomersCountSQL                = "SELECT COUNT(id) FROM customers"
	updateCustomeSQL                    = "UPDATE customers SET name = $1, email = $2, phone_number = $3, customer_type = $4, shop_id = $5, updated_at = $6 WHERE id = $7"
	deleteCustomerSQL                   = "DELETE FROM customers WHERE id = $1"
)

type (
	CustomerDomain interface {
		CreateCustomer(ctx context.Context, operations db.SQLOperations, customer *models.Customer) error
		CustomerByID(ctx context.Context, operations db.SQLOperations, customerID int64) (*models.Customer, error)
		CustomerByEmail(ctx context.Context, operations db.SQLOperations, Email string) (*models.Customer, error)
		CustomerByEmailAndPhoneNumber(ctx context.Context, operations db.SQLOperations, email, phoneNumber string) (*models.Customer, error)
		ListShopCustomers(ctx context.Context, operations db.SQLOperations, shopID string, filter *models.Filter) ([]*models.Customer, error)
		ShopCustomersCount(ctx context.Context, operations db.SQLOperations, shopID string, filter *models.Filter) (int, error)
		DeleteCustomer(ctx context.Context, operations db.SQLOperations, customer *models.Customer) error
	}

	customerDomain struct{}
)

func NewCustomerDomain() CustomerDomain {
	return &customerDomain{}
}

func (d *customerDomain) CreateCustomer(
	ctx context.Context,
	operations db.SQLOperations,
	customer *models.Customer,
) error {

	customer.Touch()
	if customer.IsNew() {
		err := operations.QueryRowContext(
			ctx,
			createCustomerSQL,
			customer.Name,
			customer.Email,
			customer.PhoneNumber,
			customer.CustomerType,
			customer.ShopID,
			customer.CreatedAt,
			customer.UpdatedAt,
		).Scan(&customer.ID)
		if err != nil {
			return apperr.NewDatabaseError(
				err,
			).LogErrorMessage("save customer query err: %v", err)
		}
		return nil
	}
	_, err := operations.ExecContext(
		ctx,
		updateCustomeSQL,
		customer.Name,
		customer.Email,
		customer.PhoneNumber,
		customer.CustomerType,
		customer.ShopID,
		customer.UpdatedAt,
		customer.ID,
	)
	if err != nil {
		return apperr.NewDatabaseError(
			err,
		).LogErrorMessage("update customer query err: %v", err)
	}
	return nil
}

func (d *customerDomain) CustomerByID(
	ctx context.Context,
	operations db.SQLOperations,
	customerID int64,
) (*models.Customer, error) {

	row := operations.QueryRowContext(
		ctx,
		getCustomerByIDSQL,
		customerID,
	)

	return d.scanRow(row)
}

func (d *customerDomain) CustomerByEmail(
	ctx context.Context,
	operations db.SQLOperations,
	Email string,
) (*models.Customer, error) {

	row := operations.QueryRowContext(
		ctx,
		getCustomerByEmailSQL,
		Email,
	)

	return d.scanRow(row)
}

func (d *customerDomain) CustomerByEmailAndPhoneNumber(
	ctx context.Context,
	operations db.SQLOperations,
	email,
	phoneNumber string,
) (*models.Customer, error) {

	row := operations.QueryRowContext(
		ctx,
		getCustomerByEmailAndPhoneNumberSQL,
		email,
		phoneNumber,
	)

	return d.scanRow(row)
}

func (d *customerDomain) ListShopCustomers(
	ctx context.Context,
	operations db.SQLOperations,
	shopID string,
	filter *models.Filter,
) ([]*models.Customer, error) {

	filter.ShopID = null.NullValue(shopID)

	query, args := d.buildquery(getCustomersSQL, filter)

	rows, err := operations.QueryContext(
		ctx,
		query,
		args...,
	)

	if err != nil {
		return []*models.Customer{}, apperr.NewDatabaseError(
			err,
		).LogErrorMessage(
			"list customers query context err: %v",
			err,
		)
	}

	defer rows.Close()

	customers := make([]*models.Customer, 0)

	for rows.Next() {
		customer, err := d.scanRow(rows)
		if err != nil {
			return []*models.Customer{}, err
		}

		customers = append(customers, customer)
	}

	return customers, nil
}

func (d *customerDomain) ShopCustomersCount(
	ctx context.Context,
	operations db.SQLOperations,
	shopID string,
	filter *models.Filter,
) (int, error) {
	filter.ShopID = null.NullValue(shopID)

	query, args := d.buildquery(getCustomersCountSQL, filter.NoPagination())

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
			"customer count by query row context err: %v",
			err,
		)
	}

	return count, nil
}

func (d *customerDomain) DeleteCustomer(
	ctx context.Context,
	operations db.SQLOperations,
	customer *models.Customer,
) error {
	_, err := operations.ExecContext(
		ctx,
		deleteCustomerSQL,
		customer.ID,
	)

	if err != nil {
		return apperr.NewDatabaseError(
			err,
		).LogErrorMessage("delete customer err: %v", err)
	}

	return nil
}

func (d *customerDomain) buildquery(
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

	if filter.Term != "" {
		textCols := []string{"name", "description"}
		likeStatements := make([]string, 0)
		term := strings.ToLower(filter.Term)

		for _, col := range textCols {
			likeStmt := fmt.Sprintf(
				" (LOWER(%s) LIKE '%%' || $%d || '%%') ", col, counter.Touch())
			likeStatements = append(likeStatements, likeStmt)
			args = append(args, term)
		}

		conditions = append(conditions, " ("+strings.Join(likeStatements, " OR ")+")")
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

func (d *customerDomain) scanRow(
	row db.RowScanner,
) (*models.Customer, error) {

	var customer models.Customer

	err := row.Scan(
		&customer.ID,
		&customer.Name,
		&customer.Email,
		&customer.PhoneNumber,
		&customer.CustomerType,
		&customer.ShopID,
		&customer.CreatedAt,
		&customer.UpdatedAt,
	)
	if err != nil {
		return &models.Customer{}, apperr.NewDatabaseError(
			err,
		).LogErrorMessage("scan customer row err: %v", err)
	}

	return &customer, nil
}
