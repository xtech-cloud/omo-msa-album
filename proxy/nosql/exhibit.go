package nosql

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type Exhibit struct {
	UID         primitive.ObjectID `bson:"_id"`
	ID          uint64             `json:"id" bson:"id"`
	CreatedTime time.Time          `json:"createdAt" bson:"createdAt"`
	UpdatedTime time.Time          `json:"updatedAt" bson:"updatedAt"`
	DeleteTime  time.Time          `json:"deleteAt" bson:"deleteAt"`
	Creator     string             `json:"creator" bson:"creator"`
	Operator    string             `json:"operator" bson:"operator"`

	Name   string   `json:"name" bson:"name"`
	Remark string   `json:"remark" bson:"remark"`
	Cover  string   `json:"cover" bson:"cover"`
	Status uint8    `json:"status" bson:"status"`
	Type   uint8    `json:"type" bson:"type"`
	Owner  string   `json:"owner" bson:"owner"`
	Assets []string `json:"assets" bson:"assets"`
	Tags   []string `json:"tags" bson:"tags"`
}

func CreateExhibit(info *Exhibit) error {
	_, err := insertOne(TableExhibit, info)
	if err != nil {
		return err
	}
	return nil
}

func GetExhibitNextID() uint64 {
	num, _ := getSequenceNext(TableExhibit)
	return num
}

func GetExhibit(uid string) (*Exhibit, error) {
	result, err := findOne(TableExhibit, uid)
	if err != nil {
		return nil, err
	}
	model := new(Exhibit)
	err1 := result.Decode(model)
	if err1 != nil {
		return nil, err1
	}
	return model, nil
}

func GetExhibitByName(name string) (*Exhibit, error) {
	filter := bson.M{"name": name, "deleteAt": new(time.Time)}
	result, err := findOneBy(TableExhibit, filter)
	if err != nil {
		return nil, err
	}
	model := new(Exhibit)
	err1 := result.Decode(model)
	if err1 != nil {
		return nil, err1
	}
	return model, nil
}

func GetAllExhibitsByOwner(owner string) ([]*Exhibit, error) {
	msg := bson.M{"owner": owner, "deleteAt": new(time.Time)}
	cursor, err1 := findMany(TableExhibit, msg, 0)
	if err1 != nil {
		return nil, err1
	}
	defer cursor.Close(context.Background())
	var items = make([]*Exhibit, 0, 100)
	for cursor.Next(context.Background()) {
		var node = new(Exhibit)
		if err := cursor.Decode(node); err != nil {
			return nil, err
		} else {
			items = append(items, node)
		}
	}
	return items, nil
}

func GetAllExhibits() ([]*Exhibit, error) {
	cursor, err1 := findAll(TableExhibit, 0)
	if err1 != nil {
		return nil, err1
	}
	defer cursor.Close(context.Background())
	var items = make([]*Exhibit, 0, 100)
	for cursor.Next(context.Background()) {
		var node = new(Exhibit)
		if err := cursor.Decode(node); err != nil {
			return nil, err
		} else {
			items = append(items, node)
		}
	}
	return items, nil
}

func UpdateExhibitBase(uid, name, remark, cover, operator string) error {
	msg := bson.M{"name": name, "remark": remark, "cover": cover, "operator": operator, "updatedAt": time.Now()}
	_, err := updateOne(TableExhibit, uid, msg)
	return err
}

func UpdateExhibitStatus(uid, operator string, st uint8) error {
	msg := bson.M{"status": st, "operator": operator, "updatedAt": time.Now()}
	_, err := updateOne(TableExhibit, uid, msg)
	return err
}

func UpdateExhibitAssets(uid, operator string, list []string) error {
	msg := bson.M{"assets": list, "operator": operator, "updatedAt": time.Now()}
	_, err := updateOne(TableExhibit, uid, msg)
	return err
}

func RemoveExhibit(uid, operator string) error {
	_, err := removeOne(TableExhibit, uid, operator)
	return err
}
