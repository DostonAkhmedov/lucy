package iblock

import (
	"fmt"
	"log"
	"savebestprice/config"
	"strings"
)

type Iblock struct {
	Id   int64
	Code string
	Name string
}

const tableName string = "b_iblock"

func GetList() []Iblock {
	var iblockIds = GetIblockIds()
	sqlStr := fmt.Sprintf("SELECT ID, CODE, NAME FROM %s WHERE ACTIVE='Y' AND ID IN (%s)", tableName, ToString(iblockIds))
	rows, err := config.DB.Query(sqlStr)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
	iblocks := []Iblock{}

	for rows.Next() {
		ib := Iblock{}
		err = rows.Scan(&ib.Id, &ib.Code, &ib.Name)
		if err != nil {
			log.Fatal(err)
		}
		iblocks = append(iblocks, ib)
	}

	return iblocks
}

func GetIblockIds(codes ...string) []int64 {
	if len(codes) == 0 {
		for k, _ := range CatalogIblocks {
			codes = append(codes, k)
		}
	}
	var result = []int64{}
	sqlStr := fmt.Sprintf("SELECT ID FROM %s WHERE ACTIVE='Y' AND XML_ID IN (%s);", tableName, ToString(codes))
	rows, err := config.DB.Query(sqlStr)
	fmt.Println(ToString(codes))
	if err != nil {
		log.Fatal(err)
	}
	for rows.Next() {
		var id int64
		err = rows.Scan(&id)
		if err != nil {
			log.Fatal(err)
		}
		result = append(result, id)
	}

	return result
}

func ToString(arr ...interface{}) string {
	return strings.Trim(strings.Join(strings.Fields(fmt.Sprint(arr)), ","), "[]{}")
}