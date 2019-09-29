package supplier

type Repository interface {
	GetList() ([]string, error)
}
