package repository

import (
	"database/sql"
	"fmt"
	"github.com/DostonAkhmedov/lucy/iblock/element/property"
	"github.com/DostonAkhmedov/lucy/models/iblock/element"
)

const tableName string = "b_iblock_element_property"

type propertyRepository struct {
	Conn *sql.DB
}

func NewPropertyRepository(Conn *sql.DB) property.Repository {
	return &propertyRepository{Conn}
}

func (p *propertyRepository) GetById(ibPropId int64, elementId int64) (*element.Property, error) {
	query := fmt.Sprintf("SELECT ID, IBLOCK_PROPERTY_ID, IBLOCK_ELEMENT_ID, VALUE "+
		"FROM %s "+
		"WHERE IBLOCK_PROPERTY_ID=%d AND IBLOCK_ELEMENT_ID=%d "+
		"ORDER BY ID DESC "+
		"LIMIT 1;",
		tableName,
		ibPropId,
		elementId)

	prop := new(element.Property)
	row := p.Conn.QueryRow(query)
	switch err := row.Scan(
		&prop.Id,
		&prop.IblockPropertyId,
		&prop.ElementId,
		&prop.Value); err {
	case sql.ErrNoRows:
		return nil, nil
	case nil:
		return prop, nil
	default:
		return nil, err
	}
}

func (p *propertyRepository) Add(prop *element.Property) (int64, error) {
	query := fmt.Sprintf("INSERT INTO %s(IBLOCK_PROPERTY_ID, IBLOCK_ELEMENT_ID, VALUE) "+
		"VALUES(?, ?, ?);",
		tableName)
	st, _ := p.Conn.Prepare(query)

	defer st.Close()

	result, err := st.Exec(
		prop.IblockPropertyId,
		prop.ElementId,
		prop.Value,
	)

	if err != nil {
		return 0, err
	}

	return result.LastInsertId()
}

func (p *propertyRepository) Update(id int64, value float64) (int64, error) {
	query := fmt.Sprintf("UPDATE %s SET VALUE=?, VALUE_NUM=? WHERE ID=?;", tableName)
	result, err := p.Conn.Exec(query, value, value, id)
	if err != nil {
		return 0, err
	}

	return result.RowsAffected()
}

func (p *propertyRepository) UpdateMultiple(props []*element.Property) error {

	// Get new Transaction.
	txn, err := p.Conn.Begin()
	if err != nil {
		return err
	}

	defer func() {
		// Rollback the transaction after the function returns.
		// If the transaction was already commited, this will do nothing.
		_ = txn.Rollback()
	}()
	for _, prop := range props {
		query := fmt.Sprintf("UPDATE %s SET VALUE=?, VALUE_NUM=? WHERE ID=?;", tableName)

		// Execute the query in the transaction.
		_, err := txn.Exec(query, prop.Value, prop.Value, prop.Id)
		if err != nil {
			return err
		}
	}

	// Commit the transaction.
	return txn.Commit()
}
