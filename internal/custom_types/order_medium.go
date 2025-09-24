package custom_types

import "database/sql/driver"

type OrderMedium string

const (
	OrderMediumOnline OrderMedium = "online"
	OrderMediumoffline OrderMedium = "offine"
)

func(o *OrderMedium) Scan(value interface{}) error {
	*o = OrderMedium(string(value.([]uint8)))
	return nil 
}

func (o OrderMedium) Value() (driver.Value, error) {
	return o.String(), nil
}

func (o OrderMedium) String() string {
	return string(o)
}