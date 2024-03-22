package cache

import (
	"errors"
	pb "github.com/xtech-cloud/omo-msp-album/proto/album"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"omo.msa.album/proxy"
	"omo.msa.album/proxy/nosql"
	"time"
)

// 文件夹或者包
type CertificateInfo struct {
	Status uint8
	Type   uint8
	baseInfo
	Remark string
	Scene  string
	Target string //目标对象
	Style  string
	SN     string
	Image  string
	Quote  string //引用的实体

	EndStamp int64 //截至时间
	Contact  *proxy.ContactInfo
	Tags     []string
	Assets   []string
}

func (mine *cacheContext) CreateCertificate(in *pb.ReqCertificateAdd) (*CertificateInfo, error) {
	db := new(nosql.Certificate)
	db.UID = primitive.NewObjectID()
	db.ID = nosql.GetCertificateNextID()
	db.Created = time.Now().Unix()
	db.Creator = in.Operator
	db.Name = in.Name
	db.Remark = in.Remark
	db.Scene = in.Scene
	db.Image = in.Image
	db.Target = in.Target
	db.SN = in.Sn
	db.Status = 0
	db.Type = uint8(in.Type)
	db.Tags = in.Tags
	db.Style = in.Style
	db.Assets = in.Assets
	db.Contact = &proxy.ContactInfo{
		Name:    in.Contact.Name,
		Phone:   in.Contact.Phone,
		Address: in.Contact.Address,
		Remark:  in.Contact.Remark,
	}
	if db.Tags == nil {
		db.Tags = make([]string, 0, 1)
	}
	if db.Assets == nil {
		db.Assets = make([]string, 0, 1)
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

func (mine *cacheContext) GetCertificatesByStyle(scene, uid string) []*CertificateInfo {
	list := make([]*CertificateInfo, 0, 20)
	if len(uid) < 2 {
		return nil
	}
	var dbs []*nosql.Certificate
	var err error
	if len(scene) > 2 {
		dbs, err = nosql.GetCertificatesBySceneStyle(scene, uid)
	} else {
		dbs, err = nosql.GetCertificatesByStyle(uid)
	}

	if err != nil {
		return list
	}
	for _, item := range dbs {
		info := new(CertificateInfo)
		info.initInfo(item)
		list = append(list, info)
	}
	return list
}

func (mine *cacheContext) GetCertificateByContact(phone, name string) []*CertificateInfo {
	list := make([]*CertificateInfo, 0, 5)
	if len(phone) < 2 || len(name) < 2 {
		return list
	}
	dbs, err := nosql.GetCertificatesByContact(phone, name)
	if err != nil {
		return list
	}
	for _, item := range dbs {
		if len(item.Target) < 2 {
			info := new(CertificateInfo)
			info.initInfo(item)
			list = append(list, info)
		}
	}
	return list
}

func (mine *cacheContext) GetCertificateByArray(arr []string) []*CertificateInfo {
	list := make([]*CertificateInfo, 0, len(arr))
	for _, key := range arr {
		info, _ := mine.GetCertificate(key)
		if info != nil {
			list = append(list, info)
		}
	}
	return list
}

func (mine *cacheContext) BatchCertificate(uid, scene, quote, operator string, num uint32, end int64) ([]*CertificateInfo, error) {
	info, err := mine.GetStyle(uid)
	if err != nil {
		return nil, err
	}
	return info.Batch(scene, quote, operator, num, end), nil
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
	mine.Quote = db.Quote
	mine.EndStamp = db.EndStamp
	mine.Contact = db.Contact
	mine.Style = db.Style
	mine.Tags = db.Tags
	mine.Assets = db.Assets
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

func (mine *CertificateInfo) UpdateContact(name, phone, addr, remark, operator string) error {
	if len(mine.Target) > 2 {
		return errors.New("the certificate had send one target")
	}
	contact := &proxy.ContactInfo{
		Name:    name,
		Phone:   phone,
		Remark:  remark,
		Address: addr,
	}
	err := nosql.UpdateCertificateContact(mine.UID, operator, contact)
	if err == nil {
		mine.Operator = operator
		mine.Contact = contact
		mine.Updated = time.Now().Unix()
	}
	return err
}

func (mine *CertificateInfo) UpdateStatus(operator string, st uint8) error {
	if mine.Status == st {
		return nil
	}
	err := nosql.UpdateCertificateStatus(mine.UID, operator, st)
	if err == nil {
		mine.Status = st
		mine.Operator = operator
	}
	return err
}

func (mine *CertificateInfo) UpdateTarget(operator, target string) error {
	if mine.Target == target {
		return nil
	}
	err := nosql.UpdateCertificateTarget(mine.UID, target, operator)
	if err == nil {
		mine.Target = target
		mine.Operator = operator
	}
	return err
}

func (mine *CertificateInfo) Remove(operator string) error {
	return nosql.RemoveCertificate(mine.UID, operator)
}
