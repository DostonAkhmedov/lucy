package main

import (
	"github.com/DostonAkhmedov/lucy/config"
	"github.com/DostonAkhmedov/lucy/driver/mysql"
	_elementRepo "github.com/DostonAkhmedov/lucy/iblock/element/repository"
	_propertyRepo "github.com/DostonAkhmedov/lucy/iblock/property/repository"
	_iblockRepo "github.com/DostonAkhmedov/lucy/iblock/repository"
	"github.com/DostonAkhmedov/lucy/models/iblock"
	"log"
)

func main() {

	dbConfig := config.Init()
	connection, err := mysql.Connection(dbConfig)
	if err != nil {
		log.Fatal(err)
	}

	defer func() {
		err := connection.Close()
		if err != nil {
			log.Fatal(err)
		}
	}()

	iblockRepo := _iblockRepo.NewIblockRepository(connection)
	iblockIds, err := iblockRepo.GetIblockIds()
	if err != nil {
		panic(err)
	}

	property := _propertyRepo.NewPropertyRepository(connection)
	propertyCode := "CML2_BESTPRICE"
	element := _elementRepo.NewElementRepository(connection)
	for _, ibId := range iblockIds {
		prop, err := property.GetByCode(ibId, propertyCode)
		if err != nil {
			panic(err)
		}
		if prop == nil {
			prop = &iblock.Property{
				Active:       "Y",
				Sort:         1,
				Code:         propertyCode,
				PropertyType: "N",
				IblockId:     ibId,
				Name:         "Цена",
				Multiple:     "N",
				MultipleCnt:  5,
				IsRequired:   "N",
			}
			prop.Id, err = property.Add(prop)
			if err != nil {
				panic(err)
			}
		}
		elementList, err := element.GetList(ibId)
		if err != nil {
			panic(err)
		}
		for _, el := range elementList {
			el.Article = element.FormatArticle(el.Article)
			if len(el.Article) > 0 && len(el.Brand) > 0 {
				log.Println(el)
			}
		}
	}
}
