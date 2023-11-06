package cache

import (
	"errors"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"omo.msa.album/proxy/nosql"
	"time"
)

// 影集样式或者模板
type PhotoFrameInfo struct {
	baseInfo
	Remark string
	Asset  string
	Owner  string
	Width  uint32
	Height uint32
}

func (mine *cacheContext) CreatePhotoFrame(name, remark, asset, owner, user string, width, height uint32) (*PhotoFrameInfo, error) {
	db := new(nosql.PhotoFrame)
	db.UID = primitive.NewObjectID()
	db.ID = nosql.GetPhotoFrameNextID()
	db.CreatedTime = time.Now()
	db.Created = time.Now().Unix()
	db.Creator = user
	db.Name = name
	db.Remark = remark
	db.Type = 0
	db.Asset = asset
	db.Width = width
	db.Height = height
	db.Owner = owner
	err := nosql.CreatePhotoFrame(db)
	if err != nil {
		return nil, err
	}
	info := new(PhotoFrameInfo)
	info.initInfo(db)
	return info, nil
}

func (mine *cacheContext) GetPhotoFrame(uid string) (*PhotoFrameInfo, error) {
	if len(uid) < 2 {
		return nil, errors.New("the PhotoTemplate uid is empty")
	}
	db, err := nosql.GetPhotoFrame(uid)
	if err != nil {
		return nil, err
	}
	info := new(PhotoFrameInfo)
	info.initInfo(db)
	return info, nil
}

func (mine *cacheContext) GetPhotoFrames(owner string) []*PhotoFrameInfo {
	list := make([]*PhotoFrameInfo, 0, 20)
	array, err := nosql.GetPhotoFramesByOwner(owner)
	if err != nil {
		return list
	}
	for _, item := range array {
		info := new(PhotoFrameInfo)
		info.initInfo(item)
		list = append(list, info)
	}
	return list
}

func (mine *PhotoFrameInfo) initInfo(db *nosql.PhotoFrame) {
	mine.Name = db.Name
	mine.UID = db.UID.Hex()
	mine.ID = db.ID
	mine.Remark = db.Remark
	mine.Created = db.Created
	mine.Updated = db.Updated
	mine.Creator = db.Creator
	mine.Operator = db.Operator
	mine.Width = db.Width
	mine.Height = db.Height
	mine.Asset = db.Asset
	mine.Owner = db.Owner
}

func (mine *PhotoFrameInfo) UpdateBase(name, remark, operator string) error {
	err := nosql.UpdatePhotoFrameBase(mine.UID, name, remark, operator)
	if err == nil {
		mine.Name = name
		mine.Remark = remark
		mine.Operator = operator
	}
	return err
}

func (mine *PhotoFrameInfo) Remove(operator string) error {
	return nosql.RemovePhotoFrame(mine.UID, operator)
}
