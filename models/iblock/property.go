package iblock

import (
	"database/sql"
	"fmt"
	"log"
	"savebestprice/config"
)

type Property struct {
	Active string
	Sort int
	Code string
	PropertyType string
	IblockId int64
	Name string
	Multiple string
	MultipleCnt int
	IsRequired string
}
