package repository

import (
	"database/sql"
	"fmt"
	"github.com/DostonAkhmedov/lucy/iblock"
	"github.com/DostonAkhmedov/lucy/models"
	"strings"
)

const tableName string = "b_iblock"

type iblockRepository struct {
	Conn *sql.DB
}

func NewIblockRepository(Conn *sql.DB) iblock.Repository {
	return &iblockRepository{Conn}
}

func (ib *iblockRepository) GetList() ([]*models.Iblock, error) {
	iblockIds, err := ib.GetIblockIds()
	if err != nil {
		return nil, err
	}
	querty := fmt.Sprintf("SELECT ID, CODE, NAME FROM %s WHERE ACTIVE='Y' AND ID IN (%s)", tableName, ToString(iblockIds))
	rows, err := ib.Conn.Query(querty)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	iblocks := make([]*models.Iblock, 0)

	for rows.Next() {
		ib := new(models.Iblock)
		err = rows.Scan(&ib.Id, &ib.Code, &ib.Name)
		if err != nil {
			return nil, err
		}
		iblocks = append(iblocks, ib)
	}

	return iblocks, nil
}

func (ib *iblockRepository) GetIblockIds(codes ...int) ([]int64, error) {
	if len(codes) == 0 {
		for _, v := range ib.CatalogIblocks() {
			codes = append(codes, v)
		}
	}

	querty := fmt.Sprintf("SELECT ID FROM %s WHERE ACTIVE='Y' AND XML_ID IN (%s);", tableName, ToString(codes))
	rows, err := ib.Conn.Query(querty)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var (
		id     int64
		result = make([]int64, 0)
	)
	for rows.Next() {
		err := rows.Scan(&id)
		if err != nil {
			return nil, err
		}
		result = append(result, id)
	}

	return result, nil
}

func (ib *iblockRepository) CatalogIblocks() []int {
	return []int{
		20102, // => 'ucenennye_tovary',
		20012, // => 'avtoelektronika--12',
		20011, // => 'masla-i-tekhnicheskie-zhidkosti--11',
		20010, // => 'avtosvet--10',
		20008, // => 'instrument',
		20007, // => 'avtokhimiya--avtokosmetika--7',
		20006, // => 'aksessuary--6',
		20005, // => 'akkumulyatornye-batarei--5',
		20014, // => 'soputstvuyushchie-tovary--14',
		20004, // => 'shiny--diski--kolpaki--4',
		//20577,// => 'avtozapchasti--577',
		20080, // => 'krepezh',
		20101, // => 'sport_i_turizm',
		22222, // => 'avtozapchasti',
	}
}

func ToString(arr ...interface{}) string {
	return strings.Trim(strings.Join(strings.Fields(fmt.Sprint(arr)), ","), "[]{}")
}
