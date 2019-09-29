package repository

import (
	"database/sql"
	"fmt"
	"github.com/DostonAkhmedov/lucy/iblock/element/property"
	"github.com/DostonAkhmedov/lucy/models/iblock/element"
	"log"
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
		"ORDER BY ID DESC"+
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

	defer func() {
		err := st.Close()
		if err != nil {
			log.Fatal(err)
		}
	}()

	result, err := st.Exec(
		prop.IblockPropertyId,
		prop.ElementId,
		prop.Value,
	)

	if err != nil {
		log.Fatal(err)
		return 0, err
	}

	id, _ := result.LastInsertId()
	return id, nil
}

func (p *propertyRepository) Update(id int64, value float64) (int64, error) {
	query := fmt.Sprintf("UPDATE %s SET VALUE=? WHERE ID=?;", tableName)
	result, err := p.Conn.Exec(query, value, id)
	if err != nil {
		log.Fatal(err)
		return 0, err
	}

	return result.RowsAffected()
}
