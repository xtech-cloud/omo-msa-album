package cache

import (
	"errors"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"omo.msa.album/proxy/nosql"
	"time"
)

// 文件夹或者包
type CertificateInfo struct {
	baseInfo
	Remark string
	Scene  string
	Target string
	Style  string
	Status uint8
	Type   uint8
	SN     string
	Image  string
	Tags   []string
}

func (mine *cacheContext) CreateCertificate(name, remark, user, sn, scene, img, target, style string, tp uint8, tags []string) (*CertificateInfo, error) {
	db := new(nosql.Certificate)
	db.UID = primitive.NewObjectID()
	db.ID = nosql.GetCertificateNextID()
	db.Created = time.Now().Unix()
	db.Creator = user
	db.Name = name
	db.Remark = remark
	db.Scene = scene
	db.Image = img
	db.Target = target
	db.SN = sn
	db.Status = 0
	db.Type = tp
	db.Tags = tags
	db.Style = style
	if db.Tags == nil {
		db.Tags = make([]string, 0, 1)
	}

	err := nosql.CreateCertificate(db)
	if err != nil {
		return nil, err
	}
	info := new(CertificateInfo)
	info.initInfo(db)
	return info, nil
}

func (mine *cacheContext) GetCertificate(uid string) (*CertificateInfo, error) {
	if len(uid) < 2 {
		return nil, errors.New("the album uid is empty")
	}
	db, err := nosql.GetCertificate(uid)
	if err != nil {
		return nil, err
	}
	info := new(CertificateInfo)
	info.initInfo(db)
	return info, nil
}

func (mine *cacheContext) GetCertificatesByScene(uid string) []*CertificateInfo {
	list := make([]*CertificateInfo, 0, 20)
	if len(uid) < 2 {
		uid = DefaultOwner
	}
	array, err := nosql.GetCertificatesByScene(uid)
	if err != nil {
		return list
	}
	for _, item := range array {
		info := new(CertificateInfo)
		info.initInfo(item)
		list = append(list, info)
	}
	return list
}

func (mine *cacheContext) GetCertificatesByStyle(uid string) []*CertificateInfo {
	list := make([]*CertificateInfo, 0, 20)
	if len(uid) < 2 {
		return nil
	}
	array, err := nosql.GetCertificatesByStyle(uid)
	if err != nil {
		return list
	}
	for _, item := range array {
		info := new(CertificateInfo)
		info.initInfo(item)
		list = append(list, info)
	}
	return list
}

func (mine *cacheContext) GetCertificatesByTarget(uid string) []*CertificateInfo {
	list := make([]*CertificateInfo, 0, 20)
	if len(uid) < 2 {
		return nil
	}
	array, err := nosql.GetCertificatesByTarget(uid)
	if err != nil {
		return list
	}
	for _, item := range array {
		info := new(CertificateInfo)
		info.initInfo(item)
		list = append(list, info)
	}
	return list
}

func (mine *CertificateInfo) initInfo(db *nosql.Certificate) {
	mine.Name = db.Name
	mine.UID = db.UID.Hex()
	mine.ID = db.ID
	mine.Created = db.Created
	mine.Creator = db.Creator
	mine.Operator = db.Operator
	mine.Updated = db.Updated

	mine.Name = db.Name
	mine.Remark = db.Remark
	mine.Scene = db.Scene
	mine.Image = db.Image
	mine.Target = db.Target
	mine.Type = db.Type
	mine.Status = db.Status
	mine.SN = db.SN
	mine.Style = db.Style
	mine.Tags = db.Tags
}

func (mine *CertificateInfo) UpdateBase(name, remark, operator string, tags []string) error {
	err := nosql.UpdateCertificateBase(mine.UID, name, remark, operator, tags)
	if err == nil {
		mine.Name = name
		mine.Remark = remark
		mine.Operator = operator
		mine.Tags = tags
		mine.Updated = time.Now().Unix()
	}
	return err
}

func (mine *CertificateInfo) Remove(operator string) error {
	return nosql.RemoveCertificate(mine.UID, operator)
}
