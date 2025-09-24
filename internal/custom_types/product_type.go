package custom_types

import "database/sql/driver"

type ProductType string

const (
	ProductTypeGoods ProductType = "goods"
	ProductTypeService ProductType = "service"  
)

func(p *ProductType) Scan(value interface{}) error {
	*p = ProductType(string(value.([]uint8)))
	return nil 
}

func (p ProductType) Value() (driver.Value, error) {
	return p.String(), nil
}

func (p ProductType) String() string {
	return string(p)
}