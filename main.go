package main

import (
	"savebestprice/config"
	"savebestprice/iblock/repository"
	"savebestprice/models/iblock"
)

func main() {
	//http.HandleFunc("/", HelloServer)
	//http.ListenAndServe(":8080", nil)
	iblockRepo := repository.NewIblockRepository(config.DB)
	iblocks, _ := iblockRepo.GetIblockIds()
	for _, iblockId := range iblocks {
		var propId int64 = iblock.property.Add(iblockId, "CML2_BESTPRICE")
		RecalcBestPrice(propId, iblockId)
	}
}

func RecalcBestPrice(propId int64, iblockId int64) {
}

//func HelloServer(w http.ResponseWriter, r *http.Request) {
//	iblocks := models.GetList()
//	fmt.Fprintf(w, "%+v\n", iblocks)
//}
