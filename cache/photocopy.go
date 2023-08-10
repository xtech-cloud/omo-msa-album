package cache

import (
	"errors"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"omo.msa.album/proxy"
	"omo.msa.album/proxy/nosql"
	"time"
)

// 影集
type PhotocopyInfo struct {
	baseInfo
	Mother   string
	Remark   string
	Template string
	Owner    string
	Count    uint32
	Tags     []string
	Slots    []proxy.PhotocopySlot
}

func (mine *cacheContext) CreatePhotocopy(name, remark, user, template, owner string) (*PhotocopyInfo, error) {
	db := new(nosql.Photocopy)
	db.UID = primitive.NewObjectID()
	db.ID = nosql.GetAlbumNextID()
	db.Created = time.Now().Unix()
	db.Creator = user
	db.Name = name
	db.Remark = remark
	db.Count = 0
	db.Mother = ""
	db.Template = template
	db.Owner = owner
	db.Slots = make([]proxy.PhotocopySlot, 0, 1)
	err := nosql.CreatePhotocopy(db)
	if err != nil {
		return nil, err
	}
	info := new(PhotocopyInfo)
	info.initInfo(db)
	return info, nil
}

func (mine *cacheContext) ClonePhotocopy(uid, user string) (*PhotocopyInfo, error) {
	master, err1 := mine.GetPhotocopy(uid)
	if err1 != nil {
		return nil, err1
	}
	db := new(nosql.Photocopy)
	db.UID = primitive.NewObjectID()
	db.ID = nosql.GetAlbumNextID()
	db.CreatedTime = time.Now()
	db.Creator = user
	db.Name = master.Name
	db.Remark = master.Remark
	db.Count = 0
	db.Owner = DefaultOwner
	db.Mother = master.UID
	db.Template = master.Template
	db.Slots = make([]proxy.PhotocopySlot, 0, 1)
	err := nosql.CreatePhotocopy(db)
	if err != nil {
		return nil, err
	}
	_ = master.increaseCopy()
	info := new(PhotocopyInfo)
	info.initInfo(db)
	return info, nil
}

func (mine *cacheContext) GetPhotocopy(uid string) (*PhotocopyInfo, error) {
	if len(uid) < 2 {
		return nil, errors.New("the photocopy uid is empty")
	}
	db, err := nosql.GetPhotocopy(uid)
	if err != nil {
		return nil, err
	}
	info := new(PhotocopyInfo)
	info.initInfo(db)
	return info, nil
}

func (mine *cacheContext) GetPhotocopiesByCreator(user string) ([]*PhotocopyInfo, error) {
	list := make([]*PhotocopyInfo, 0, 20)
	array, err := nosql.GetPhotocopiesByCreator(user)
	if err != nil {
		return list, err
	}
	for _, item := range array {
		info := new(PhotocopyInfo)
		info.initInfo(item)
		list = append(list, info)
	}
	return list, nil
}

func (mine *cacheContext) GetPhotocopiesByMaster(uid string) ([]*PhotocopyInfo, error) {
	list := make([]*PhotocopyInfo, 0, 20)
	array, err := nosql.GetPhotocopiesByMaster(uid)
	if err != nil {
		return list, err
	}
	for _, item := range array {
		info := new(PhotocopyInfo)
		info.initInfo(item)
		list = append(list, info)
	}
	return list, nil
}

func (mine *cacheContext) GetPhotocopiesByTemplate(uid string) ([]*PhotocopyInfo, error) {
	list := make([]*PhotocopyInfo, 0, 20)
	array, err := nosql.GetPhotocopiesByTemplate(uid)
	if err != nil {
		return list, err
	}
	for _, item := range array {
		info := new(PhotocopyInfo)
		info.initInfo(item)
		list = append(list, info)
	}
	return list, nil
}

func (mine *PhotocopyInfo) initInfo(db *nosql.Photocopy) {
	mine.Name = db.Name
	mine.UID = db.UID.Hex()
	mine.ID = db.ID
	mine.Remark = db.Remark
	mine.Created = db.Created
	mine.Updated = db.Updated
	mine.Creator = db.Creator
	mine.Template = db.Template
	mine.Count = db.Count
	mine.Mother = db.Mother
	mine.Tags = db.Tags
	mine.Slots = db.Slots
	if mine.Slots == nil {
		mine.Slots = make([]proxy.PhotocopySlot, 0, 1)
	}
}

func (mine *PhotocopyInfo) GetPages() []proxy.PhotocopySlot {
	return mine.Slots
}

func (mine *PhotocopyInfo) UpdateBase(name, remark, operator string) error {
	err := nosql.UpdatePhotocopyBase(mine.UID, name, remark, operator)
	if err == nil {
		mine.Name = name
		mine.Remark = remark
		mine.Operator = operator
	}
	return err
}

func (mine *PhotocopyInfo) Remove(operator string) error {
	return nosql.RemovePhotocopy(mine.UID, operator)
}

func (mine *PhotocopyInfo) increaseCopy() error {
	sum := mine.Count + 1
	err := nosql.UpdatePhotocopyCount(mine.UID, sum)
	if err == nil {
		mine.Count = sum
	}
	return err
}

func (mine *PhotocopyInfo) getPage(index uint8) *proxy.PhotocopySlot {
	for _, page := range mine.Slots {
		if page.Index == index {
			return &page
		}
	}
	return nil
}

func (mine *PhotocopyInfo) hadPage(index uint8) bool {
	for _, page := range mine.Slots {
		if page.Index == index {
			return true
		}
	}
	return false
}

func (mine *PhotocopyInfo) UpdatePage(info proxy.PhotocopySlot) error {
	if info.Index < 1 {
		return errors.New("the photocopy page is empty")
	}
	page := mine.getPage(info.Index)
	if page == nil {
		return errors.New("the photocopy page not existed")
	}

	err := mine.subtractPage(info.Index)
	if err != nil {
		return err
	}
	return mine.appendPage(info)
}

/**
从相册中删除一个页面
*/
func (mine *PhotocopyInfo) subtractPage(index uint8) error {
	if !mine.hadPage(index) {
		return nil
	}
	err := nosql.SubtractPhotocopyPage(mine.UID, index)
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
func (mine *PhotocopyInfo) appendPage(page proxy.PhotocopySlot) error {
	if mine.hadPage(page.Index) {
		return nil
	}
	err := nosql.AppendPhotocopyPage(mine.UID, page)
	if err == nil {
		mine.Slots = append(mine.Slots, page)
	}
	return err
}
