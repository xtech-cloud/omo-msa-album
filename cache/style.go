package cache

import (
	"errors"
	"fmt"
	pb "github.com/xtech-cloud/omo-msp-album/proto/album"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"omo.msa.album/proxy"
	"omo.msa.album/proxy/nosql"
	"omo.msa.album/tool"
	"strconv"
	"time"
)

const (
	StyleAll       StyleType = 0
	StyleForPerson StyleType = 1
	StyleForGroup  StyleType = 2
)

type StyleType uint8

// 证书模板
type CertificateStyleInfo struct {
	baseInfo
	Count      uint32
	Year       int
	Remark     string
	Type       StyleType
	Cover      string
	Background string
	Prefix     string

	Width  int
	Height int

	Tags    []string
	Scenes  []string
	Relates []proxy.StyleRelate //关联的实体
	Slots   []proxy.StyleSlot
}

func (mine *cacheContext) CreateStyle(in *pb.ReqStyleAdd, slots []proxy.StyleSlot) (*CertificateStyleInfo, error) {
	db := new(nosql.CertificateStyle)
	db.UID = primitive.NewObjectID()
	db.ID = nosql.GetCertificateStyleNextID()
	db.Created = time.Now().Unix()
	db.CreatedTime = time.Now()
	db.Creator = in.Operator
	db.Name = in.Name
	db.Remark = in.Remark
	db.Type = uint8(in.Type)
	db.Tags = in.Tags
	db.Cover = in.Cover
	db.Prefix = in.Prefix
	db.Scenes = in.Scenes
	db.Background = in.Background
	db.Slots = slots
	db.Count = 0
	db.Width = int(in.Width)
	db.Height = int(in.Height)
	db.Year = time.Now().Year()
	db.Relates = make([]proxy.StyleRelate, 0, 1)
	if db.Slots == nil {
		db.Slots = make([]proxy.StyleSlot, 0, 1)
	}

	if db.Tags == nil {
		db.Tags = make([]string, 0, 1)
	}
	if db.Scenes == nil {
		db.Scenes = make([]string, 0, 1)
	}
	err := nosql.CreateCertificateStyle(db)
	if err != nil {
		return nil, err
	}
	info := new(CertificateStyleInfo)
	info.initInfo(db)
	return info, nil
}

func (mine *cacheContext) GetStyle(uid string) (*CertificateStyleInfo, error) {
	if len(uid) < 2 {
		return nil, errors.New("the PhotoTemplate uid is empty")
	}
	db, err := nosql.GetCertificateStyle(uid)
	if err != nil {
		return nil, err
	}
	info := new(CertificateStyleInfo)
	info.initInfo(db)
	return info, nil
}

func (mine *cacheContext) GetStyleByEntity(uid string) (*CertificateStyleInfo, error) {
	if len(uid) < 2 {
		return nil, errors.New("the entity uid is empty")
	}
	db, err := nosql.GetCertificateStyleByEntity(uid)
	if err != nil {
		return nil, err
	}
	info := new(CertificateStyleInfo)
	info.initInfo(db)
	return info, nil
}

func (mine *cacheContext) GetCertificatesCountByStyle(uid string) uint32 {
	if len(uid) < 2 {
		return 0
	}
	num, err := nosql.GetCertificatesCountByStyle(uid)
	if err != nil {
		return 0
	}
	return num
}

func (mine *cacheContext) GetCertificatesCountBySceneStyle(scene, uid string) uint32 {
	if len(uid) < 2 {
		return 0
	}
	num, err := nosql.GetCertificatesCountBySceneStyle(scene, uid)
	if err != nil {
		return 0
	}
	return num
}

func (mine *cacheContext) GetStyles(page, number uint32) (uint32, uint32, []*CertificateStyleInfo) {
	list := make([]*CertificateStyleInfo, 0, 20)
	array, err := nosql.GetAllCertificateStyles()
	if err != nil {
		return 0, 0, list
	}
	for _, item := range array {
		info := new(CertificateStyleInfo)
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

func (mine *cacheContext) GetStylesByScene(scene string) []*CertificateStyleInfo {
	list := make([]*CertificateStyleInfo, 0, 20)
	dbs, err := nosql.GetAllCertificateStyles()
	if err != nil {
		return list
	}
	for _, db := range dbs {
		if tool.HasItem(db.Scenes, scene) {
			info := new(CertificateStyleInfo)
			info.initInfo(db)
			list = append(list, info)
		}
	}
	return list
}

func (mine *cacheContext) GetStylesByArray(arr []string) []*CertificateStyleInfo {
	list := make([]*CertificateStyleInfo, 0, len(arr))
	for _, key := range arr {
		item, _ := mine.GetStyle(key)
		if item != nil {
			list = append(list, item)
		}
	}
	return list
}

func (mine *CertificateStyleInfo) initInfo(db *nosql.CertificateStyle) {
	mine.Name = db.Name
	mine.UID = db.UID.Hex()
	mine.ID = db.ID
	mine.Remark = db.Remark
	mine.Created = db.Created
	mine.Operator = db.Operator
	mine.Creator = db.Creator
	mine.Background = db.Background
	mine.Cover = db.Cover
	mine.Type = StyleType(db.Type)
	mine.Tags = db.Tags
	mine.Scenes = db.Scenes
	mine.Slots = db.Slots
	mine.Count = db.Count
	mine.Year = db.Year
	mine.Width = db.Width
	mine.Height = db.Height
	if mine.Slots == nil {
		mine.Slots = make([]proxy.StyleSlot, 0, 1)
	}
	mine.Relates = db.Relates
	if mine.Relates == nil {
		mine.Relates = make([]proxy.StyleRelate, 0, 1)
	}
}

func (mine *CertificateStyleInfo) UpdateBase(operator string, in *pb.ReqStyleUpdate) error {
	slots := make([]proxy.StyleSlot, 0, len(in.Slots))
	for _, slot := range in.Slots {
		slots = append(slots, proxy.StyleSlot{
			Key:    slot.Key,
			X:      slot.X,
			Y:      slot.Y,
			Width:  slot.Width,
			Height: slot.Height,
			Bold:   slot.Bold,
			Size:   slot.Size,
		})
	}
	err := nosql.UpdateCertificateStyleBase(mine.UID, operator, in.Name, in.Remark, in.Background, int(in.Width), int(in.Height), in.Tags, in.Scenes, slots, uint8(in.Type))
	if err == nil {
		mine.Name = in.Name
		mine.Remark = in.Remark
		mine.Operator = operator
		mine.Background = in.Background
		mine.Tags = in.Tags
		mine.Scenes = in.Scenes
		mine.Slots = slots
		mine.Width = int(in.Width)
		mine.Height = int(in.Height)
		mine.Type = StyleType(in.Type)
	}
	return err
}

func (mine *CertificateStyleInfo) updateCount(operator string, num uint32) error {
	year := time.Now().Year()
	if mine.Year != year {
		num = 1
	}
	err := nosql.UpdateCertificateStyleCount(mine.UID, operator, year, int(num))
	if err == nil {
		mine.Count = num
		mine.Year = year
		mine.Operator = operator
	}
	return err
}

func (mine *CertificateStyleInfo) UpdateCover(cover, operator string) error {
	if mine.Cover == cover {
		return nil
	}
	err := nosql.UpdateCertificateStyleCover(mine.UID, cover, operator)
	if err == nil {
		mine.Cover = cover
		mine.Operator = operator
	}
	return err
}

func (mine *CertificateStyleInfo) AppendRelate(entity, operator string, way uint32) error {
	arr := make([]proxy.StyleRelate, 0, len(mine.Relates)+1)
	arr = append(arr, mine.Relates...)
	add := true
	for i, relate := range arr {
		if relate.Entity == entity {
			if way == 0 {
				add = false
				if i == len(arr)-1 {
					arr = append(arr[:i])
				} else {
					arr = append(arr[:i], arr[i+1:]...)
				}
			} else {
				if uint32(relate.Way) == way {
					return nil
				}
				add = false
				relate.Way = uint8(way)
			}
			break
		}
	}
	if add {
		arr = append(arr, proxy.StyleRelate{Entity: entity, Way: uint8(way)})
	}
	err := nosql.UpdateCertificateStyleRelates(mine.UID, operator, arr)
	if err == nil {
		mine.Relates = arr
		mine.Operator = operator
	}
	return err
}

func (mine *CertificateStyleInfo) GetSN(operator string) string {
	err := mine.updateCount(operator, mine.Count+1)
	if err != nil {
		return ""
	}
	tp := StringPad(strconv.Itoa(int(mine.Type)), 3, "0", PadLeft)
	s := StringPad(strconv.Itoa(int(mine.Count)), 5, "0", PadLeft)
	prefix := mine.Prefix
	if prefix == "" {
		prefix = "SCCD"
	}
	return fmt.Sprintf("%s-%s-%d-%s", prefix, tp, mine.Year, s)
}

func (mine *CertificateStyleInfo) Remove(operator string) error {
	return nosql.RemoveCertificateStyle(mine.UID, operator)
}

func (mine *CertificateStyleInfo) CreateCertificate(scene, quote, operator string, end int64) (*CertificateInfo, error) {
	db := new(nosql.Certificate)
	db.UID = primitive.NewObjectID()
	db.ID = nosql.GetCertificateNextID()
	db.Created = time.Now().Unix()
	db.Creator = operator
	db.Name = mine.Name
	db.Remark = mine.Remark
	db.Scene = scene
	db.Image = ""
	db.Target = ""
	db.Quote = quote
	db.SN = mine.GetSN(operator)
	db.Status = 0
	db.Type = uint8(mine.Type)
	db.Tags = mine.Tags
	db.Style = mine.UID
	db.EndStamp = end
	db.Assets = make([]string, 0, 1)
	db.Contact = nil
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

func (mine *CertificateStyleInfo) Batch(scene, quote, operator string, num uint32, end int64) []*CertificateInfo {
	list := make([]*CertificateInfo, 0, num)
	for i := 0; i < int(num); i += 1 {
		info, er := mine.CreateCertificate(scene, quote, operator, end)
		if er == nil {
			list = append(list, info)
		}
	}
	return list
}
