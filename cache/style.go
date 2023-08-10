package cache

import (
	"errors"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"omo.msa.album/proxy"
	"omo.msa.album/proxy/nosql"
	"time"
)

const (
	StyleAll       StyleType = 0
	StyleForPerson StyleType = 1
	StyleForGroup  StyleType = 2
)

type StyleType uint8

// 影集样式或者模板
type PhotoStyleInfo struct {
	baseInfo
	Remark string
	Type   StyleType
	Cover  string
	Size   uint8
	Price  uint32
	Slots  []proxy.PhotocopySlot
	Tags   []string
}

func (mine *cacheContext) CreatePhotoTemplate(name, remark, user string, kind uint8) (*PhotoStyleInfo, error) {
	db := new(nosql.PhotoStyle)
	db.UID = primitive.NewObjectID()
	db.ID = nosql.GetAlbumNextID()
	db.Created = time.Now().Unix()
	db.Creator = user
	db.Name = name
	db.Remark = remark
	db.Type = kind
	db.Slots = make([]proxy.PhotocopySlot, 0, 1)
	err := nosql.CreatePhotoStyle(db)
	if err != nil {
		return nil, err
	}
	info := new(PhotoStyleInfo)
	info.initInfo(db)
	return info, nil
}

func (mine *cacheContext) GetPhotoTemplate(uid string) (*PhotoStyleInfo, error) {
	if len(uid) < 2 {
		return nil, errors.New("the PhotoTemplate uid is empty")
	}
	db, err := nosql.GetPhotoStyle(uid)
	if err != nil {
		return nil, err
	}
	info := new(PhotoStyleInfo)
	info.initInfo(db)
	return info, nil
}

func (mine *cacheContext) GetPhotoStyles(page, number uint32) (uint32, uint32, []*PhotoStyleInfo) {
	list := make([]*PhotoStyleInfo, 0, 20)
	array, err := nosql.GetAllPhotoStyles()
	if err != nil {
		return 0, 0, list
	}
	for _, item := range array {
		info := new(PhotoStyleInfo)
		info.initInfo(item)
		list = append(list, info)
	}
	if page < 1 {
		return uint32(len(list)), 0, list
	}
	if number < 1 {
		number = 10
	}
	total, maxPage, set := CheckPage(page, number, list)
	return total, maxPage, set
}

func (mine *PhotoStyleInfo) initInfo(db *nosql.PhotoStyle) {
	mine.Name = db.Name
	mine.UID = db.UID.Hex()
	mine.ID = db.ID
	mine.Remark = db.Remark
	mine.Created = db.Created
	mine.Operator = db.Operator
	mine.Creator = db.Creator
	mine.Size = db.Size
	mine.Cover = db.Cover
	mine.Type = StyleType(db.Type)

	mine.Slots = db.Slots
	if mine.Slots == nil {
		mine.Slots = make([]proxy.PhotocopySlot, 0, 1)
	}
}

func (mine *PhotoStyleInfo) UpdateBase(name, remark, operator string, tp uint8) error {
	err := nosql.UpdatePhotoStyleBase(mine.UID, name, remark, operator, tp)
	if err == nil {
		mine.Name = name
		mine.Remark = remark
		mine.Operator = operator
		mine.Type = StyleType(tp)
	}
	return err
}

func (mine *PhotoStyleInfo) Remove(operator string) error {
	return nosql.RemovePhotoStyle(mine.UID, operator)
}

func (mine *PhotoStyleInfo) getPage(index uint8) *proxy.PhotocopySlot {
	for _, page := range mine.Slots {
		if page.Index == index {
			return &page
		}
	}
	return nil
}

func (mine *PhotoStyleInfo) hadPage(index uint8) bool {
	for _, page := range mine.Slots {
		if page.Index == index {
			return true
		}
	}
	return false
}

/**
从相册中删除一个页面
*/
func (mine *PhotoStyleInfo) SubtractPage(index uint8) error {
	if !mine.hadPage(index) {
		return nil
	}
	err := nosql.SubtractPhotoStylePage(mine.UID, index)
	if err == nil {
		for i := 0; i < len(mine.Slots); i += 1 {
			if mine.Slots[i].Index == index {
				mine.Slots = append(mine.Slots[:i], mine.Slots[i+1:]...)
				break
			}
		}
	}
	return err
}

/**
向相册中添加一个页面
*/
func (mine *PhotoStyleInfo) AppendPage(page proxy.PhotocopySlot) error {
	if mine.hadPage(page.Index) {
		return nil
	}
	err := nosql.AppendPhotoStylePage(mine.UID, page)
	if err == nil {
		mine.Slots = append(mine.Slots, page)
	}
	return err
}
