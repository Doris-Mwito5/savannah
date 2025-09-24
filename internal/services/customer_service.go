package services

import (
	"context"
	"fmt"
	"github/Doris-Mwito5/savannah-pos/internal/apperr"
	"github/Doris-Mwito5/savannah-pos/internal/custom_types"
	"github/Doris-Mwito5/savannah-pos/internal/db"
	"github/Doris-Mwito5/savannah-pos/internal/domain"
	"github/Doris-Mwito5/savannah-pos/internal/dtos"
	"github/Doris-Mwito5/savannah-pos/internal/models"
	"github/Doris-Mwito5/savannah-pos/internal/null"
)

type (
	CustomerService interface {
		CreateCustomer(ctx context.Context, dB db.DB, form *dtos.CreateCustomerForm) (*models.Customer, error)
		UpdateCustomer(ctx context.Context, dB db.DB, customerID int64, form *dtos.UpdateCustomerForm) (*models.Customer, error)
		CustomerByID(ctx context.Context, dB db.DB, customerID int64) (*models.Customer, error)
		CustomerByEmail(ctx context.Context, dB db.DB, Email string) (*models.Customer, error)
		DeleteCustomer(ctx context.Context, dB db.DB, customerID int64) (*models.Customer, error)
		ListShopCustomers(ctx context.Context, dB db.DB, shopID string, filter *models.Filter) (*models.CustomerList, error)
	}

	customerService struct {
		store *domain.Store
	}
)

func NewCustomerService(store *domain.Store) CustomerService {
	return &customerService{
		store: store,
	}
}

func (s *customerService) CreateCustomer(
    ctx context.Context,
    dB db.DB, form *dtos.CreateCustomerForm,
) (*models.Customer, error) {

    _, err := s.store.CustomerDomain.CustomerByEmailAndPhoneNumber(ctx, dB, form.Email, form.PhoneNumber)

    if err == nil {
        return nil, apperr.NewErrorWithType(
            fmt.Errorf("customer with email %s or phone number %s already exists", form.Email, form.PhoneNumber),
            apperr.Conflict,
        )
    }

    // A non-nil error that is NOT a "no rows" error should be returned.
    if !apperr.IsNoRowsErr(err) {
        return nil, err
    }

    // If we've reached this point, the customer does not exist in the database.
    customer := &models.Customer{
        Name:         form.Name,
        Email:        form.Email,
        PhoneNumber:  form.PhoneNumber,
        CustomerType: form.CustomerType,
        ShopID:       form.ShopID,
    }

    err = s.store.CustomerDomain.CreateCustomer(ctx, dB, customer)
    if err != nil {
        return nil, err
    }

    return customer, nil
}

func (s *customerService) UpdateCustomer(
	ctx context.Context,
	dB db.DB,
	customerID int64,
	form *dtos.UpdateCustomerForm,
) (*models.Customer, error) {

	customer, err := s.store.CustomerDomain.CustomerByID(ctx, dB, customerID)
	if err != nil {
		return &models.Customer{}, err
	}

	//if customer is found, update customer details
	if form.Name != nil {
		customer.Name = null.ValueFromNull(form.Name)
	}

	if form.Email != nil {
		customer.Email = null.ValueFromNull(form.Email)
	}

	if form.PhoneNumber != nil {
		customer.PhoneNumber = null.ValueFromNull(form.PhoneNumber)
	}
	if form.CustomerType != nil {
		customer.CustomerType = custom_types.CustomerType(null.ValueFromNull(form.CustomerType))
	}

	err = s.store.CustomerDomain.CreateCustomer(ctx, dB, customer)
	if err != nil {
		return &models.Customer{}, err
	}

	return customer, nil
}

func (s *customerService) CustomerByID(
	ctx context.Context,
	dB db.DB,
	customerID int64,
) (*models.Customer, error) {

	return s.store.CustomerDomain.CustomerByID(ctx, dB, customerID)
}

func (s *customerService) CustomerByEmail(
	ctx context.Context,
	dB db.DB,
	Email string,
) (*models.Customer, error) {

	return s.store.CustomerDomain.CustomerByEmail(ctx, dB, Email)
}

func (s *customerService) ListShopCustomers(
	ctx context.Context,
	dB db.DB,
	shopID string,
	filter *models.Filter,
) (*models.CustomerList, error) {

	customers, err := s.store.CustomerDomain.ListShopCustomers(ctx, dB, shopID, filter)
	if err != nil {
		return &models.CustomerList{}, err
	}

	count, err := s.store.CustomerDomain.ShopCustomersCount(ctx, dB, shopID, filter)
	if err != nil {
		return &models.CustomerList{}, err
	}

	customerList := &models.CustomerList{
		Customers: customers,
		Pagination: models.NewPagination(
			count,
			filter.Page,
			filter.Per,
		),
	}

	return customerList, nil
}

func (s *customerService) DeleteCustomer(
	ctx context.Context,
	dB db.DB,
	customerID int64,
) (*models.Customer, error) {

	customer, err := s.store.CustomerDomain.CustomerByID(ctx, dB, customerID)
	if err != nil {
		return &models.Customer{}, err
	}

	err = s.store.CustomerDomain.DeleteCustomer(ctx, dB, customer)
	if err != nil {
		return &models.Customer{}, err
	}

	return customer, nil
}
