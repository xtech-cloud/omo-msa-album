package nosql

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"omo.msa.album/proxy"
	"time"
)

//影集
type Photocopy struct {
	UID         primitive.ObjectID `bson:"_id"`
	ID          uint64             `json:"id" bson:"id"`
	CreatedTime time.Time          `json:"createdAt" bson:"createdAt"`
	UpdatedTime time.Time          `json:"updatedAt" bson:"updatedAt"`
	DeleteTime  time.Time          `json:"deleteAt" bson:"deleteAt"`
	Creator     string             `json:"creator" bson:"creator"`
	Operator    string             `json:"operator" bson:"operator"`

	// 克隆的对象或者母版
	Mother string `json:"mother" bson:"mother"`
	Name   string `json:"name" bson:"name"`
	Remark string `json:"remark" bson:"remark"`
	// 样式模板
	Template string `json:"template" bson:"template"`
	// 被克隆的数量
	Count uint32                `json:"count" bson:"count"`
	Owner string                `json:"owner" bson:"owner"`
	Tags  []string              `json:"tags" bson:"tags"`
	Slots []proxy.PhotocopySlot `json:"slots" bson:"slots"`
}

func CreatePhotocopy(info *Photocopy) error {
	_, err := insertOne(TablePhotocopy, info)
	if err != nil {
		return err
	}
	return nil
}

func GetPhotocopyNextID() uint64 {
	num, _ := getSequenceNext(TablePhotocopy)
	return num
}

func GetPhotocopy(uid string) (*Photocopy, error) {
	result, err := findOne(TablePhotocopy, uid)
	if err != nil {
		return nil, err
	}
	model := new(Photocopy)
	err1 := result.Decode(model)
	if err1 != nil {
		return nil, err1
	}
	return model, nil
}

func GetPhotocopiesByCreator(user string) ([]*Photocopy, error) {
	msg := bson.M{"creator": user, "deleteAt": new(time.Time)}
	cursor, err1 := findMany(TablePhotocopy, msg, 0)
	if err1 != nil {
		return nil, err1
	}
	defer cursor.Close(context.Background())
	var items = make([]*Photocopy, 0, 50)
	for cursor.Next(context.Background()) {
		var node = new(Photocopy)
		if err := cursor.Decode(node); err != nil {
			return nil, err
		} else {
			items = append(items, node)
		}
	}
	return items, nil
}

// 根据母版影集UID获取克隆的所有影集
func GetPhotocopiesByMaster(master string) ([]*Photocopy, error) {
	msg := bson.M{"master": master, "deleteAt": new(time.Time)}
	cursor, err1 := findMany(TablePhotocopy, msg, 0)
	if err1 != nil {
		return nil, err1
	}
	defer cursor.Close(context.Background())
	var items = make([]*Photocopy, 0, 50)
	for cursor.Next(context.Background()) {
		var node = new(Photocopy)
		if err := cursor.Decode(node); err != nil {
			return nil, err
		} else {
			items = append(items, node)
		}
	}
	return items, nil
}

func GetPhotocopiesByTemplate(uid string) ([]*Photocopy, error) {
	msg := bson.M{"template": uid, "deleteAt": new(time.Time)}
	cursor, err1 := findMany(TablePhotocopy, msg, 0)
	if err1 != nil {
		return nil, err1
	}
	defer cursor.Close(context.Background())
	var items = make([]*Photocopy, 0, 50)
	for cursor.Next(context.Background()) {
		var node = new(Photocopy)
		if err := cursor.Decode(node); err != nil {
			return nil, err
		} else {
			items = append(items, node)
		}
	}
	return items, nil
}

func GetAllPhotocopies() ([]*Photocopy, error) {
	cursor, err1 := findAll(TablePhotocopy, 0)
	if err1 != nil {
		return nil, err1
	}
	defer cursor.Close(context.Background())
	var items = make([]*Photocopy, 0, 100)
	for cursor.Next(context.Background()) {
		var node = new(Photocopy)
		if err := cursor.Decode(node); err != nil {
			return nil, err
		} else {
			items = append(items, node)
		}
	}
	return items, nil
}

func UpdatePhotocopyBase(uid, name, remark, operator string) error {
	msg := bson.M{"name": name, "remark": remark, "operator": operator, "updatedAt": time.Now()}
	_, err := updateOne(TablePhotocopy, uid, msg)
	return err
}

func UpdatePhotocopyCover(uid, cover, operator string) error {
	msg := bson.M{"cover": cover, "operator": operator, "updatedAt": time.Now()}
	_, err := updateOne(TablePhotocopy, uid, msg)
	return err
}

func RemovePhotocopy(uid, operator string) error {
	_, err := removeOne(TablePhotocopy, uid, operator)
	return err
}

func UpdatePhotocopyCount(uid string, num uint32) error {
	msg := bson.M{"count": num, "updatedAt": time.Now()}
	_, err := updateOne(TablePhotocopy, uid, msg)
	return err
}

func AppendPhotocopyPage(uid string, slot proxy.PhotocopySlot) error {
	msg := bson.M{"slots": slot}
	_, err := appendElement(TablePhotocopy, uid, msg)
	return err
}

func SubtractPhotocopyPage(uid string, slot uint8) error {
	msg := bson.M{"slots": bson.M{"index": slot}}
	_, err := removeElement(TablePhotocopy, uid, msg)
	return err
}
