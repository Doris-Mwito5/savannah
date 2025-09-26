package services_test

import (
	"context"
	"database/sql"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github/Doris-Mwito5/savannah-pos/internal/custom_types"
	"github/Doris-Mwito5/savannah-pos/internal/db"
	"github/Doris-Mwito5/savannah-pos/internal/domain"
	"github/Doris-Mwito5/savannah-pos/internal/dtos"
	"github/Doris-Mwito5/savannah-pos/internal/models"
	"github/Doris-Mwito5/savannah-pos/internal/services"
)

type MockCustomerDomain struct {
	mock.Mock
}

func (m *MockCustomerDomain) CustomerByEmailAndPhoneNumber(ctx context.Context, dB db.SQLOperations, email, phone string) (*models.Customer, error) {
	args := m.Called(ctx, dB, email, phone)
	if customer, ok := args.Get(0).(*models.Customer); ok {
		return customer, args.Error(1)
	}
	return nil, args.Error(1)
}
func (m *MockCustomerDomain) CreateCustomer(ctx context.Context, dB db.SQLOperations, customer *models.Customer) error {
	args := m.Called(ctx, dB, customer)
	return args.Error(0)
}
func (m *MockCustomerDomain) CustomerByID(ctx context.Context, dB db.SQLOperations, customerID int64) (*models.Customer, error) {
	args := m.Called(ctx, dB, customerID)
	if customer, ok := args.Get(0).(*models.Customer); ok {
		return customer, args.Error(1)
	}
	return nil, args.Error(1)
}
func (m *MockCustomerDomain) CustomerByEmail(ctx context.Context, dB db.SQLOperations, email string) (*models.Customer, error) {
	args := m.Called(ctx, dB, email)
	if customer, ok := args.Get(0).(*models.Customer); ok {
		return customer, args.Error(1)
	}
	return nil, args.Error(1)
}
func (m *MockCustomerDomain) ListShopCustomers(ctx context.Context, dB db.SQLOperations, shopID string, filter *models.Filter) ([]*models.Customer, error) {
	args := m.Called(ctx, dB, shopID, filter)
	if customers, ok := args.Get(0).([]*models.Customer); ok {
		return customers, args.Error(1)
	}
	return nil, args.Error(1)
}
func (m *MockCustomerDomain) ShopCustomersCount(ctx context.Context, dB db.SQLOperations, shopID string, filter *models.Filter) (int, error) {
	args := m.Called(ctx, dB, shopID, filter)
	return args.Int(0), args.Error(1)
}
func (m *MockCustomerDomain) DeleteCustomer(ctx context.Context, dB db.SQLOperations, customer *models.Customer) error {
	args := m.Called(ctx, dB, customer)
	return args.Error(0)
}

// --- Tests ---

func TestCreateCustomer_Success(t *testing.T) {
	ctx := context.Background()
	mockDomain := new(MockCustomerDomain)
	store := &domain.Store{CustomerDomain: mockDomain}
	service := services.NewCustomerService(store)

	form := &dtos.CreateCustomerForm{
		Name:         "John Doe",
		Email:        "john@example.com",
		PhoneNumber:  "1234567890",
		CustomerType: "REGULAR",
		ShopID:       "shop-1",
	}

	mockDomain.On("CustomerByEmailAndPhoneNumber", ctx, mock.Anything, form.Email, form.PhoneNumber).
		Return(nil, sql.ErrNoRows)
	mockDomain.On("CreateCustomer", ctx, mock.Anything, mock.AnythingOfType("*models.Customer")).
		Return(nil)

	customer, err := service.CreateCustomer(ctx, nil, form)

	assert.NoError(t, err)
	assert.NotNil(t, customer)
	assert.Equal(t, "John Doe", customer.Name)
	assert.Equal(t, "john@example.com", customer.Email)
	assert.Equal(t, "1234567890", customer.PhoneNumber)

	mockDomain.AssertExpectations(t)
}

func TestUpdateCustomer_Success(t *testing.T) {
	ctx := context.Background()
	mockDomain := new(MockCustomerDomain)
	store := &domain.Store{CustomerDomain: mockDomain}
	service := services.NewCustomerService(store)

	customerID := int64(1)
	existing := &models.Customer{
		Name: "Old Name", Email: "old@mail.com", PhoneNumber: "000", CustomerType: "REGULAR",
	}

	form := &dtos.UpdateCustomerForm{
		Name:        ptr("New Name"),
		Email:       ptr("new@mail.com"),
		PhoneNumber: ptr("999"),
		CustomerType: ptr(string(custom_types.CustomerTypeIndividual)),
	}

	mockDomain.On("CustomerByID", ctx, mock.Anything, customerID).
		Return(existing, nil)
	mockDomain.On("CreateCustomer", ctx, mock.Anything, existing).
		Return(nil)

	updated, err := service.UpdateCustomer(ctx, nil, customerID, form)

	assert.NoError(t, err)
	assert.Equal(t, "New Name", updated.Name)
	assert.Equal(t, "new@mail.com", updated.Email)
	assert.Equal(t, "999", updated.PhoneNumber)
	assert.Equal(t, custom_types.CustomerType(custom_types.CustomerTypeIndividual), updated.CustomerType)

	mockDomain.AssertExpectations(t)
}

func TestCustomerByID_Success(t *testing.T) {
	ctx := context.Background()
	mockDomain := new(MockCustomerDomain)
	store := &domain.Store{CustomerDomain: mockDomain}
	service := services.NewCustomerService(store)

	customerID := int64(10)
	expected := &models.Customer{Name: "Jane"}

	mockDomain.On("CustomerByID", ctx, mock.Anything, customerID).
		Return(expected, nil)

	result, err := service.CustomerByID(ctx, nil, customerID)

	assert.NoError(t, err)
	assert.Equal(t, "Jane", result.Name)
	mockDomain.AssertExpectations(t)
}

func TestCustomerByEmail_Success(t *testing.T) {
	ctx := context.Background()
	mockDomain := new(MockCustomerDomain)
	store := &domain.Store{CustomerDomain: mockDomain}
	service := services.NewCustomerService(store)

	email := "test@mail.com"
	expected := &models.Customer{Name: "Tom"}

	mockDomain.On("CustomerByEmail", ctx, mock.Anything, email).
		Return(expected, nil)

	result, err := service.CustomerByEmail(ctx, nil, email)

	assert.NoError(t, err)
	assert.Equal(t, "Tom", result.Name)
	mockDomain.AssertExpectations(t)
}

func TestListShopCustomers_Success(t *testing.T) {
	ctx := context.Background()
	mockDomain := new(MockCustomerDomain)
	store := &domain.Store{CustomerDomain: mockDomain}
	service := services.NewCustomerService(store)

	shopID := "shop-1"
	filter := &models.Filter{Page: 1, Per: 10}
	customers := []*models.Customer{{Name: "Alice"}, {Name: "Bob"}}

	mockDomain.On("ListShopCustomers", ctx, mock.Anything, shopID, filter).
		Return(customers, nil)
	mockDomain.On("ShopCustomersCount", ctx, mock.Anything, shopID, filter).
		Return(2, nil)

	result, err := service.ListShopCustomers(ctx, nil, shopID, filter)

	assert.NoError(t, err)
	assert.Len(t, result.Customers, 2)
	assert.Equal(t, "Alice", result.Customers[0].Name)
	assert.Equal(t, "Bob", result.Customers[1].Name)
	assert.Equal(t, 2, result.Pagination.Count)

	mockDomain.AssertExpectations(t)
}

func TestDeleteCustomer_Success(t *testing.T) {
	ctx := context.Background()
	mockDomain := new(MockCustomerDomain)
	store := &domain.Store{CustomerDomain: mockDomain}
	service := services.NewCustomerService(store)

	customerID := int64(5)
	existing := &models.Customer{Name: "Del Target"}

	mockDomain.On("CustomerByID", ctx, mock.Anything, customerID).
		Return(existing, nil)
	mockDomain.On("DeleteCustomer", ctx, mock.Anything, existing).
		Return(nil)

	result, err := service.DeleteCustomer(ctx, nil, customerID)

	assert.NoError(t, err)
	assert.Equal(t, "Del Target", result.Name)
	mockDomain.AssertExpectations(t)
}

// helper for pointers
func ptr[T any](v T) *T {
	return &v
}
