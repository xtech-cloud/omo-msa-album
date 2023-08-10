package nosql

import (
	"context"
	"errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"omo.msa.album/proxy"
	"time"
)

type Page struct {
	UID      primitive.ObjectID `bson:"_id"`
	ID       uint64             `json:"id" bson:"id"`
	Created  int64              `json:"created" bson:"created"`
	Updated  int64              `json:"updated" bson:"updated"`
	Deleted  int64              `json:"deleted" bson:"deleted"`
	Creator  string             `json:"creator" bson:"creator"`
	Operator string             `json:"operator" bson:"operator"`

	Name        string `json:"name" bson:"name"`
	Remark      string `json:"remark" bson:"remark"`
	Owner       string `json:"owner" bson:"owner"`
	Type        uint8  `json:"type" bson:"type"`
	Lifecycle   uint32 `json:"lifecycle" bson:"lifecycle"`
	Composition string `json:"composition" bson:"composition"`

	Tags     []string              `json:"tags" bson:"tags"`
	Contents []*proxy.PageContents `json:"contents" bson:"contents"`
}

func CreatePage(info *Page) error {
	_, err := insertOne(TablePage, info)
	if err != nil {
		return err
	}
	return nil
}

func GetPageNextID() uint64 {
	num, _ := getSequenceNext(TablePage)
	return num
}

func GetPage(uid string) (*Page, error) {
	result, err := findOne(TablePage, uid)
	if err != nil {
		return nil, err
	}
	model := new(Page)
	err1 := result.Decode(model)
	if err1 != nil {
		return nil, err1
	}
	return model, nil
}

func GetPagesByOwner(user string) ([]*Page, error) {
	msg := bson.M{"owner": user, TimeDeleted: 0}
	cursor, err1 := findMany(TablePage, msg, 0)
	if err1 != nil {
		return nil, err1
	}
	defer cursor.Close(context.Background())
	var items = make([]*Page, 0, 50)
	for cursor.Next(context.Background()) {
		var node = new(Page)
		if err := cursor.Decode(node); err != nil {
			return nil, err
		} else {
			items = append(items, node)
		}
	}
	return items, nil
}

func GetAllPages() ([]*Page, error) {
	cursor, err1 := findAllEnable(TablePage, 0)
	if err1 != nil {
		return nil, err1
	}
	defer cursor.Close(context.Background())
	var items = make([]*Page, 0, 100)
	for cursor.Next(context.Background()) {
		var node = new(Page)
		if err := cursor.Decode(node); err != nil {
			return nil, err
		} else {
			items = append(items, node)
		}
	}
	return items, nil
}

func UpdatePageBase(uid, name, remark, operator string, life uint32) error {
	msg := bson.M{"name": name, "remark": remark, "lifecycle": life,
		"operator": operator, TimeUpdated: time.Now().Unix()}
	_, err := updateOne(TablePage, uid, msg)
	return err
}

func RemovePage(uid, operator string) error {
	_, err := removeOne(TablePage, uid, operator)
	return err
}

func UpdatePageContents(uid, operator string, list []*proxy.PageContents) error {
	if len(uid) < 2 {
		return errors.New("the uid is empty")
	}
	msg := bson.M{"contents": list, "operator": operator, TimeUpdated: time.Now().Unix()}
	_, err := updateOne(TablePage, uid, msg)
	return err
}

func AppendPageContent(uid string, slot *proxy.PageContents) error {
	if len(uid) < 2 {
		return errors.New("the uid is empty")
	}
	msg := bson.M{"contents": slot}
	_, err := appendElement(TablePage, uid, msg)
	return err
}

func SubtractPageContent(uid string, slot uint32) error {
	if len(uid) < 2 {
		return errors.New("the uid is empty")
	}
	msg := bson.M{"contents.slot": slot}
	_, err := removeElement(TablePage, uid, msg)
	return err
}
