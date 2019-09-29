package repository

import (
	"database/sql"
	"fmt"
	"github.com/DostonAkhmedov/lucy/iblock/element"
	"github.com/DostonAkhmedov/lucy/models/iblock"
	"log"
	"regexp"
	"strings"
)

const tableName string = "b_iblock_element"

type elementRepository struct {
	Conn *sql.DB
}

func NewElementRepository(Conn *sql.DB) element.Repository {
	return &elementRepository{Conn}
}

func (el *elementRepository) GetList(iblockId int64) ([]*iblock.Element, error) {
	query := fmt.Sprintf("SELECT BE.ID AS ID, PV0.VALUE AS BRAND, PV1.VALUE AS ARTICLE, PR.QUANTITY AS QUANTITY, MIN(CEIL(P.PRICE)) AS PRICE "+
		"FROM %s BE "+
		"LEFT JOIN b_iblock B ON B.ID = IBLOCK_ID "+
		"RIGHT JOIN b_catalog_product PR ON PR.ID = BE.ID "+
		"RIGHT JOIN b_catalog_price P ON P.PRODUCT_ID = BE.ID "+
		"LEFT JOIN b_iblock_property P0 ON P0.IBLOCK_ID = B.ID AND P0.CODE='BRAND' "+
		"LEFT JOIN b_iblock_property P1 ON P1.IBLOCK_ID = B.ID AND P1.CODE='CML2_ARTICLE' "+
		"RIGHT JOIN b_iblock_element_property PV0 ON PV0.IBLOCK_PROPERTY_ID = P0.ID AND PV0.IBLOCK_ELEMENT_ID = BE.ID "+
		"RIGHT JOIN b_iblock_element_property PV1 ON PV1.IBLOCK_PROPERTY_ID = P1.ID AND PV1.IBLOCK_ELEMENT_ID = BE.ID "+
		"WHERE BE.ACTIVE='Y' AND BE.IBLOCK_ID=%d AND PR.QUANTITY > 0 AND CEIL(P.PRICE) > 0 "+
		"GROUP BY BE.ID;",
		tableName,
		iblockId,
	)
	rows, err := el.Conn.Query(query)
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

	var elements = make([]*iblock.Element, 0)
	for rows.Next() {
		el := new(iblock.Element)
		err = rows.Scan(&el.Id, &el.Brand, &el.Article, &el.Quantity, &el.MinPrice)
		if err != nil {
			log.Fatal(err)
			return nil, err
		}

		elements = append(elements, el)
	}

	return elements, nil
}

func (el *elementRepository) FormatArticle(article string) string {
	r := strings.NewReplacer(" ", "",
		"-", "",
		"/", "",
		"\\", "",
		".", "",
		"\"", "",
		"'", "",
		"\r", "",
		"\n", "",
		"\t", "")

	article = strings.ToLower(r.Replace(article))
	article = regexp.MustCompile(`/\s*\([^)]*\)/`).ReplaceAllString(article, "")
	article = regexp.MustCompile("[^A-Za-zА-Яа-яЁё0-9)(_]").ReplaceAllString(article, "")

	return article
}
