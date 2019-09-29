package main

import (
	_wordFormRepo "github.com/DostonAkhmedov/lucy/brand/wordforms/repository"
	"github.com/DostonAkhmedov/lucy/config"
	"github.com/DostonAkhmedov/lucy/driver/mysql"
	_elementRepo "github.com/DostonAkhmedov/lucy/iblock/element/repository"
	_propertyRepo "github.com/DostonAkhmedov/lucy/iblock/property/repository"
	_iblockRepo "github.com/DostonAkhmedov/lucy/iblock/repository"
	_linemediaRepo "github.com/DostonAkhmedov/lucy/linemedia/repository"
	_supplierRepo "github.com/DostonAkhmedov/lucy/linemedia/supplier/repository"
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

	brandsRepo := _wordFormRepo.NewWordFormsRepository(connection)
	brandsForms, err := brandsRepo.GetWordForms()
	if err != nil {
		panic(err)
	}

	supplierRepo := _supplierRepo.NewSupplierRepository(connection)
	suppliers, err := supplierRepo.GetList()
	if err != nil {
		panic(err)
	}

	propertyRepo := _propertyRepo.NewPropertyRepository(connection)
	elementRepo := _elementRepo.NewElementRepository(connection)
	linemediaRepo := _linemediaRepo.NewLinemediaRepository(connection)

	iblockRepo := _iblockRepo.NewIblockRepository(connection)
	iblockIds, err := iblockRepo.GetIblockIds()
	if err != nil {
		panic(err)
	}

	propertyCode := "CML2_BESTPRICE"
	for _, ibId := range iblockIds {
		prop, err := propertyRepo.GetByCode(ibId, propertyCode)
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
			prop.Id, err = propertyRepo.Add(prop)
			if err != nil {
				panic(err)
			}
		}
		elementList, err := elementRepo.GetList(ibId)
		if err != nil {
			panic(err)
		}
		for _, el := range elementList {
			el.Article = elementRepo.FormatArticle(el.Article)
			if len(el.Article) > 0 && len(el.Brand) > 0 {
				var brands []string
				if _, prs := brandsForms[el.Brand]; prs {
					brands = brandsForms[el.Brand]
				} else {
					brands = []string{el.Brand}
				}
				lmProducts, err := linemediaRepo.GetList(el.Article, brands, suppliers)
				if err != nil {
					panic(err)
				}
				log.Println(lmProducts)
			}
		}
	}
}
