package repository

import (
	"database/sql"
	"fmt"
	"github.com/DostonAkhmedov/lucy/iblock/property"
	"github.com/DostonAkhmedov/lucy/models/iblock"
	"log"
)

const tableName string = "b_iblock_property"

type propertyRepository struct {
	Conn *sql.DB
}

func NewPropertyRepository(Conn *sql.DB) property.Repository {
	return &propertyRepository{Conn}
}

func (p *propertyRepository) GetByCode(ibId int64, propCode string) (*iblock.Property, error) {
	query := fmt.Sprintf("SELECT ID, ACTIVE, SORT, CODE, PROPERTY_TYPE, IBLOCK_ID, NAME, MULTIPLE, MULTIPLE_CNT, IS_REQUIRED "+
		"FROM %s WHERE IBLOCK_ID=%d AND CODE='%s';", tableName, ibId, propCode)

	prop := new(iblock.Property)
	row := p.Conn.QueryRow(query)
	switch err := row.Scan(
		&prop.Id,
		&prop.Active,
		&prop.Sort,
		&prop.Code,
		&prop.PropertyType,
		&prop.IblockId,
		&prop.Name,
		&prop.Multiple,
		&prop.MultipleCnt,
		&prop.IsRequired); err {
	case sql.ErrNoRows:
		return nil, nil
	case nil:
		return prop, nil
	default:
		return nil, err
	}

}

func (p *propertyRepository) Add(prop *iblock.Property) (int64, error) {
	var id int64

	query := fmt.Sprintf("INSERT INTO "+
		"%s(ACTIVE, SORT, CODE, PROPERTY_TYPE, IBLOCK_ID, NAME, MULTIPLE, MULTIPLE_CNT, IS_REQUIRED) "+
		"VALUES(?, ?, ?, ?, ?, ?, ?, ?, ?);",
		tableName)
	st, _ := p.Conn.Prepare(query)

	defer func() {
		err := st.Close()
		if err != nil {
			log.Fatal(err)
		}
	}()

	result, err := st.Exec(
		prop.Active,
		prop.Sort,
		prop.Code,
		prop.PropertyType,
		prop.IblockId,
		prop.Name,
		prop.Multiple,
		prop.MultipleCnt,
		prop.IsRequired,
	)
	if err != nil {
		log.Fatal(err)
		return -1, err
	}
	id, _ = result.LastInsertId()
	return id, nil
}

func (p *propertyRepository) Update(prop *iblock.Property) error {
	panic("implement me")
}
