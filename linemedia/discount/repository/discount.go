package repository

import (
	"database/sql"
	"fmt"
	"github.com/DostonAkhmedov/lucy/linemedia/discount"
	"github.com/DostonAkhmedov/lucy/models/linemedia"
	"log"
	"strings"
)

const iblockId int64 = 8

type discountRepository struct {
	Conn *sql.DB
}

func NewDiscountRepository(Conn *sql.DB) discount.Repository {
	return &discountRepository{Conn}
}

func (d *discountRepository) GetList(suppliers []string) (map[string][]*linemedia.Discount, error) {
	defaultUserGroups := []string{"2", // not authorized
		"37",    // user group SPB City for price
		"guest", // user group for guest users
	}
	query := fmt.Sprintf(
		"SELECT BE.ID AS ID, BE.IBLOCK_ID AS IBLOCK_ID, BE.NAME AS NAME, SPIDV.VALUE AS SUPPLIER_ID, "+
			"PV0.VALUE AS MIN_PRICE, PV1.VALUE AS MAX_PRICE, FORMAT(PV2.VALUE, 'N') AS PERCENT "+
			"FROM b_iblock_element BE "+
			"LEFT JOIN b_iblock B ON B.ID = IBLOCK_ID "+
			"LEFT JOIN b_iblock_property SP ON SP.IBLOCK_ID = B.ID AND SP.CODE='supplier_id' "+
			"LEFT JOIN b_iblock_property SPID ON SPID.IBLOCK_ID = 5 AND SPID.CODE='supplier_id' "+
			"LEFT JOIN b_iblock_property P0 ON P0.IBLOCK_ID = B.ID AND P0.CODE='price_min' "+
			"LEFT JOIN b_iblock_property P1 ON P1.IBLOCK_ID = B.ID AND P1.CODE='price_max' "+
			"LEFT JOIN b_iblock_property P2 ON P2.IBLOCK_ID = B.ID AND P2.CODE='discount' "+
			"LEFT JOIN b_iblock_property P3 ON P3.IBLOCK_ID = B.ID AND P3.CODE='user_group' "+
			"RIGHT JOIN b_iblock_element_property SPV ON SPV.IBLOCK_PROPERTY_ID = SP.ID AND SPV.IBLOCK_ELEMENT_ID = BE.ID "+
			"RIGHT JOIN b_iblock_element_property SPIDV ON SPIDV.IBLOCK_PROPERTY_ID = SPID.ID AND SPIDV.IBLOCK_ELEMENT_ID = SPV.VALUE "+
			"LEFT JOIN b_iblock_element_property PV0 ON PV0.IBLOCK_PROPERTY_ID = P0.ID AND PV0.IBLOCK_ELEMENT_ID = BE.ID "+
			"LEFT JOIN b_iblock_element_property PV1 ON PV1.IBLOCK_PROPERTY_ID = P1.ID AND PV1.IBLOCK_ELEMENT_ID = BE.ID "+
			"RIGHT JOIN b_iblock_element_property PV2 ON PV2.IBLOCK_PROPERTY_ID = P2.ID AND PV2.IBLOCK_ELEMENT_ID = BE.ID "+
			"LEFT JOIN b_iblock_element_property PV3 ON PV3.IBLOCK_PROPERTY_ID = P3.ID AND PV3.IBLOCK_ELEMENT_ID = BE.ID "+
			"WHERE BE.ACTIVE='Y' AND BE.IBLOCK_ID=%d AND (PV3.VALUE IS NULL OR PV3.VALUE IN('%s')) "+
			"ORDER BY ID ASC;",
		iblockId,
		strings.Join(defaultUserGroups, "','"),
	)

	rows, err := d.Conn.Query(query)
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

	var discounts = make(map[string][]*linemedia.Discount)
	for rows.Next() {
		disc := new(linemedia.Discount)
		err = rows.Scan(
			&disc.Id,
			&disc.IblockId,
			&disc.Name,
			&disc.SupplierId,
			&disc.MinPrice,
			&disc.MaxPrice,
			&disc.Percent)
		if err != nil {
			log.Fatal(err)
			return nil, err
		}

		discounts[disc.SupplierId] = append(discounts[disc.SupplierId], disc)
	}

	return discounts, nil
}
