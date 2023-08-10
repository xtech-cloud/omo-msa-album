package nosql

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type Panorama struct {
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

	Name    string `json:"name" bson:"name"`
	Remark  string `json:"remark" bson:"remark"`
	Content string `json:"content" bson:"content"`
	Owner   string `json:"owner" bson:"owner"`
}

func CreatePanorama(info *Panorama) error {
	_, err := insertOne(TablePanorama, info)
	if err != nil {
		return err
	}
	return nil
}

func GetPanoramaNextID() uint64 {
	num, _ := getSequenceNext(TablePanorama)
	return num
}

func GetPanorama(uid string) (*Panorama, error) {
	result, err := findOne(TablePanorama, uid)
	if err != nil {
		return nil, err
	}
	model := new(Panorama)
	err1 := result.Decode(model)
	if err1 != nil {
		return nil, err1
	}
	return model, nil
}

func GetAllPanoramasByOwner(owner string) ([]*Panorama, error) {
	msg := bson.M{"owner": owner, TimeDeleted: 0}
	cursor, err1 := findMany(TablePanorama, msg, 0)
	if err1 != nil {
		return nil, err1
	}
	defer cursor.Close(context.Background())
	var items = make([]*Panorama, 0, 100)
	for cursor.Next(context.Background()) {
		var node = new(Panorama)
		if err := cursor.Decode(node); err != nil {
			return nil, err
		} else {
			items = append(items, node)
		}
	}
	return items, nil
}

func GetAllPanoramas() ([]*Panorama, error) {
	cursor, err1 := findAllEnable(TablePanorama, 0)
	if err1 != nil {
		return nil, err1
	}
	defer cursor.Close(context.Background())
	var items = make([]*Panorama, 0, 100)
	for cursor.Next(context.Background()) {
		var node = new(Panorama)
		if err := cursor.Decode(node); err != nil {
			return nil, err
		} else {
			items = append(items, node)
		}
	}
	return items, nil
}

func UpdatePanoramaBase(uid, name, remark string) error {
	msg := bson.M{"name": name, "remark": remark, TimeUpdated: time.Now().Unix()}
	_, err := updateOne(TablePanorama, uid, msg)
	return err
}

func UpdatePanoramaContent(uid, content string) error {
	msg := bson.M{"content": content, TimeUpdated: time.Now().Unix()}
	_, err := updateOne(TablePanorama, uid, msg)
	return err
}

func UpdatePanoramaOwner(uid, owner string) error {
	msg := bson.M{"owner": owner}
	_, err := updateOne(TablePanorama, uid, msg)
	return err
}

func RemovePanorama(uid, operator string) error {
	_, err := removeOne(TablePanorama, uid, operator)
	return err
}
