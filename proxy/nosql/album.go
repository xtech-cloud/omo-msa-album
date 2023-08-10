package nosql

import (
	"context"
	"errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type Album struct {
	UID         primitive.ObjectID `bson:"_id"`
	ID          uint64             `json:"id" bson:"id"`
	CreatedTime time.Time          `json:"createdAt" bson:"createdAt"`
	UpdatedTime time.Time          `json:"TimeCreatedAt" bson:"TimeCreatedAt"`
	DeleteTime  time.Time          `json:"deletedAt" bson:"deletedAt"`
	Created     int64              `json:"created" bson:"created"`
	Updated     int64              `json:"updated" bson:"updated"`
	Deleted     int64              `json:"deleted" bson:"deleted"`
	Creator     string             `json:"creator" bson:"creator"`
	Operator    string             `json:"operator" bson:"operator"`

	Status    uint8  `json:"status" bson:"status"`
	Name      string `json:"name" bson:"name"`
	Remark    string `json:"remark" bson:"remark"`
	Kind      uint8  `json:"kind" bson:"kind"`
	Style     uint16 `json:"style" bson:"style"`
	Star      uint32 `json:"star" bson:"star"`
	Cover     string `json:"cover" bson:"cover"`
	Passwords string `json:"passwords" bson:"passwords"`
	Location  string `json:"location" bson:"location"`
	MaxCount  uint16 `json:"max" bson:"max"`
	Size      uint64 `json:"size" bson:"size"`
	// 可访问相册的组织，机构，家庭等UID
	Targets []string `json:"targets" bson:"targets"`
	Assets  []string `json:"assets" bson:"assets"`
	Tags    []string `json:"tags" bson:"tags"`
}

func CreateAlbum(info *Album) error {
	_, err := insertOne(TableAlbum, info)
	if err != nil {
		return err
	}
	return nil
}

func GetAlbumNextID() uint64 {
	num, _ := getSequenceNext(TableAlbum)
	return num
}

func GetAlbum(uid string) (*Album, error) {
	result, err := findOne(TableAlbum, uid)
	if err != nil {
		return nil, err
	}
	model := new(Album)
	err1 := result.Decode(model)
	if err1 != nil {
		return nil, err1
	}
	return model, nil
}

func GetAlbumsByCreator(user string) ([]*Album, error) {
	msg := bson.M{"creator": user, TimeDeleted: 0}
	cursor, err1 := findMany(TableAlbum, msg, 0)
	if err1 != nil {
		return nil, err1
	}
	defer cursor.Close(context.Background())
	var items = make([]*Album, 0, 50)
	for cursor.Next(context.Background()) {
		var node = new(Album)
		if err := cursor.Decode(node); err != nil {
			return nil, err
		} else {
			items = append(items, node)
		}
	}
	return items, nil
}

func GetAlbumsByStatus(user string, st uint8) ([]*Album, error) {
	msg := bson.M{"creator": user, "status": st, TimeDeleted: 0}
	cursor, err1 := findMany(TableAlbum, msg, 0)
	if err1 != nil {
		return nil, err1
	}
	defer cursor.Close(context.Background())
	var items = make([]*Album, 0, 50)
	for cursor.Next(context.Background()) {
		var node = new(Album)
		if err := cursor.Decode(node); err != nil {
			return nil, err
		} else {
			items = append(items, node)
		}
	}
	return items, nil
}

func GetAllAlbums() ([]*Album, error) {
	cursor, err1 := findAllEnable(TableAlbum, 0)
	if err1 != nil {
		return nil, err1
	}
	defer cursor.Close(context.Background())
	var items = make([]*Album, 0, 100)
	for cursor.Next(context.Background()) {
		var node = new(Album)
		if err := cursor.Decode(node); err != nil {
			return nil, err
		} else {
			items = append(items, node)
		}
	}
	return items, nil
}

func UpdateAlbumBase(uid, name, remark, operator, psw, loc string, style uint16) error {
	msg := bson.M{"name": name, "remark": remark, "passwords": psw,
		"location": loc, "style": style, "operator": operator, TimeUpdated: time.Now()}
	_, err := updateOne(TableAlbum, uid, msg)
	return err
}

func UpdateAlbumBase2(uid, name, remark, operator, psw string) error {
	msg := bson.M{"name": name, "remark": remark, "passwords": psw, "operator": operator, TimeUpdated: time.Now().Unix()}
	_, err := updateOne(TableAlbum, uid, msg)
	return err
}

func UpdateAlbumCover(uid, cover, operator string) error {
	msg := bson.M{"cover": cover, "operator": operator, TimeUpdated: time.Now().Unix()}
	_, err := updateOne(TableAlbum, uid, msg)
	return err
}

func UpdateAlbumStatus(uid, operator string, st uint8) error {
	msg := bson.M{"status": st, "operator": operator, TimeUpdated: time.Now().Unix()}
	_, err := updateOne(TableAlbum, uid, msg)
	return err
}

func UpdateAlbumSize(uid, operator string, size uint64) error {
	msg := bson.M{"size": size, "operator": operator, TimeUpdated: time.Now().Unix()}
	_, err := updateOne(TableAlbum, uid, msg)
	return err
}

func UpdateAlbumStyle(uid, operator string, st uint16) error {
	msg := bson.M{"style": st, "operator": operator, TimeUpdated: time.Now().Unix()}
	_, err := updateOne(TableAlbum, uid, msg)
	return err
}

func RemoveAlbum(uid, operator string) error {
	_, err := removeOne(TableAlbum, uid, operator)
	return err
}

func UpdateAlbumAssets(uid, operator string, assets []string) error {
	if len(uid) < 2 {
		return errors.New("the uid is empty")
	}
	msg := bson.M{"assets": assets, "operator": operator, TimeUpdated: time.Now()}
	_, err := updateOne(TableAlbum, uid, msg)
	return err
}

func UpdateAlbumTargets(uid, operator string, targets []string) error {
	if len(uid) < 2 {
		return errors.New("the uid is empty")
	}
	msg := bson.M{"targets": targets, "operator": operator, TimeUpdated: time.Now().Unix()}
	_, err := updateOne(TableAlbum, uid, msg)
	return err
}

func GetAlbumsByTarget(target string) ([]*Album, error) {
	msg := bson.M{"targets": bson.M{"$elemMatch": bson.M{"$eq": target}}, TimeDeleted: 0}
	cursor, err1 := findMany(TableAlbum, msg, 0)
	if err1 != nil {
		return nil, err1
	}
	defer cursor.Close(context.Background())
	var items = make([]*Album, 0, 10)
	for cursor.Next(context.Background()) {
		var node = new(Album)
		if err := cursor.Decode(node); err != nil {
			return nil, err
		} else {
			items = append(items, node)
		}
	}
	return items, nil
}

func AppendAlbumTarget(uid string, target string) error {
	if len(uid) < 2 {
		return errors.New("the uid is empty")
	}
	msg := bson.M{"targets": target}
	_, err := appendElement(TableAlbum, uid, msg)
	return err
}

func SubtractAlbumTarget(uid, target string) error {
	if len(uid) < 2 {
		return errors.New("the uid is empty")
	}
	msg := bson.M{"targets": target}
	_, err := removeElement(TableAlbum, uid, msg)
	return err
}

func AppendAlbumAsset(uid string, asset string) error {
	if len(uid) < 2 {
		return errors.New("the uid is empty")
	}
	msg := bson.M{"assets": asset}
	_, err := appendElement(TableAlbum, uid, msg)
	return err
}

func SubtractAlbumAsset(uid, asset string) error {
	if len(uid) < 2 {
		return errors.New("the uid is empty")
	}
	msg := bson.M{"assets": asset}
	_, err := removeElement(TableAlbum, uid, msg)
	return err
}
