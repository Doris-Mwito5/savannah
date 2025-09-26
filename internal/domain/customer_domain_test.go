package domain

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"testing"
	"time"

	"github/Doris-Mwito5/savannah-pos/internal/custom_types"
	"github/Doris-Mwito5/savannah-pos/internal/loggers"
	"github/Doris-Mwito5/savannah-pos/internal/models"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

type dbWrapper struct {
	DB *sql.DB
}

func (w dbWrapper) ValidForPostgres() bool { return true }
func (w dbWrapper) QueryRowContext(ctx context.Context, query string, args ...interface{}) *sql.Row {
	return w.DB.QueryRowContext(ctx, query, args...)
}
func (w dbWrapper) ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error) {
	return w.DB.ExecContext(ctx, query, args...)
}
func (w dbWrapper) QueryContext(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error) {
	return w.DB.QueryContext(ctx, query, args...)
}

func TestMain(m *testing.M) {
    loggers.InitLogger("test")
    
    code := m.Run()
    os.Exit(code)
}

func TestCustomerDomain_CreateCustomer(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	customerDomain := NewCustomerDomain()
	ctx := context.Background()

	customer := &models.Customer{
		Name:         "Alice",
		Email:        "alice@example.com",
		PhoneNumber:  "+254700000000",
		CustomerType: custom_types.CustomerType("individual"),
		ShopID:       "shop123",
	}

	// -------- New Customer (INSERT) --------
	mock.ExpectQuery("INSERT INTO customers").
		WithArgs(customer.Name, customer.Email, customer.PhoneNumber, customer.CustomerType,
			customer.ShopID, sqlmock.AnyArg(), sqlmock.AnyArg()).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))

	err = customerDomain.CreateCustomer(ctx, dbWrapper{DB: db}, customer)
	assert.NoError(t, err)
	assert.Equal(t, int64(1), customer.ID)

	// -------- Update Customer (UPDATE) --------
	customer.Name = "Alice Updated"
	customer.Touch()
	mock.ExpectExec("UPDATE customers SET .* WHERE id = \\$7").
		WithArgs(customer.Name, customer.Email, customer.PhoneNumber, customer.CustomerType,
			customer.ShopID, sqlmock.AnyArg(), customer.ID).
		WillReturnResult(sqlmock.NewResult(1, 1))

	err = customerDomain.CreateCustomer(ctx, dbWrapper{DB: db}, customer)
	assert.NoError(t, err)

	// -------- Error on INSERT --------
	newCust := &models.Customer{
		Name:         "Bob",
		Email:        "bob@example.com",
		CustomerType: custom_types.CustomerType("individual"),
		ShopID:       "shop123",
	}
	mock.ExpectQuery("INSERT INTO customers").
		WillReturnError(fmt.Errorf("insert failed"))

	err = customerDomain.CreateCustomer(ctx, dbWrapper{DB: db}, newCust)
	assert.Error(t, err)

	// -------- Error on UPDATE --------
	customer.ID = 99
	mock.ExpectExec("UPDATE customers SET .* WHERE id = \\$7").
		WillReturnError(fmt.Errorf("update failed"))

	err = customerDomain.CreateCustomer(ctx, dbWrapper{DB: db}, customer)
	assert.Error(t, err)

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestCustomerDomain_CustomerByID(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()
	customerDomain := NewCustomerDomain()
	ctx := context.Background()
	now := time.Now()

	mock.ExpectQuery("SELECT .* FROM customers WHERE id = \\$1").
		WithArgs(1).
		WillReturnRows(sqlmock.NewRows([]string{
			"id", "name", "email", "phone_number", "customer_type", "shop_id", "created_at", "updated_at",
		}).AddRow(1, "Alice", "alice@example.com", "+254700000000", []byte("individual"), "shop123", now, now))

	cust, err := customerDomain.CustomerByID(ctx, dbWrapper{DB: db}, 1)
	assert.NoError(t, err)
	assert.Equal(t, "Alice", cust.Name)
}

func TestCustomerDomain_CustomerByEmail(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()
	customerDomain := NewCustomerDomain()
	ctx := context.Background()
	now := time.Now()

	mock.ExpectQuery("SELECT .* FROM customers WHERE email = \\$1").
		WithArgs("alice@example.com").
		WillReturnRows(sqlmock.NewRows([]string{
			"id", "name", "email", "phone_number", "customer_type", "shop_id", "created_at", "updated_at",
		}).AddRow(1, "Alice", "alice@example.com", "+254700000000", []byte("individual"), "shop123", now, now))

	cust, err := customerDomain.CustomerByEmail(ctx, dbWrapper{DB: db}, "alice@example.com")
	assert.NoError(t, err)
	assert.Equal(t, "Alice", cust.Name)
}

func TestCustomerDomain_CustomerByEmailAndPhoneNumber(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()
	customerDomain := NewCustomerDomain()
	ctx := context.Background()
	now := time.Now()

	mock.ExpectQuery("SELECT .* FROM customers WHERE email = \\$1 AND phone_number = \\$2").
		WithArgs("alice@example.com", "+254700000000").
		WillReturnRows(sqlmock.NewRows([]string{
			"id", "name", "email", "phone_number", "customer_type", "shop_id", "created_at", "updated_at",
		}).AddRow(1, "Alice", "alice@example.com", "+254700000000", []byte("individual"), "shop123", now, now))

	cust, err := customerDomain.CustomerByEmailAndPhoneNumber(ctx, dbWrapper{DB: db}, "alice@example.com", "+254700000000")
	assert.NoError(t, err)
	assert.Equal(t, "Alice", cust.Name)
}

func TestCustomerDomain_ListShopCustomers(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()
	customerDomain := NewCustomerDomain()
	ctx := context.Background()
	now := time.Now()

	mock.ExpectQuery("SELECT .* FROM customers").
		WillReturnRows(sqlmock.NewRows([]string{
			"id", "name", "email", "phone_number", "customer_type", "shop_id", "created_at", "updated_at",
		}).AddRow(1, "Alice", "alice@example.com", "+254700000000", []byte("individual"), "shop123", now, now))

	customers, err := customerDomain.ListShopCustomers(ctx, dbWrapper{DB: db}, "shop123", &models.Filter{})
	assert.NoError(t, err)
	assert.Len(t, customers, 1)
	assert.Equal(t, "Alice", customers[0].Name)
}

func TestCustomerDomain_ShopCustomersCount(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()
	customerDomain := NewCustomerDomain()
	ctx := context.Background()

	mock.ExpectQuery("SELECT COUNT\\(id\\) FROM customers").
		WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(5))

	count, err := customerDomain.ShopCustomersCount(ctx, dbWrapper{DB: db}, "shop123", &models.Filter{})
	assert.NoError(t, err)
	assert.Equal(t, 5, count)
}

func TestCustomerDomain_DeleteCustomer(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()
	customerDomain := NewCustomerDomain()
	ctx := context.Background()

	cust := &models.Customer{SequentialIdentifier: custom_types.SequentialIdentifier{ID: 1}}
	mock.ExpectExec("DELETE FROM customers WHERE id = \\$1").
		WithArgs(cust.ID).
		WillReturnResult(sqlmock.NewResult(0, 1))

	err := customerDomain.DeleteCustomer(ctx, dbWrapper{DB: db}, cust)
	assert.NoError(t, err)
}
