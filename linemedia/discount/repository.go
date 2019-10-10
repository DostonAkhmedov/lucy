package discount

import "github.com/DostonAkhmedov/lucy/models/linemedia"

type Repository interface {
	GetList(suppliers []string) (map[string][]*linemedia.Discount, error)
}
