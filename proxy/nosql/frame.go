package nosql

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type PhotoFrame struct {
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
	Owner  string `json:"owner" bson:"owner"`
	// 类型：个人或者组织集体
	Type  uint8  `json:"type" bson:"type"`
	Asset string `json:"asset" bson:"asset"`
	//纸张的大小
	Width  uint32 `json:"width" bson:"width"`
	Height uint32 `json:"height" bson:"height"`
}

func CreatePhotoFrame(info *PhotoFrame) error {
	_, err := insertOne(TablePhotoFrame, info)
	if err != nil {
		return err
	}
	return nil
}

func GetPhotoFrameNextID() uint64 {
	num, _ := getSequenceNext(TablePhotoFrame)
	return num
}

func GetPhotoFrame(uid string) (*PhotoFrame, error) {
	result, err := findOne(TablePhotoFrame, uid)
	if err != nil {
		return nil, err
	}
	model := new(PhotoFrame)
	err1 := result.Decode(model)
	if err1 != nil {
		return nil, err1
	}
	return model, nil
}

func GetPhotoFramesByOwner(owner string) ([]*PhotoFrame, error) {
	msg := bson.M{"owner": owner, TimeDeleted: 0}
	cursor, err1 := findMany(TablePhotoFrame, msg, 0)
	if err1 != nil {
		return nil, err1
	}
	defer cursor.Close(context.Background())
	var items = make([]*PhotoFrame, 0, 100)
	for cursor.Next(context.Background()) {
		var node = new(PhotoFrame)
		if err := cursor.Decode(node); err != nil {
			return nil, err
		} else {
			items = append(items, node)
		}
	}
	return items, nil
}

func UpdatePhotoFrameBase(uid, name, remark, operator string) error {
	msg := bson.M{"name": name, "remark": remark, "operator": operator, TimeUpdated: time.Now().Unix()}
	_, err := updateOne(TablePhotoFrame, uid, msg)
	return err
}

func UpdatePhotoFrameCover(uid, asset, operator string, width, height uint32) error {
	msg := bson.M{"asset": asset, "width": width, "height": height, "operator": operator, TimeUpdated: time.Now().Unix()}
	_, err := updateOne(TablePhotoFrame, uid, msg)
	return err
}

func RemovePhotoFrame(uid, operator string) error {
	_, err := removeOne(TablePhotoFrame, uid, operator)
	return err
}
