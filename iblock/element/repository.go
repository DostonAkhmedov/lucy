package element

import (
	"github.com/DostonAkhmedov/lucy/models/iblock"
)

type Repository interface {
	GetList(iblockId int64) ([]*iblock.Element, error)
	FormatArticle(article string) string
}
