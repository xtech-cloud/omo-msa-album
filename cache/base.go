package cache

import (
	"omo.msa.album/config"
	"omo.msa.album/proxy/nosql"
)

const DefaultOwner = "system"

type baseInfo struct {
	ID       uint64 `json:"-"`
	UID      string `json:"uid"`
	Name     string `json:"name"`
	Creator  string
	Operator string
	Created  int64
	Updated  int64
}

type cacheContext struct {
}

var cacheCtx *cacheContext

func InitData() error {
	cacheCtx = &cacheContext{}

	err := nosql.InitDB(config.Schema.Database.IP, config.Schema.Database.Port, config.Schema.Database.Name, config.Schema.Database.Type)
	if err == nil {
		//num := nosql.GetAssetCount()
		//count := nosql.GetThumbCount()
		//logger.Infof("the asset count = %d and the thumb count = %d", num, count)
		//checkOwner()
		nosql.CheckTimes()
	}
	return err
}

func Context() *cacheContext {
	return cacheCtx
}

func checkOwner() {
	dbs, _ := nosql.GetAllPanoramas()
	for _, db := range dbs {
		if db.Owner == "" {
			_ = nosql.UpdatePanoramaOwner(db.UID.Hex(), DefaultOwner)
		}
	}
	dbs2, _ := nosql.GetAllExhibits()
	for _, db := range dbs2 {
		if db.Owner == "" {
			_ = nosql.UpdateExhibitOwner(db.UID.Hex(), DefaultOwner)
		}
	}
}

//func checkPage(page, number uint32, all interface{}) (uint32, uint32, interface{}) {
//	if number < 1 {
//		number = 10
//	}
//	array := reflect.ValueOf(all)
//	total := uint32(array.Len())
//	maxPage := total / number
//	if total%number != 0 {
//		maxPage = total/number + 1
//	}
//	if page < 1 {
//		return total, maxPage, all
//	}
//	if page > maxPage {
//		page = maxPage
//	}
//
//	var start = (page - 1) * number
//	var end = start + number
//	if end > total {
//		end = total
//	}
//
//	list := array.Slice(int(start), int(end))
//	return total, maxPage, list.Interface()
//}

func CheckPage[T any](page, number uint32, all []T) (uint32, uint32, []T) {
	if len(all) < 1 {
		return 0, 0, make([]T, 0, 1)
	}
	if number < 1 {
		number = 10
	}
	total := uint32(len(all))
	if len(all) <= int(number) {
		return total, 1, all
	}
	maxPage := total/number + 1
	if page < 1 {
		return total, maxPage, all
	}

	var start = (page - 1) * number
	var end = start + number
	if end > total {
		end = total
	}
	list := make([]T, 0, number)
	list = append(all[start:end])
	return total, maxPage, list
}
