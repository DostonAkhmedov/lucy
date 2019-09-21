package main

import (
	"github.com/DostonAkhmedov/lucy/config"
	"github.com/DostonAkhmedov/lucy/driver/mysql"
	propertyRepo "github.com/DostonAkhmedov/lucy/iblock/property/repository"
	iblockRepo "github.com/DostonAkhmedov/lucy/iblock/repository"
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

	_iblock := iblockRepo.NewIblockRepository(connection)
	iblockIds, err := _iblock.GetIblockIds()
	if err != nil {
		panic(err)
	}

	_property := propertyRepo.NewPropertyRepository(connection)
	propertyCode := "CML2_BESTPRICE"
	for _, ibId := range iblockIds {
		prop, err := _property.GetByCode(ibId, propertyCode)
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
			prop.Id, err = _property.Add(prop)
			if err != nil {
				panic(err)
			}
		}
	}
}
