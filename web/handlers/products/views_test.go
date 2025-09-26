package products

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github/Doris-Mwito5/savannah-pos/internal/custom_types"
	"github/Doris-Mwito5/savannah-pos/internal/db"
	"github/Doris-Mwito5/savannah-pos/internal/dtos"
	"github/Doris-Mwito5/savannah-pos/internal/models"
)


type MockProductService struct {
	mock.Mock
}

func (m *MockProductService) CreateProduct(ctx context.Context, dB db.DB, form *dtos.CreateProductForm) (*models.Product, error) {
	args := m.Called(ctx, dB, form)
	if prod, ok := args.Get(0).(*models.Product); ok {
		return prod, args.Error(1)
	}
	return nil, args.Error(1)
}
func (m *MockProductService) UpdateProduct(ctx context.Context, dB db.DB, id int64, form *dtos.UpdateProductForm) (*models.Product, error) {
	args := m.Called(ctx, dB, id, form)
	if prod, ok := args.Get(0).(*models.Product); ok {
		return prod, args.Error(1)
	}
	return nil, args.Error(1)
}
func (m *MockProductService) ProductByID(ctx context.Context, dB db.DB, id int64) (*models.Product, error) {
	args := m.Called(ctx, dB, id)
	if prod, ok := args.Get(0).(*models.Product); ok {
		return prod, args.Error(1)
	}
	return nil, args.Error(1)
}
func (m *MockProductService) ListProducts(ctx context.Context, dB db.DB, shopID string, filter *models.Filter) (*models.ProductList, error) {
	args := m.Called(ctx, dB, shopID, filter)
	if list, ok := args.Get(0).(*models.ProductList); ok {
		return list, args.Error(1)
	}
	return nil, args.Error(1)
}
func (m *MockProductService) DeleteProduct(ctx context.Context, dB db.DB, id int64) (*models.Product, error) {
	args := m.Called(ctx, dB, id)
	if prod, ok := args.Get(0).(*models.Product); ok {
		return prod, args.Error(1)
	}
	return nil, args.Error(1)
}
func (m *MockProductService) GetAveragePriceByCategory(ctx context.Context, dB db.DB, categoryID int64) (float64, error) {
	args := m.Called(ctx, dB, categoryID)
	return args.Get(0).(float64), args.Error(1)
}

// --- Tests ---

func TestCreateProduct_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockSvc := new(MockProductService)
	router := gin.New()
	router.POST("/products", createProduct(nil, mockSvc))

	form := dtos.CreateProductForm{Name: "Product A"}
	expected := &models.Product{SequentialIdentifier: custom_types.SequentialIdentifier{ID: 1}, Name: "Product A"}

	mockSvc.On("CreateProduct", mock.Anything, nil, &form).Return(expected, nil)

	body, _ := json.Marshal(form)
	req := httptest.NewRequest(http.MethodPost, "/products", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)
	var got models.Product
	_ = json.Unmarshal(w.Body.Bytes(), &got)
	assert.Equal(t, expected.Name, got.Name)

	mockSvc.AssertExpectations(t)
}

func TestUpdateProduct_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockSvc := new(MockProductService)
	router := gin.New()
	router.PUT("/products/:id", updateProduct(nil, mockSvc))

	form := dtos.UpdateProductForm{Name: ptr("Updated")}
	expected := &models.Product{SequentialIdentifier: custom_types.SequentialIdentifier{ID: 1}, Name: "Updated"}

	mockSvc.On("UpdateProduct", mock.Anything, nil, int64(1), &form).Return(expected, nil)

	body, _ := json.Marshal(form)
	req := httptest.NewRequest(http.MethodPut, "/products/1", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	var got models.Product
	_ = json.Unmarshal(w.Body.Bytes(), &got)
	assert.Equal(t, "Updated", got.Name)

	mockSvc.AssertExpectations(t)
}

func TestGetProduct_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockSvc := new(MockProductService)
	router := gin.New()
	router.GET("/products/:id", getProduct(nil, mockSvc))

	expected := &models.Product{SequentialIdentifier: custom_types.SequentialIdentifier{ID: 1}, Name: "Prod X"}
	mockSvc.On("ProductByID", mock.Anything, nil, int64(2)).Return(expected, nil)

	req := httptest.NewRequest(http.MethodGet, "/products/2", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	var got models.Product
	_ = json.Unmarshal(w.Body.Bytes(), &got)
	assert.Equal(t, expected.Name, got.Name)

	mockSvc.AssertExpectations(t)
}

func TestDeleteProduct_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockSvc := new(MockProductService)
	router := gin.New()
	router.DELETE("/products/:id", deleteProduct(nil, mockSvc))

	expected := &models.Product{SequentialIdentifier: custom_types.SequentialIdentifier{ID: 1}, Name: "Deleted"}
	mockSvc.On("DeleteProduct", mock.Anything, nil, int64(3)).Return(expected, nil)

	req := httptest.NewRequest(http.MethodDelete, "/products/3", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	var got models.Product
	_ = json.Unmarshal(w.Body.Bytes(), &got)
	assert.Equal(t, expected.Name, got.Name)

	mockSvc.AssertExpectations(t)
}

func TestGetAveragePriceByCategory_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockSvc := new(MockProductService)
	router := gin.New()
	router.GET("/products/category/:category_id/average", getAveragePriceByCategory(nil, mockSvc))

	mockSvc.On("GetAveragePriceByCategory", mock.Anything, nil, int64(5)).Return(99.5, nil)

	req := httptest.NewRequest(http.MethodGet, "/products/category/5/average", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	var got float64
	_ = json.Unmarshal(w.Body.Bytes(), &got)
	assert.Equal(t, 99.5, got)

	mockSvc.AssertExpectations(t)
}

func ptr(s string) *string { return &s }
