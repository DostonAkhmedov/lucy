package iblock

import "database/sql"

type Element struct {
	Id       int64
	MinPrice float64
	Quantity int
	Article  sql.NullString
	Brand    sql.NullString
}
