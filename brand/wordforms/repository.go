package wordforms

type Repository interface {
	GetWordForms() (map[string][]string, error)
	ClearBrand(brand string) string
}
