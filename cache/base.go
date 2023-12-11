package cache

import (
	"bytes"
	"omo.msa.album/config"
	"omo.msa.album/proxy/nosql"
)

const DefaultOwner = "system"

const (
	PadLeft  PadType = 0
	PadRight PadType = 1
	PadBoth  PadType = 2
)

type PadType uint8

type baseInfo struct {
	ID       uint64 `json:"-"`
	UID      string `json:"uid"`
	Name     string `json:"name"`
	Creator  string
	Operator string
	Created  int64
	Updated  int64
	Deleted  int64
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
		//nosql.CheckTimes()
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

//StringPad tp: 0=left; 1=right; 2=both
func StringPad(input string, length int, pad string, tp PadType) string {
	var left, right = 0, 0
	chars := length - len(input)
	if chars <= 0 {
		return input
	}
	var buffer bytes.Buffer
	buffer.WriteString(input)
	switch tp {
	case PadLeft:
		left = chars
		right = 0
	case PadRight:
		left = 0
		right = chars
	case PadBoth:
		right = chars / 2
		left = chars - right
	}
	var leftBuffer bytes.Buffer
	for i := 0; i < left; i += 1 {
		leftBuffer.WriteString(pad)
		if leftBuffer.Len() > left {
			leftBuffer.Truncate(left)
			break
		}
	}
	for i := 0; i < right; i += 1 {
		buffer.WriteString(pad)
		if buffer.Len() > length {
			buffer.Truncate(length)
			break
		}
	}
	leftBuffer.WriteString(buffer.String())
	return leftBuffer.String()
}

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
