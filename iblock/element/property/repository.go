package property

import "github.com/DostonAkhmedov/lucy/models/iblock/element"

type Repository interface {
	GetById(ibPropId int64, elementId int64) (*element.Property, error)
	Add(prop *element.Property) (int64, error)
	Update(id int64, value float64) (int64, error)
}
