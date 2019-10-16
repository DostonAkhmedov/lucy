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
	slacklogger := logger.SlcLogger(conf.GetSlcWebHook())

	logtofile.Println("Start!")
	slacklogger.Info("Lucy Start!")

	connection, err := mysql.Connection(conf)
	if err != nil {
		logtofile.Println(err)
		slacklogger.Error(err, "DB connection error!")
	} else {
		slacklogger.Info("DB connection success!")
	}

	defer func() {
		err := connection.Close()
		if err != nil {
			logtofile.Println(err)
			slacklogger.Error(err)
		}
	}()

	supplierRepo := _supplierRepo.NewSupplierRepository(connection)
	suppliers, err := supplierRepo.GetList()
	if err != nil {
		logtofile.Println(err)
		slacklogger.Error(err, "Get supplier error!")
	}
	logtofile.Printf("%d suppliers found.", len(suppliers))

	discountRepo := _discountRepo.NewDiscountRepository(connection)
	discounts, err := discountRepo.GetList(suppliers)
	if err != nil {
		logtofile.Println(err)
		slacklogger.Error(err, "Get discounts error!")
	}

	propertyRepo := _propertyRepo.NewPropertyRepository(connection)
	elementRepo := _elementRepo.NewElementRepository(connection)
	linemediaRepo := _linemediaRepo.NewLinemediaRepository(connection)
	elementPropertyRepo := _elementPropertyRepo.NewPropertyRepository(connection)
	brandsRepo := _wordFormRepo.NewWordFormsRepository(connection)

	iblockRepo := _iblockRepo.NewIblockRepository(connection)
	iblockIds, err := iblockRepo.GetIblockIds()
	if err != nil {
		logtofile.Println(err)
		slacklogger.Error(err, "Get iblockIds error!")
	}
	logtofile.Printf("%d iblocks found.", len(iblockIds))

	var brandsWordForms = make(map[string][]string)

	var cntadded, cntupdated = 0, 0
	propertyCode := "CML2_BESTPRICE"
	for _, ibId := range iblockIds {
		prop, err := propertyRepo.GetByCode(ibId, propertyCode)
		if err != nil {
			logtofile.Println(err)
			slacklogger.Error(err, "Get property error!")
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
				logtofile.Println(err)
				slacklogger.Error(err, "Add property error!")
			}
		}
		elementList, err := elementRepo.GetList(ibId)
		if err != nil {
			logtofile.Println(err)
			slacklogger.Error(err, "Get element list error!")
		}
		logtofile.Printf("%d products found with iblock=%d.", len(elementList), ibId)
		slacklogger.Info(fmt.Sprintf("%d products found with iblock=%d.", len(elementList), ibId))
		for _, el := range elementList {
			el.Article = elementRepo.FormatArticle(el.Article)
			var minPrice = math.MaxFloat64
			if el.Quantity > 0 {
				minPrice = el.MinPrice
			}
			if len(el.Article) > 0 && len(el.Brand) > 0 {
				if _, prs := brandsWordForms[el.Brand]; !prs {
					brandsWordForms[el.Brand], err = brandsRepo.GetWordForms(el.Brand)
					if err != nil {
						logtofile.Println(err)
						slacklogger.Error(err, "Get brand forms error!")
					}
				}
				lmProducts, err := linemediaRepo.GetList(el.Article, brandsWordForms[el.Brand], suppliers)
				if err != nil {
					logtofile.Println(err)
					slacklogger.Error(err, "Get linemedia part error!")
				}
				for _, part := range lmProducts {
					if _, prs := discounts[part.Supplier]; prs {
						part.Price = recalc(part, discounts[part.Supplier])
					}
					minPrice = math.Min(minPrice, part.Price)
				}

				if minPrice == math.MaxFloat64 {
					continue
				}

				minPrice = math.Ceil(minPrice)

				elementProperty, err := elementPropertyRepo.GetById(prop.Id, el.Id)
				if err != nil {
					logtofile.Println(err)
					slacklogger.Error(err, "Get element property error!")
				}
				if elementProperty == nil {
					elementProperty.Id, err = elementPropertyRepo.Add(&element.Property{
						IblockPropertyId: prop.Id,
						ElementId:        el.Id,
						Value:            minPrice,
					})
					if err != nil {
						logtofile.Println(err)
						slacklogger.Error(err, "Add min price error!")
					}
					cntadded++
				} else if elementProperty.Value != minPrice {
					logtofile.Println(el)
					logtofile.Println(minPrice)
					logtofile.Println(elementProperty.Value)
					_, err := elementPropertyRepo.Update(elementProperty.Id, minPrice)
					if err != nil {
						logtofile.Println(err)
						slacklogger.Error(err, "Update min price error!")
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
	slacklogger.Info(fmt.Sprintf("Lucy Finished! runned: %s", elapsed))
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
