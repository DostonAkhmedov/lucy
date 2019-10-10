package linemedia

import "database/sql"

type Discount struct {
	Id         int64
	IblockId   int64
	Name       string
	SupplierId string
	MinPrice   sql.NullFloat64
	MaxPrice   sql.NullFloat64
	Percent    float64
}
