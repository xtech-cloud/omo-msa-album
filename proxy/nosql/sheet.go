package nosql

import (
	"context"
	"errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"omo.msa.album/proxy"
	"time"
)

type Sheet struct {
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
	Cover  string `json:"cover" bson:"cover"`
	Target string `json:"target" bson:"target"`
	Aspect string `json:"aspect" bson:"aspect"`
	Size   uint32 `json:"size" bson:"size"`
	Type   uint8  `json:"type" bson:"type"`
	Status uint8  `json:"status" bson:"status"`

	Tags  []string           `json:"tags" bson:"tags"`
	Pages []*proxy.SheetPage `json:"pages" bson:"pages"`
}

func CreateSheet(info *Sheet) error {
	_, err := insertOne(TableSheet, info)
	if err != nil {
		return err
	}
	return nil
}

func GetSheetNextID() uint64 {
	num, _ := getSequenceNext(TableSheet)
	return num
}

func GetSheet(uid string) (*Sheet, error) {
	result, err := findOne(TableSheet, uid)
	if err != nil {
		return nil, err
	}
	model := new(Sheet)
	err1 := result.Decode(model)
	if err1 != nil {
		return nil, err1
	}
	return model, nil
}

func GetSheetsByOwner(user string) ([]*Sheet, error) {
	msg := bson.M{"owner": user, TimeDeleted: 0}
	cursor, err1 := findMany(TableSheet, msg, 0)
	if err1 != nil {
		return nil, err1
	}
	defer cursor.Close(context.Background())
	var items = make([]*Sheet, 0, 50)
	for cursor.Next(context.Background()) {
		var node = new(Sheet)
		if err := cursor.Decode(node); err != nil {
			return nil, err
		} else {
			items = append(items, node)
		}
	}
	return items, nil
}

func GetSheetsByTarget(target string) ([]*Sheet, error) {
	msg := bson.M{"target": target, TimeDeleted: 0}
	cursor, err1 := findMany(TableSheet, msg, 0)
	if err1 != nil {
		return nil, err1
	}
	defer cursor.Close(context.Background())
	var items = make([]*Sheet, 0, 50)
	for cursor.Next(context.Background()) {
		var node = new(Sheet)
		if err := cursor.Decode(node); err != nil {
			return nil, err
		} else {
			items = append(items, node)
		}
	}
	return items, nil
}

func GetSheetsByPage(page string) ([]*Sheet, error) {
	msg := bson.M{"pages.uid": page, TimeDeleted: 0}
	cursor, err1 := findMany(TableSheet, msg, 0)
	if err1 != nil {
		return nil, err1
	}
	defer cursor.Close(context.Background())
	var items = make([]*Sheet, 0, 50)
	for cursor.Next(context.Background()) {
		var node = new(Sheet)
		if err := cursor.Decode(node); err != nil {
			return nil, err
		} else {
			items = append(items, node)
		}
	}
	return items, nil
}

func GetAllSheets() ([]*Sheet, error) {
	cursor, err1 := findAllEnable(TableSheet, 0)
	if err1 != nil {
		return nil, err1
	}
	defer cursor.Close(context.Background())
	var items = make([]*Sheet, 0, 100)
	for cursor.Next(context.Background()) {
		var node = new(Sheet)
		if err := cursor.Decode(node); err != nil {
			return nil, err
		} else {
			items = append(items, node)
		}
	}
	return items, nil
}

func UpdateSheetBase(uid, name, remark, operator string) error {
	msg := bson.M{"name": name, "remark": remark,
		"operator": operator, TimeUpdated: time.Now().Unix()}
	_, err := updateOne(TableSheet, uid, msg)
	return err
}

func UpdateSheetCover(uid, cover, operator string) error {
	msg := bson.M{"cover": cover, "operator": operator, TimeUpdated: time.Now().Unix()}
	_, err := updateOne(TableSheet, uid, msg)
	return err
}

func UpdateSheetTarget(uid, target, operator string) error {
	msg := bson.M{"target": target, "operator": operator, TimeUpdated: time.Now().Unix()}
	_, err := updateOne(TableSheet, uid, msg)
	return err
}

func RemoveSheet(uid, operator string) error {
	_, err := removeOne(TableSheet, uid, operator)
	return err
}

func UpdateSheetPages(uid, operator string, list []*proxy.SheetPage) error {
	if len(uid) < 2 {
		return errors.New("the uid is empty")
	}
	msg := bson.M{"pages": list, "operator": operator, TimeUpdated: time.Now().Unix()}
	_, err := updateOne(TableSheet, uid, msg)
	return err
}

func AppendSheetPages(uid string, page *proxy.SheetPage) error {
	if len(uid) < 2 {
		return errors.New("the uid is empty")
	}
	msg := bson.M{"pages": page}
	_, err := appendElement(TableSheet, uid, msg)
	return err
}

func SubtractSheetPage(uid, page string) error {
	if len(uid) < 2 {
		return errors.New("the uid is empty")
	}
	msg := bson.M{"pages.uid": page}
	_, err := removeElement(TableSheet, uid, msg)
	return err
}
