package main

import (
	"fmt"
	_wordFormRepo "github.com/DostonAkhmedov/lucy/brand/wordforms/repository"
	"github.com/DostonAkhmedov/lucy/config"
	"github.com/DostonAkhmedov/lucy/driver/mysql"
	_elementPropertyRepo "github.com/DostonAkhmedov/lucy/iblock/element/property/repository"
	_elementRepo "github.com/DostonAkhmedov/lucy/iblock/element/repository"
	_propertyRepo "github.com/DostonAkhmedov/lucy/iblock/property/repository"
	_iblockRepo "github.com/DostonAkhmedov/lucy/iblock/repository"
	_discountRepo "github.com/DostonAkhmedov/lucy/linemedia/discount/repository"
	_linemediaRepo "github.com/DostonAkhmedov/lucy/linemedia/repository"
	_supplierRepo "github.com/DostonAkhmedov/lucy/linemedia/supplier/repository"
	"github.com/DostonAkhmedov/lucy/logger"
	"github.com/DostonAkhmedov/lucy/models"
	"github.com/DostonAkhmedov/lucy/models/iblock"
	"github.com/DostonAkhmedov/lucy/models/iblock/element"
	"github.com/DostonAkhmedov/lucy/models/linemedia"
	"math"
	"time"
)

func main() {

	start := time.Now()

	conf := config.Init()

	logtofile := logger.ToFile("status.log", "")
	slacklogger := logger.SlcLogger(conf)

	logtofile.Println("Start!")
	slacklogger.Info("Start!")

	connection, err := mysql.Connection(conf)
	if err != nil {
		logtofile.Fatal(err)
		slacklogger.Error(err)
	}

	defer func() {
		err := connection.Close()
		if err != nil {
			logtofile.Fatal(err)
			slacklogger.Error(err)
		}
	}()

	supplierRepo := _supplierRepo.NewSupplierRepository(connection)
	suppliers, err := supplierRepo.GetList()
	if err != nil {
		logtofile.Fatal(err)
		slacklogger.Error(err)
	}
	logtofile.Printf("%d suppliers found.", len(suppliers))

	discountRepo := _discountRepo.NewDiscountRepository(connection)
	discounts, err := discountRepo.GetList(suppliers)
	if err != nil {
		logtofile.Fatal(err)
		slacklogger.Error(err)
	}

	propertyRepo := _propertyRepo.NewPropertyRepository(connection)
	elementRepo := _elementRepo.NewElementRepository(connection)
	linemediaRepo := _linemediaRepo.NewLinemediaRepository(connection)
	elementPropertyRepo := _elementPropertyRepo.NewPropertyRepository(connection)
	brandsRepo := _wordFormRepo.NewWordFormsRepository(connection)

	iblockRepo := _iblockRepo.NewIblockRepository(connection)
	iblockIds, err := iblockRepo.GetIblockIds()
	if err != nil {
		logtofile.Fatal(err)
		slacklogger.Error(err)
	}
	logtofile.Printf("%d iblocks found.", len(iblockIds))

	var cntadded, cntupdated = 0, 0
	propertyCode := "CML2_BESTPRICE"
	for _, ibId := range iblockIds {
		prop, err := propertyRepo.GetByCode(ibId, propertyCode)
		if err != nil {
			logtofile.Fatal(err)
			slacklogger.Error(err)
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
				logtofile.Fatal(err)
				slacklogger.Error(err)
			}
		}
		elementList, err := elementRepo.GetList(ibId)
		if err != nil {
			logtofile.Fatal(err)
			slacklogger.Error(err)
		}
		logtofile.Printf("%d products found with iblock=%d.", len(elementList), ibId)
		slacklogger.Info(fmt.Sprintf("%d products found with iblock=%d.", len(elementList), ibId))
		for _, el := range elementList {
			el.Article = elementRepo.FormatArticle(el.Article)
			var minPrice = el.MinPrice
			if len(el.Article) > 0 && len(el.Brand) > 0 {
				brands, err := brandsRepo.GetWordForms(el.Brand)
				if err != nil {
					logtofile.Fatal(err)
					slacklogger.Error(err)
				}

				lmProducts, err := linemediaRepo.GetList(el.Article, brands, suppliers)
				if err != nil {
					logtofile.Fatal(err)
					slacklogger.Error(err)
				}
				for _, part := range lmProducts {
					if _, prs := discounts[part.Supplier]; prs {
						part.Price = recalc(part, discounts[part.Supplier])
					}
					minPrice = math.Min(minPrice, part.Price)
				}
				minPrice = math.Ceil(minPrice)

				elementProperty, err := elementPropertyRepo.GetById(prop.Id, el.Id)
				if err != nil {
					logtofile.Fatal(err)
					slacklogger.Error(err)
				}
				if elementProperty == nil {
					elementProperty.Id, err = elementPropertyRepo.Add(&element.Property{
						IblockPropertyId: prop.Id,
						ElementId:        el.Id,
						Value:            minPrice,
					})
					if err != nil {
						logtofile.Fatal(err)
						slacklogger.Error(err)
					}
					cntadded++
				} else if elementProperty.Value != minPrice {
					_, err := elementPropertyRepo.Update(elementProperty.Id, minPrice)
					if err != nil {
						logtofile.Fatal(err)
						slacklogger.Error(err)
					}
					cntupdated++
				}
			}
		}
	}
	logtofile.Printf("Added new rows: %d", cntadded)
	logtofile.Printf("Updated rows: %d", cntupdated)

	slacklogger.Info(fmt.Sprintf("Added new rows: %d", cntadded))
	slacklogger.Info(fmt.Sprintf("Updated rows: %d", cntupdated))

	elapsed := time.Since(start)
	logtofile.Printf("Finished! runned: %s", elapsed)
	slacklogger.Info(fmt.Sprintf("Finished! runned: %s", elapsed))
}

func recalc(part *models.LMProduct, discounts []*linemedia.Discount) float64 {
	price := part.Price
	for _, discount := range discounts {
		if part.Supplier != discount.SupplierId {
			continue
		}
		if discount.MinPrice.Valid && part.Price < discount.MinPrice.Float64 {
			continue
		}
		if discount.MaxPrice.Valid && part.Price > discount.MaxPrice.Float64 {
			continue
		}
		price += part.Price * (discount.Percent / 100)
	}

	return price
}
