package repository

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"savebestprice/config"
	"savebestprice/iblock"
	"savebestprice/models"
	"strings"
)

const tableName string = "b_iblock"

type iblockRepository struct {
	Conn *sql.DB
}

func NewIblockRepository(Conn *sql.DB) iblock.Repository {
	return &iblockRepository{Conn}
}

func (ib *iblockRepository) GetList(ctx context.Context) ([]*models.Iblock, error) {
	iblockIds, err := ib.GetIblockIds()
	if err != nil {
		log.Fatal(err)
		return nil, err
	}
	sqlStr := fmt.Sprintf("SELECT ID, CODE, NAME FROM %s WHERE ACTIVE='Y' AND ID IN (%s)", tableName, ToString(iblockIds))
	rows, err := ib.Conn.Query(sqlStr)

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
	iblocks := make([]*models.Iblock, 0)

	for rows.Next() {
		ib := new(models.Iblock)
		err = rows.Scan(&ib.Id, &ib.Code, &ib.Name)
		if err != nil {
			log.Fatal(err)
			return nil, err
		}
		iblocks = append(iblocks, ib)
	}

	return iblocks, nil
}

func (ib *iblockRepository) GetIblockIds(codes ...string) ([]int64, error) {
	if len(codes) == 0 {
		for k, _ := range config.CatalogIblocks {
			codes = append(codes, k)
		}
	}

	var result = make([]int64, 0)
	sqlStr := fmt.Sprintf("SELECT ID FROM %s WHERE ACTIVE='Y' AND XML_ID IN (%s);", tableName, ToString(codes))
	rows, err := ib.Conn.Query(sqlStr)

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

	for rows.Next() {
		var id int64
		err = rows.Scan(&id)
		if err != nil {
			log.Fatal(err)
			return nil, err
		}
		result = append(result, id)
	}

	return result, nil
}

func ToString(arr ...interface{}) string {
	return strings.Trim(strings.Join(strings.Fields(fmt.Sprint(arr)), ","), "[]{}")
}