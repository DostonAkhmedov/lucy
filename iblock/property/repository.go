package property

import (
	"github.com/DostonAkhmedov/lucy/models/iblock"
)

type Repository interface {
	Add(prop *iblock.Property) (id int64, err error)
	GetByCode(ibId int64, propCode string) (property *iblock.Property, err error)
	Update(prop *iblock.Property) error
}
