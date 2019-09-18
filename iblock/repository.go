package iblock

import (
	"context"
	"savebestprice/models"
)

type Repository interface {
	GetList(ctx context.Context) ([]*models.Iblock, error)
	GetIblockIds(codes ...string) ([]int64, error)
}
