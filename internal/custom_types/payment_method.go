package custom_types

import "database/sql/driver"

type PaymentMethod string

const (
	PaymentMethodCash PaymentMethod = "cash"
	PaymentMethodCard PaymentMethod = "card"
	PaymentMethodMpesa PaymentMethod = "mpesa"
)

func(p *PaymentMethod) Scan(value interface{}) error {
	*p = PaymentMethod(string(value.([]uint8)))
	return nil 
}

func (p PaymentMethod) Value() (driver.Value, error) {
	return p.String(), nil
}

func (p PaymentMethod) String() string {
	return string(p)
}