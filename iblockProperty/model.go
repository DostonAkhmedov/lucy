package iblockProperty

import (
	"database/sql"
	"fmt"
	"log"
	"savebestprice/config"
)

type Property struct {
	Active string
	Sort int
	Code string
	PropertyType string
	IblockId int64
	Name string
	Multiple string
	MultipleCnt int
	IsRequired string
}

const tableName string = "b_iblock_property"

func AddProperty(iblockId int64, propertyCode string) int64 {
	var property = Property{
		Active:"Y",
		Sort:1,
		Code:propertyCode,
		PropertyType:"N",
		IblockId:iblockId,
		Name:"Цена",
		Multiple:"N",
		MultipleCnt:5,
		IsRequired:"N"}

	var id int64

	sqlStr := fmt.Sprintf("SELECT ID FROM %s WHERE ACTIVE='%s' AND IBLOCK_ID=%d AND CODE='%s';",
		tableName, property.Active, property.IblockId, property.Code)
	row := config.DB.QueryRow(sqlStr)
	switch err := row.Scan(&id); err {
	case sql.ErrNoRows:
		sqlStr = fmt.Sprintf("INSERT INTO " +
			"%s(ACTIVE, SORT, CODE, PROPERTY_TYPE, IBLOCK_ID, NAME, MULTIPLE, MULTIPLE_CNT, IS_REQUIRED) " +
			"VALUES(?, ?, ?, ?, ?, ?, ?, ?, ?);",
			tableName)
		st, _ := config.DB.Prepare(sqlStr)
		result, err := st.Exec(
			property.Active,
			property.Sort,
			property.Code,
			property.PropertyType,
			property.IblockId,
			property.Name,
			property.Multiple,
			property.MultipleCnt,
			property.IsRequired)
		if err != nil {
			log.Fatal(err)
		}
		id, _ = result.LastInsertId()
		return id
	case nil:
		return id
	default:
		panic(err)
	}
}
