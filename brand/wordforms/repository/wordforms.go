package repository

import (
	"database/sql"
	"fmt"
	"github.com/DostonAkhmedov/lucy/brand/wordforms"
	"log"
	"strings"
)

const tableName string = "b_lm_wordforms"

type wordFormsRepository struct {
	Conn *sql.DB
}

func NewWordFormsRepository(Conn *sql.DB) wordforms.Repository {
	return &wordFormsRepository{Conn}
}

func (wf *wordFormsRepository) GetWordForms() (map[string][]string, error) {
	query := fmt.Sprintf("SELECT UPPER(`brand_title`) AS `brand_title`, UPPER(`group`) AS `group` "+
		"FROM %s ORDER BY `group`;",
		tableName)
	rows, err := wf.Conn.Query(query)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	defer func() {
		err := rows.Close()
		if err != nil {
			log.Fatal(err)
		}
	}()

	var (
		brandsForms = make(map[string][]string)
		group, word string
	)
	for rows.Next() {
		err := rows.Scan(&word, &group)
		if err != nil {
			log.Fatal(err)
			return nil, err
		}

		if _, prs := brandsForms[group]; !prs {
			brandsForms[group] = append(brandsForms[group], wf.ClearBrand(group))
		}
		brandsForms[group] = append(brandsForms[group], wf.ClearBrand(word))
	}

	return brandsForms, nil
}

func (wf *wordFormsRepository) ClearBrand(brand string) string {
	replace := map[string]string{"\\": "\\\\", "'": `''`, "\\0": "\\\\0", "\n": "\\n", "\r": "\\r", `"`: `\"`, "\x1a": "\\Z"}

	for b, a := range replace {
		brand = strings.Replace(brand, b, a, -1)
	}

	return brand
}
