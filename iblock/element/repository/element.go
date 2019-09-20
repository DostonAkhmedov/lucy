package repository

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/DostonAkhmedov/lucy/iblock/element"
	"github.com/DostonAkhmedov/lucy/models/iblock"
	"log"
)

const tableName string = "b_iblock_element"

type elementRepository struct {
	Conn *sql.DB
}

func NewElementRepository(Conn *sql.DB) element.Repository {
	return &elementRepository{Conn}
}

func (el *elementRepository) GetList(ctx context.Context, iblockId int64) ([]*iblock.Element, error) {
	sqlStr := fmt.Sprintf("SELECT ID FROM %s WHERE ACTIVE='Y' AND IBLOCK_ID=%d", tableName, iblockId)
	rows, err := el.Conn.Query(sqlStr)
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

	elements := make([]*iblock.Element, 0)

	for rows.Next() {
		el := new(iblock.Element)
		err = rows.Scan(&el.Id)
		if err != nil {
			log.Fatal(err)
			return nil, err
		}
		elements = append(elements, el)
	}

	return elements, nil
}
