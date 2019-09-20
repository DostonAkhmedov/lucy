package element

import (
	"context"
	"github.com/DostonAkhmedov/lucy/models/iblock"
)

type Repository interface {
	GetList(ctx context.Context, iblockId int64) ([]*iblock.Element, error)
}
