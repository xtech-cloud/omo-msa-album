package nosql

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"omo.msa.album/proxy"
	"time"
)

type Collective struct {
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

	Status   uint8              `json:"status" bson:"status"`
	Name     string             `json:"name" bson:"name"`
	Remark   string             `json:"remark" bson:"remark"`
	MaxCount uint16             `json:"max" bson:"max"`
	Type     uint8              `json:"type" bson:"type"`
	Group    string             `json:"group" bson:"group"`
	Size     uint64             `json:"size" bson:"size"`
	Cover    string             `json:"cover" bson:"cover"`
	Date     proxy.DurationInfo `json:"date" bson:"date"`
	Assets   []string           `json:"assets" bson:"assets"`
	Tags     []string           `json:"tags" bson:"tags"`
}

func CreateCollective(info *Collective) error {
	_, err := insertOne(TableCollective, info)
	if err != nil {
		return err
	}
	return nil
}

func GetCollectiveNextID() uint64 {
	num, _ := getSequenceNext(TableCollective)
	return num
}

func GetCollective(uid string) (*Collective, error) {
	result, err := findOne(TableCollective, uid)
	if err != nil {
		return nil, err
	}
	model := new(Collective)
	err1 := result.Decode(model)
	if err1 != nil {
		return nil, err1
	}
	return model, nil
}

func GetCollectiveByName(owner, name string) (*Collective, error) {
	msg := bson.M{"group": owner, "name": name, TimeDeleted: 0}
	result, err := findOneBy(TableCollective, msg)
	if err != nil {
		return nil, err
	}
	model := new(Collective)
	err1 := result.Decode(model)
	if err1 != nil {
		return nil, err1
	}
	return model, nil
}

func GetCollectivesByCreator(user string) ([]*Collective, error) {
	msg := bson.M{"creator": user, TimeDeleted: 0}
	cursor, err1 := findMany(TableCollective, msg, 0)
	if err1 != nil {
		return nil, err1
	}
	defer cursor.Close(context.Background())
	var items = make([]*Collective, 0, 50)
	for cursor.Next(context.Background()) {
		var node = new(Collective)
		if err := cursor.Decode(node); err != nil {
			return nil, err
		} else {
			items = append(items, node)
		}
	}
	return items, nil
}

func GetCollectivesByGroup(group string) ([]*Collective, error) {
	msg := bson.M{"group": group, TimeDeleted: 0}
	cursor, err1 := findMany(TableCollective, msg, 0)
	if err1 != nil {
		return nil, err1
	}
	defer cursor.Close(context.Background())
	var items = make([]*Collective, 0, 50)
	for cursor.Next(context.Background()) {
		var node = new(Collective)
		if err := cursor.Decode(node); err != nil {
			return nil, err
		} else {
			items = append(items, node)
		}
	}
	return items, nil
}

func GetAllCollectives() ([]*Collective, error) {
	cursor, err1 := findAllEnable(TableCollective, 0)
	if err1 != nil {
		return nil, err1
	}
	defer cursor.Close(context.Background())
	var items = make([]*Collective, 0, 100)
	for cursor.Next(context.Background()) {
		var node = new(Collective)
		if err := cursor.Decode(node); err != nil {
			return nil, err
		} else {
			items = append(items, node)
		}
	}
	return items, nil
}

func UpdateCollectiveBase(uid, name, remark, operator string) error {
	msg := bson.M{"name": name, "remark": remark, "operator": operator, TimeUpdated: time.Now().Unix()}
	_, err := updateOne(TableCollective, uid, msg)
	return err
}

func UpdateCollectiveCover(uid, cover, operator string) error {
	msg := bson.M{"cover": cover, "operator": operator, TimeUpdated: time.Now().Unix()}
	_, err := updateOne(TableCollective, uid, msg)
	return err
}

func UpdateCollectiveStatus(uid, operator string, st uint8) error {
	msg := bson.M{"status": st, "operator": operator, TimeUpdated: time.Now().Unix()}
	_, err := updateOne(TableCollective, uid, msg)
	return err
}

func UpdateCollectiveSize(uid, operator string, size uint64) error {
	msg := bson.M{"size": size, "operator": operator, TimeUpdated: time.Now().Unix()}
	_, err := updateOne(TableCollective, uid, msg)
	return err
}

func UpdateCollectiveDate(uid, operator string, begin int64) error {
	msg := bson.M{"date.start": begin, "operator": operator, TimeUpdated: time.Now().Unix()}
	_, err := updateOne(TableCollective, uid, msg)
	return err
}

func RemoveCollective(uid, operator string) error {
	_, err := removeOne(TableCollective, uid, operator)
	return err
}

func UpdateCollectiveAssets(uid, operator string, assets []string) error {
	msg := bson.M{"assets": assets, "operator": operator, TimeUpdated: time.Now().Unix()}
	_, err := updateOne(TableCollective, uid, msg)
	return err
}

func AppendCollectiveAsset(uid string, asset string) error {
	msg := bson.M{"assets": asset}
	_, err := appendElement(TableCollective, uid, msg)
	return err
}

func SubtractCollectiveAsset(uid, asset string) error {
	msg := bson.M{"assets": asset}
	_, err := removeElement(TableCollective, uid, msg)
	return err
}
