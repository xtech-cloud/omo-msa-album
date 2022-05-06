package cache

import (
	"omo.msa.album/config"
	"omo.msa.album/proxy/nosql"
	"reflect"
	"time"
)

type baseInfo struct {
	ID         uint64 `json:"-"`
	UID        string `json:"uid"`
	Name       string `json:"name"`
	Creator    string
	Operator   string
	CreateTime time.Time
	UpdateTime time.Time
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
	}
	return err
}

func Context() *cacheContext {
	return cacheCtx
}

func checkPage(page, number uint32, all interface{}) (uint32, uint32, interface{}) {
	if number < 1 {
		number = 10
	}
	array := reflect.ValueOf(all)
	total := uint32(array.Len())
	maxPage := total / number
	if total%number != 0 {
		maxPage = total/number + 1
	}
	if page < 1 {
		return total, maxPage, all
	}
	if page > maxPage {
		page = maxPage
	}

	var start = (page - 1) * number
	var end = start + number
	if end > total {
		end = total
	}

	list := array.Slice(int(start), int(end))
	return total, maxPage, list.Interface()
}
