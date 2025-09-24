package services

import (
	"context"
	"github/Doris-Mwito5/savannah-pos/internal/apperr"
	"github/Doris-Mwito5/savannah-pos/internal/custom_types"
	"github/Doris-Mwito5/savannah-pos/internal/db"
	"github/Doris-Mwito5/savannah-pos/internal/domain"
	"github/Doris-Mwito5/savannah-pos/internal/dtos"
	"github/Doris-Mwito5/savannah-pos/internal/loggers"
	"github/Doris-Mwito5/savannah-pos/internal/models"
	"github/Doris-Mwito5/savannah-pos/internal/notification"
	"github/Doris-Mwito5/savannah-pos/internal/null"
	"github/Doris-Mwito5/savannah-pos/internal/utils"
)

type (
	OrderService interface {
		CreateOrder(ctx context.Context, dB db.DB, form *dtos.CreateOrderForm) (*models.Order, error)
		ListShopOrders(ctx context.Context, dB db.DB, shopID string, filter *models.Filter) (*models.OrderList, error)
		OrderByID(ctx context.Context, dB db.DB, orderID int64) (*models.Order, error)
	}

	orderService struct {
		customerService   CustomerService
		store             *domain.Store
		orderNotification notification.OrderNotification
	}
)

func NewOrderService(
	customerService CustomerService,
	store *domain.Store,
	orderNotification notification.OrderNotification,
) OrderService {
	return &orderService{
		customerService:   customerService,
		store:             store,
		orderNotification: orderNotification,
	}
}

func (s *orderService) CreateOrder(
    ctx context.Context,
    dB db.DB,
    form *dtos.CreateOrderForm,
) (*models.Order, error) {

    // 1. Initialize the order object with basic data from the form.
    order := &models.Order{
        ReferenceNumber: utils.GenerateTransactionRef(),
        OrderStatus:     custom_types.OrderStatus(form.OrderStatus),
        OrderMedium:     custom_types.OrderMedium(form.OrderMedium),
        PaymentMethod:   custom_types.PaymentMethod(form.PaymentMethod),
        TotalItems:      len(form.Items),
        ShopID:          form.ShopID,
        PhoneNumber:     form.PhoneNumber,
    }

    // These variables will be populated inside the transaction.
    var orderItems []*models.OrderItem
    var totalOrderCost float64

    // 2. Start a database transaction to ensure all operations (customer, order, items) succeed or fail together.
    err := dB.InTransaction(ctx, func(ctx context.Context, operations db.SQLOperations) error {
        
        // 3. Check for an existing customer by email and phone number using the transaction's operations.
        customer, err := s.store.CustomerDomain.CustomerByEmailAndPhoneNumber(ctx, operations, form.CustomerEmail, form.PhoneNumber)
        
        // 4. Handle the customer creation based on the result of the check.
        if err == nil {
            // Customer exists, link the order to the existing customer ID.
            order.CustomerID = null.NullValue(customer.ID)
        } else if apperr.IsNoRowsErr(err) {
            // Customer does not exist, so we create a new one.
            formCustomer := &dtos.CreateCustomerForm{
                Name:         form.CustomerName,
                Email:        form.CustomerEmail,
                PhoneNumber:  form.PhoneNumber,
                ShopID:       form.ShopID,
                CustomerType: "individual",
            }
            
            createdCustomer, err := s.customerService.CreateCustomer(ctx, dB, formCustomer)
            if err != nil {
                loggers.Errorf("failed to create new customer: [%+v]", err)
                return err
            }
            
            order.CustomerID = null.NullValue(createdCustomer.ID)
        } else {
            // An unexpected database error occurred.
            loggers.Errorf("database error when checking for customer: [%+v]", err)
            return err
        }
    
        // 5. Calculate order items and total cost.
        orderItems, totalOrderCost, err = s.getPriceAndOrderItems(ctx, operations, form)
        if err != nil {
            loggers.Errorf("failed to get price and order items: [%+v]", err)
            return err
        }
    
        // 6. Finalize order details.
        order.TotalAmount = totalOrderCost
        if form.Discount != nil {
            order.Discount = null.NullValue(totalOrderCost - null.ValueFromNull(form.Discount))
        }
    
        // 7. Save the order to the database.
        err = s.store.OrderDomain.CreateOrder(ctx, operations, order)
        if err != nil {
            loggers.Errorf("failed to save order: [%+v]", err)
            return err
        }
    
        // 8. Save the order's items in batches.
        if len(orderItems) > 0 {
            for _, orderItem := range orderItems {
                orderItem.OrderID = order.ID
            }
            orderItemBatches := s.createOrderItemsBatches(orderItems)
            for _, batch := range orderItemBatches {
                err = s.store.OrderItemDomain.InsertOrderItems(ctx, operations, batch)
                if err != nil {
                    loggers.Errorf("failed to insert bulk order items: [%+v]", err)
                    return err
                }
            }
        }
    
        return nil
    })

    if err != nil {
        // If the transaction failed, the entire operation is rolled back.
        return nil, err
    }

    // 9. After the successful transaction, send the SMS notification.
    if order.PhoneNumber != "" {
        err = s.orderNotification.SendOrderNotifications(order)
        if err != nil {
            loggers.Errorf("failed to send order confirmation SMS: [%+v]", err)
            // Note: We log the error but don't return it, as the order itself was successfully created.
        }
    }

    return order, nil
}

func (s *orderService) ListShopOrders(
	ctx context.Context,
	dB db.DB,
	shopID string,
	filter *models.Filter,
) (*models.OrderList, error) {

	orders, err := s.store.OrderDomain.LisOrders(ctx, dB, shopID, filter)
	if err != nil {
		return &models.OrderList{}, err
	}

	count, err := s.store.OrderDomain.OrderCount(ctx, dB, shopID, filter)
	if err != nil {
		return &models.OrderList{}, err
	}

	orderList := &models.OrderList{
		Orders: orders,
		Pagination: models.NewPagination(
			count,
			filter.Page,
			filter.Per,
		),
	}

	return orderList, nil
}

func (s *orderService) OrderByID(
	ctx context.Context,
	dB db.DB,
	orderID int64,
) (*models.Order, error) {
	return s.store.OrderDomain.OrderByID(ctx, dB, orderID)
}

func (s *orderService) getPriceAndOrderItems(
	ctx context.Context,
	operations db.SQLOperations,
	form *dtos.CreateOrderForm,
) ([]*models.OrderItem, float64, error) {

	var totalPrice float64
	orderItems := make([]*models.OrderItem, 0, len(form.Items))

	for _, item := range form.Items {
		product, err := s.store.ProductDomain.ProductByID(ctx, operations, item.ProductID)
		if err != nil {
			loggers.Errorf("failed to get product by id [%d], err: [%+v]", item.ProductID, err)
			return nil, 0, err
		}

		// cost = retail_price * quantity
		cost := product.RetailPrice * float64(item.Quantity)

		orderItem := &models.OrderItem{
			ProductID:   product.ID,
			Quantity:    item.Quantity,
			UnitPrice:   product.RetailPrice,
			TotalAmount: cost,
		}

		orderItems = append(orderItems, orderItem)
		totalPrice += cost
	}

	return orderItems, totalPrice, nil
}

func (s *orderService) createOrderItemsBatches(
	orderItems []*models.OrderItem,
) [][]*models.OrderItem {

	batchSize := 6000

	expectedBatches := utils.CreateBatches(len(orderItems), batchSize)

	orderItemsBatches := make([][]*models.OrderItem, expectedBatches)

	for i := 0; i < expectedBatches; i++ {
		start := i * batchSize

		end := ((i + 1) * batchSize)
		if end > len(orderItems) {
			end = len(orderItems)
		}

		orderItemsBatches[i] = orderItems[start:end]
	}

	return orderItemsBatches
}
