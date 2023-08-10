package nosql

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"omo.msa.album/proxy"
	"time"
)

type PhotoStyle struct {
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
	// 类型：个人或者组织集体
	Type  uint8  `json:"type" bson:"type"`
	Cover string `json:"cover" bson:"cover"`
	//纸张的大小
	Size  uint8                 `json:"size" bson:"size"`
	Price uint32                `json:"price" bson:"price"`
	Slots []proxy.PhotocopySlot `json:"slots" bson:"slots"`
}

func CreatePhotoStyle(info *PhotoStyle) error {
	_, err := insertOne(TablePhotoStyle, info)
	if err != nil {
		return err
	}
	return nil
}

func GetPhotoStyleNextID() uint64 {
	num, _ := getSequenceNext(TablePhotoStyle)
	return num
}

func GetPhotoStyle(uid string) (*PhotoStyle, error) {
	result, err := findOne(TablePhotoStyle, uid)
	if err != nil {
		return nil, err
	}
	model := new(PhotoStyle)
	err1 := result.Decode(model)
	if err1 != nil {
		return nil, err1
	}
	return model, nil
}

func GetAllPhotoStyles() ([]*PhotoStyle, error) {
	cursor, err1 := findAllEnable(TablePhotoStyle, 0)
	if err1 != nil {
		return nil, err1
	}
	defer cursor.Close(context.Background())
	var items = make([]*PhotoStyle, 0, 100)
	for cursor.Next(context.Background()) {
		var node = new(PhotoStyle)
		if err := cursor.Decode(node); err != nil {
			return nil, err
		} else {
			items = append(items, node)
		}
	}
	return items, nil
}

func UpdatePhotoStyleBase(uid, name, remark, operator string, tp uint8) error {
	msg := bson.M{"name": name, "remark": remark, "type": tp, "operator": operator, TimeUpdated: time.Now().Unix()}
	_, err := updateOne(TablePhotoStyle, uid, msg)
	return err
}

func UpdatePhotoStyleCover(uid, cover, operator string) error {
	msg := bson.M{"cover": cover, "operator": operator, TimeUpdated: time.Now().Unix()}
	_, err := updateOne(TablePhotoStyle, uid, msg)
	return err
}

func UpdatePhotoStylePrice(uid, operator string, price uint32) error {
	msg := bson.M{"price": price, "operator": operator, TimeUpdated: time.Now().Unix()}
	_, err := updateOne(TablePhotoStyle, uid, msg)
	return err
}

func RemovePhotoStyle(uid, operator string) error {
	_, err := removeOne(TablePhotoStyle, uid, operator)
	return err
}

func AppendPhotoStylePage(uid string, slot proxy.PhotocopySlot) error {
	msg := bson.M{"slots": slot}
	_, err := appendElement(TablePhotoStyle, uid, msg)
	return err
}

func SubtractPhotoStylePage(uid string, index uint8) error {
	msg := bson.M{"slots": bson.M{"index": index}}
	_, err := removeElement(TablePhotoStyle, uid, msg)
	return err
}
