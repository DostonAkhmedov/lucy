package iblock

import (
	"context"
	"github.com/DostonAkhmedov/lucy/models"
)

type Repository interface {
	GetList(ctx context.Context) ([]*models.Iblock, error)
	GetIblockIds(codes ...int) ([]int64, error)
}
