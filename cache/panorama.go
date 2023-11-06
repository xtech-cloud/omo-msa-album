package cache

import (
	"errors"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"omo.msa.album/proxy/nosql"
	"time"
)

type PanoramaInfo struct {
	baseInfo
	Remark  string
	Content string
	Owner   string
}

func (mine *cacheContext) CreatePanorama(name, remark, json, owner, operator string) (*PanoramaInfo, error) {
	db := new(nosql.Panorama)
	db.UID = primitive.NewObjectID()
	db.ID = nosql.GetPanoramaNextID()
	db.CreatedTime = time.Now()
	db.Created = time.Now().Unix()
	db.Creator = operator
	db.Name = name
	db.Remark = remark
	db.Content = json
	if owner == "" {
		owner = DefaultOwner
	}
	db.Owner = owner
	err := nosql.CreatePanorama(db)
	if err != nil {
		return nil, err
	}
	info := new(PanoramaInfo)
	info.initInfo(db)
	return info, nil
}

func (mine *cacheContext) GetPanorama(uid string) (*PanoramaInfo, error) {
	db, err := nosql.GetPanorama(uid)
	if err != nil {
		return nil, err
	}
	info := new(PanoramaInfo)
	info.initInfo(db)
	return info, nil
}

func (mine *cacheContext) RemovePanorama(uid, operator string) error {
	if uid == "" {
		return errors.New("the panorama uid is empty")
	}
	err := nosql.RemovePanorama(uid, operator)
	return err
}

func (mine *cacheContext) GetPanoramas(owner string) ([]*PanoramaInfo, error) {
	if owner == "" {
		owner = DefaultOwner
	}
	array, err := nosql.GetAllPanoramasByOwner(owner)
	if err != nil {
		return make([]*PanoramaInfo, 0, 1), err
	}
	list := make([]*PanoramaInfo, 0, len(array))
	for _, item := range array {
		info := new(PanoramaInfo)
		info.initInfo(item)
		list = append(list, info)
	}
	return list, nil
}

func (mine *PanoramaInfo) initInfo(db *nosql.Panorama) {
	mine.UID = db.UID.Hex()
	mine.Name = db.Name
	mine.Created = db.Created
	mine.Remark = db.Remark
	mine.Content = db.Content
}

func (mine *PanoramaInfo) UpdateContent(data, operator string) error {
	err := nosql.UpdatePanoramaContent(mine.UID, data)
	if err == nil {
		mine.Content = data
		mine.Operator = operator
	}
	return err
}
