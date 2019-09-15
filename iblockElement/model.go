package iblockElement

import (
	"fmt"
	"log"
	"savebestprice/config"
)

type Element struct {
	Id int64
}

const tableName string = "b_iblock_element"

func GetList(iblockId int64) []Element {
	sqlStr := fmt.Sprintf("SELECT ID FROM %s WHERE ACTIVE='Y' AND IBLOCK_ID=%d", tableName, iblockId)
	rows, err := config.DB.Query(sqlStr)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
	elements := []Element{}

	for rows.Next() {
		el := Element{}
		err = rows.Scan(&el.Id)
		if err != nil {
			log.Fatal(err)
		}
		elements = append(elements, el)
	}

	return elements
}
