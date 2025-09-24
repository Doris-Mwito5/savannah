package custom_types

import "database/sql/driver"

type CustomerType string

const (
	CustomerTypeIndividual CustomerType = "individual"
	CustomerTypeBusiness CustomerType = "business"  
)

func(c *CustomerType) Scan(value interface{}) error {
	*c = CustomerType(string(value.([]uint8)))
	return nil 
}

func (c CustomerType) Value() (driver.Value, error) {
	return c.String(), nil
}

func (c CustomerType) String() string {
	return string(c)
}