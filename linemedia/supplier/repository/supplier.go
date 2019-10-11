package repository

import (
	"database/sql"
	"fmt"
	"github.com/DostonAkhmedov/lucy/linemedia/supplier"
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
		return nil, err
	}

	defer rows.Close()

	var (
		suppliers = make([]string, 0)
		sp        string
	)
	for rows.Next() {
		err := rows.Scan(&sp)
		if err != nil {
			return nil, err
		}

		suppliers = append(suppliers, sp)
	}

	return suppliers, nil
}
