package repository

import (
	"database/sql"
	"fmt"
	"github.com/DostonAkhmedov/lucy/linemedia/supplier"
	"log"
)

const tableName string = "b_lm_sphinx_active_suppliers"

type supplierRepository struct {
	Conn *sql.DB
}

func NewSupplierRepository(Conn *sql.DB) supplier.Repository {
	return &supplierRepository{Conn}
}

func (s *supplierRepository) GetList() ([]string, error) {
	query := fmt.Sprintf("SELECT id FROM %s", tableName)
	rows, err := s.Conn.Query(query)
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
		suppliers = make([]string, 0)
		sp        string
	)
	for rows.Next() {
		err := rows.Scan(&sp)
		if err != nil {
			log.Fatal(err)
			return nil, err
		}

		suppliers = append(suppliers, sp)
	}

	return suppliers, nil
}
