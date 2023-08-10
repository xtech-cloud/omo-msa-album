package nosql

import (
	"context"
	"errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"omo.msa.album/proxy"
	"time"
)

type Composition struct {
	UID      primitive.ObjectID `bson:"_id"`
	ID       uint64             `json:"id" bson:"id"`
	Created  int64              `json:"created" bson:"created"`
	Updated  int64              `json:"updated" bson:"updated"`
	Deleted  int64              `json:"deleted" bson:"deleted"`
	Creator  string             `json:"creator" bson:"creator"`
	Operator string             `json:"operator" bson:"operator"`

	Name   string `json:"name" bson:"name"`
	Remark string `json:"remark" bson:"remark"`
	Owner  string `json:"owner" bson:"owner"`
	Aspect uint8  `json:"aspect" bson:"aspect"`
	Cover  string `json:"cover" bson:"cover"`

	Tags  []string          `json:"tags" bson:"tags"`
	Slots []*proxy.SlotInfo `json:"slots" bson:"slots"`
}

func CreateComposition(info *Composition) error {
	_, err := insertOne(TableComposition, info)
	if err != nil {
		return err
	}
	return nil
}

func GetCompositionNextID() uint64 {
	num, _ := getSequenceNext(TableComposition)
	return num
}

func GetComposition(uid string) (*Composition, error) {
	result, err := findOne(TableComposition, uid)
	if err != nil {
		return nil, err
	}
	model := new(Composition)
	err1 := result.Decode(model)
	if err1 != nil {
		return nil, err1
	}
	return model, nil
}

func GetCompositionsByOwner(user string) ([]*Composition, error) {
	msg := bson.M{"owner": user, TimeDeleted: 0}
	cursor, err1 := findMany(TableComposition, msg, 0)
	if err1 != nil {
		return nil, err1
	}
	defer cursor.Close(context.Background())
	var items = make([]*Composition, 0, 50)
	for cursor.Next(context.Background()) {
		var node = new(Composition)
		if err := cursor.Decode(node); err != nil {
			return nil, err
		} else {
			items = append(items, node)
		}
	}
	return items, nil
}

func GetAllCompositions() ([]*Composition, error) {
	cursor, err1 := findAllEnable(TableComposition, 0)
	if err1 != nil {
		return nil, err1
	}
	defer cursor.Close(context.Background())
	var items = make([]*Composition, 0, 100)
	for cursor.Next(context.Background()) {
		var node = new(Composition)
		if err := cursor.Decode(node); err != nil {
			return nil, err
		} else {
			items = append(items, node)
		}
	}
	return items, nil
}

func UpdateCompositionBase(uid, name, remark, operator string, aspect uint8) error {
	msg := bson.M{"name": name, "remark": remark, "aspect": aspect,
		"operator": operator, TimeUpdated: time.Now().Unix()}
	_, err := updateOne(TableComposition, uid, msg)
	return err
}

func RemoveComposition(uid, operator string) error {
	_, err := removeOne(TableComposition, uid, operator)
	return err
}

func UpdateCompositionCover(uid, cover, operator string) error {
	msg := bson.M{"cover": cover, "operator": operator, TimeUpdated: time.Now().Unix()}
	_, err := updateOne(TableComposition, uid, msg)
	return err
}

func UpdateCompositionSlots(uid, operator string, list []*proxy.SlotInfo) error {
	if len(uid) < 2 {
		return errors.New("the uid is empty")
	}
	msg := bson.M{"contents": list, "operator": operator, TimeUpdated: time.Now().Unix()}
	_, err := updateOne(TableComposition, uid, msg)
	return err
}
