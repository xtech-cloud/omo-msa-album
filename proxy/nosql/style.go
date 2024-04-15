package nosql

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"omo.msa.album/proxy"
	"time"
)

type CertificateStyle struct {
	UID         primitive.ObjectID `bson:"_id"`
	ID          uint64             `json:"id" bson:"id"`
	CreatedTime time.Time          `json:"createdAt" bson:"createdAt"`
	UpdatedTime time.Time          `json:"updatedAt" bson:"updatedAt"`
	DeleteTime  time.Time          `json:"deleteAt" bson:"deleteAt"`
	Created     int64              `json:"created" bson:"created"`
	Updated     int64              `json:"updated" bson:"updated"`
	Deleted     int64              `json:"deleted" bson:"deleted"`
	Creator     string             `json:"creator" bson:"creator"`
	Operator    string             `json:"operator" bson:"operator"`

	Name   string `json:"name" bson:"name"`
	Remark string `json:"remark" bson:"remark"`
	Type   uint8  `json:"type" bson:"type"`
	Cover  string `json:"cover" bson:"cover"`
	Prefix string `json:"prefix" bson:"prefix"`
	Count  uint32 `json:"count" bson:"count"`
	Year   int    `json:"year" bson:"year"`
	Width  int    `json:"width" bson:"width"`
	Height int    `json:"height" bson:"height"`

	Background string              `json:"background" bson:"background"`
	Tags       []string            `json:"tags" bson:"tags"`
	Scenes     []string            `json:"scenes" bson:"scenes"`
	Slots      []proxy.StyleSlot   `json:"slots" bson:"slots"`
	Relates    []proxy.StyleRelate `json:"relates" bson:"relates"`
}

func CreateCertificateStyle(info *CertificateStyle) error {
	_, err := insertOne(TableCertificateStyle, info)
	if err != nil {
		return err
	}
	return nil
}

func GetCertificateStyleNextID() uint64 {
	num, _ := getSequenceNext(TableCertificateStyle)
	return num
}

func GetCertificateStyle(uid string) (*CertificateStyle, error) {
	result, err := findOne(TableCertificateStyle, uid)
	if err != nil {
		return nil, err
	}
	model := new(CertificateStyle)
	err1 := result.Decode(model)
	if err1 != nil {
		return nil, err1
	}
	return model, nil
}

func GetCertificateStyleByEntity(uid string) (*CertificateStyle, error) {
	filter := bson.M{"relates.entity": uid, TimeDeleted: 0}
	result, err := findOneBy(TableCertificateStyle, filter)
	if err != nil {
		return nil, err
	}
	model := new(CertificateStyle)
	err1 := result.Decode(model)
	if err1 != nil {
		return nil, err1
	}
	return model, nil
}

func GetAllCertificateStyles() ([]*CertificateStyle, error) {
	cursor, err1 := findAllEnable(TableCertificateStyle, 0)
	if err1 != nil {
		return nil, err1
	}
	defer cursor.Close(context.Background())
	var items = make([]*CertificateStyle, 0, 100)
	for cursor.Next(context.Background()) {
		var node = new(CertificateStyle)
		if err := cursor.Decode(node); err != nil {
			return nil, err
		} else {
			items = append(items, node)
		}
	}
	return items, nil
}

func UpdateCertificateStyleBase(uid, operator, name, remark, bg string, w, h int, tags, scenes []string, slots []proxy.StyleSlot, tp uint8) error {
	msg := bson.M{"name": name, "remark": remark, "background": bg, "tags": tags, "width": w, "height": h,
		"scenes": scenes, "slots": slots, "type": tp, "operator": operator, TimeUpdated: time.Now().Unix()}
	_, err := updateOne(TableCertificateStyle, uid, msg)
	return err
}

func UpdateCertificateStyleCount(uid, operator string, year, num int) error {
	msg := bson.M{"count": num, "year": year, "operator": operator, TimeUpdated: time.Now().Unix()}
	_, err := updateOne(TableCertificateStyle, uid, msg)
	return err
}

func UpdateCertificateStyleCover(uid, cover, operator string) error {
	msg := bson.M{"cover": cover, "operator": operator, TimeUpdated: time.Now().Unix()}
	_, err := updateOne(TableCertificateStyle, uid, msg)
	return err
}

func UpdateCertificateStyleSlots(uid, operator string, slots []proxy.StyleSlot) error {
	msg := bson.M{"slots": slots, "operator": operator, TimeUpdated: time.Now().Unix()}
	_, err := updateOne(TableCertificateStyle, uid, msg)
	return err
}

func UpdateCertificateStyleRelates(uid, operator string, arr []proxy.StyleRelate) error {
	msg := bson.M{"relates": arr, "operator": operator, TimeUpdated: time.Now().Unix()}
	_, err := updateOne(TableCertificateStyle, uid, msg)
	return err
}

func UpdateCertificateStyleScenes(uid, operator string, scenes []string) error {
	msg := bson.M{"scenes": scenes, "operator": operator, TimeUpdated: time.Now().Unix()}
	_, err := updateOne(TableCertificateStyle, uid, msg)
	return err
}

func RemoveCertificateStyle(uid, operator string) error {
	_, err := removeOne(TableCertificateStyle, uid, operator)
	return err
}
