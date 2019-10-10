package repository

import (
	"database/sql"
	"fmt"
	"github.com/DostonAkhmedov/lucy/linemedia"
	"github.com/DostonAkhmedov/lucy/models"
	"log"
	"strings"
)

const tableName string = "b_lm_products"

type linemediaRepository struct {
	Conn *sql.DB
}

func NewLinemediaRepository(Conn *sql.DB) linemedia.Repository {
	return &linemediaRepository{Conn}
}

func (lm *linemediaRepository) GetList(article string, brands []string, suppliers []string) ([]*models.LMProduct, error) {
	articles := []string{article, "0" + article}
	query := fmt.Sprintf("SELECT id, article, UPPER(brand_title) AS brand, price, supplier_id "+
		"FROM %s "+
		"WHERE price > 0 AND quantity > 0 AND "+
		"article IN('%s') AND UPPER(brand_title) IN('%s') AND supplier_id IN('%s') "+
		"LIMIT 20;",
		tableName,
		strings.Join(articles, "','"),
		strings.Join(brands, "','"),
		strings.Join(suppliers, "','"),
	)
	rows, err := lm.Conn.Query(query)
	if err != nil {
		log.Println(query)
		log.Fatal(err)
		return nil, err
	}

	defer func() {
		err := rows.Close()
		if err != nil {
			log.Fatal(err)
		}
	}()

	parts := make([]*models.LMProduct, 0)
	for rows.Next() {
		part := new(models.LMProduct)
		err := rows.Scan(&part.Id, &part.Article, &part.Brand, &part.Price, &part.Supplier)
		if err != nil {
			log.Fatal(err)
			return nil, err
		}

		parts = append(parts, part)
	}

	return parts, nil
}
