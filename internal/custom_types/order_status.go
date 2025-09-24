package custom_types

import "database/sql/driver"

type OrderStatus string

const (
	OrderStatusPending OrderStatus = "pending"
	OrderStatusCancelled OrderStatus = "cancelled"
	OrderStatusPaid OrderStatus = "paid"
	OrderStatusReturned OrderStatus = "returned"
)

func(o *OrderStatus) Scan(value interface{}) error {
	*o = OrderStatus(string(value.([]uint8)))
	return nil 
}

func (o OrderStatus) Value() (driver.Value, error) {
	return o.String(), nil
}

func (o OrderStatus) String() string {
	return string(o)
}