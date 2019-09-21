package iblock

import (
	"github.com/DostonAkhmedov/lucy/models"
)

type Repository interface {
	GetList() ([]*models.Iblock, error)
	GetIblockIds(codes ...int) ([]int64, error)
}
