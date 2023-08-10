package nosql

import (
	"context"
	"errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"omo.msa.album/proxy"
	"time"
)

type Folder struct {
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
	Parent string `json:"parent" bson:"parent"`
	Cover  string `json:"cover" bson:"cover"`
	Access uint8  `json:"access" bson:"access"`

	Tags     []string               `json:"tags" bson:"tags"`
	Contents []*proxy.FolderContent `json:"contents" bson:"contents"`
}

func CreateFolder(info *Folder) error {
	_, err := insertOne(TableFolder, info)
	if err != nil {
		return err
	}
	return nil
}

func GetFolderNextID() uint64 {
	num, _ := getSequenceNext(TableFolder)
	return num
}

func GetFolder(uid string) (*Folder, error) {
	result, err := findOne(TableFolder, uid)
	if err != nil {
		return nil, err
	}
	model := new(Folder)
	err1 := result.Decode(model)
	if err1 != nil {
		return nil, err1
	}
	return model, nil
}

func GetFoldersByOwner(user string) ([]*Folder, error) {
	msg := bson.M{"owner": user, TimeDeleted: 0}
	cursor, err1 := findMany(TableFolder, msg, 0)
	if err1 != nil {
		return nil, err1
	}
	defer cursor.Close(context.Background())
	var items = make([]*Folder, 0, 50)
	for cursor.Next(context.Background()) {
		var node = new(Folder)
		if err := cursor.Decode(node); err != nil {
			return nil, err
		} else {
			items = append(items, node)
		}
	}
	return items, nil
}

func GetFoldersByParent(uid string) ([]*Folder, error) {
	msg := bson.M{"parent": uid, TimeDeleted: 0}
	cursor, err1 := findMany(TableFolder, msg, 0)
	if err1 != nil {
		return nil, err1
	}
	defer cursor.Close(context.Background())
	var items = make([]*Folder, 0, 50)
	for cursor.Next(context.Background()) {
		var node = new(Folder)
		if err := cursor.Decode(node); err != nil {
			return nil, err
		} else {
			items = append(items, node)
		}
	}
	return items, nil
}

func GetAllFolders() ([]*Folder, error) {
	cursor, err1 := findAllEnable(TableFolder, 0)
	if err1 != nil {
		return nil, err1
	}
	defer cursor.Close(context.Background())
	var items = make([]*Folder, 0, 100)
	for cursor.Next(context.Background()) {
		var node = new(Folder)
		if err := cursor.Decode(node); err != nil {
			return nil, err
		} else {
			items = append(items, node)
		}
	}
	return items, nil
}

func UpdateFolderBase(uid, name, remark, operator, parent string) error {
	msg := bson.M{"name": name, "remark": remark, "parent": parent,
		"operator": operator, TimeUpdated: time.Now()}
	_, err := updateOne(TableFolder, uid, msg)
	return err
}

func UpdateFolderCover(uid, cover, operator string) error {
	msg := bson.M{"cover": cover, "operator": operator, TimeUpdated: time.Now().Unix()}
	_, err := updateOne(TableFolder, uid, msg)
	return err
}

func RemoveFolder(uid, operator string) error {
	_, err := removeOne(TableFolder, uid, operator)
	return err
}

func UpdateFolderContents(uid, operator string, list []*proxy.FolderContent) error {
	if len(uid) < 2 {
		return errors.New("the uid is empty")
	}
	msg := bson.M{"contents": list, "operator": operator, TimeUpdated: time.Now().Unix()}
	_, err := updateOne(TableFolder, uid, msg)
	return err
}

func AppendFolderContent(uid string, asset string) error {
	if len(uid) < 2 {
		return errors.New("the uid is empty")
	}
	msg := bson.M{"assets": asset}
	_, err := appendElement(TableFolder, uid, msg)
	return err
}

func SubtractFolderContent(uid, asset string) error {
	if len(uid) < 2 {
		return errors.New("the uid is empty")
	}
	msg := bson.M{"contents.uid": asset}
	_, err := removeElement(TableFolder, uid, msg)
	return err
}
