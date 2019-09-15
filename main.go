package main

import (
	"savebestprice/iblock"
	"savebestprice/iblockProperty"
)

func main() {
	//http.HandleFunc("/", HelloServer)
	//http.ListenAndServe(":8080", nil)
	iblocks := iblock.GetIblockIds()
	for _, iblockId := range iblocks {
		var propId int64 = iblockProperty.AddProperty(iblockId, "CML2_BESTPRICE")
		RecalcBestPrice(propId, iblockId)
	}
}

func RecalcBestPrice(propId int64, iblockId int64) {

}

//func HelloServer(w http.ResponseWriter, r *http.Request) {
//	iblocks := model.GetList()
//	fmt.Fprintf(w, "%+v\n", iblocks)
//}
