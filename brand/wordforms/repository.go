package wordforms

type Repository interface {
	GetWordForms(brand string) ([]string, error)
	GetByGroup(group string) ([]string, error)
	GetGroup(brand string) (string, error)
	ClearBrand(brand string) string
}
