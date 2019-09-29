package linemedia

import "github.com/DostonAkhmedov/lucy/models"

type Repository interface {
	GetList(article string, brands []string, suppliers []string) ([]*models.LMProduct, error)
}
