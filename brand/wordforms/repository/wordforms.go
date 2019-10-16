package repository

import (
	"database/sql"
	"fmt"
	"github.com/DostonAkhmedov/lucy/brand/wordforms"
	"strings"
)

const tableName string = "b_lm_wordforms"

type wordFormsRepository struct {
	Conn *sql.DB
}

func NewWordFormsRepository(Conn *sql.DB) wordforms.Repository {
	return &wordFormsRepository{Conn}
}

func (wf *wordFormsRepository) GetWordForms(brand string) ([]string, error) {

	brand = wf.ClearBrand(brand)

	brands, err := wf.GetByGroup(brand)

	if err != nil {
		return nil, err
	}

	if len(brands) > 0 {
		brands = append(brands, brand)
		return brands, nil
	}

	group, err := wf.GetGroup(brand)
	if err != nil {
		return nil, err
	}

	if len(group) > 0 {
		brands, err = wf.GetByGroup(group)
		if err == nil {
			brands = append(brands, group)
		}
		return brands, err
	}

	return []string{brand}, nil
}

func (wf *wordFormsRepository) GetByGroup(brand string) ([]string, error) {
	query := fmt.Sprintf("SELECT UPPER(`brand_title`) AS `brand_title`, UPPER(`group`) AS `group` "+
		"FROM %s "+
		"WHERE `group`='%s' "+
		"ORDER BY `group`;",
		tableName,
		brand)
	rows, err := wf.Conn.Query(query)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var (
		brandsForms = make([]string, 0)
		group, word string
	)
	for rows.Next() {
		err := rows.Scan(&word, &group)
		if err != nil {
			return nil, err
		}

		brandsForms = append(brandsForms, wf.ClearBrand(word))
	}

	return brandsForms, nil
}

func (wf *wordFormsRepository) GetGroup(brand string) (string, error) {
	query := fmt.Sprintf("SELECT UPPER(`group`) AS `group` "+
		"FROM %s "+
		"WHERE `brand_title`='%s' "+
		"ORDER BY `group` "+
		"LIMIT 1;",
		tableName,
		brand)
	row := wf.Conn.QueryRow(query)
	var group string
	switch err := row.Scan(&group); err {
	case sql.ErrNoRows:
		return "", nil
	case nil:
		return wf.ClearBrand(group), nil
	default:
		return "", err
	}
}

func (wf *wordFormsRepository) ClearBrand(brand string) string {
	replace := map[string]string{"\\": "\\\\", "'": `''`, "\\0": "\\\\0", "\n": "\\n", "\r": "\\r", `"`: `\"`, "\x1a": "\\Z"}

	for b, a := range replace {
		brand = strings.Replace(brand, b, a, -1)
	}

	return brand
}
