package nosql

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type Certificate struct {
	UID      primitive.ObjectID `bson:"_id"`
	ID       uint64             `json:"id" bson:"id"`
	Created  int64              `json:"created" bson:"created"`
	Updated  int64              `json:"updated" bson:"updated"`
	Deleted  int64              `json:"deleted" bson:"deleted"`
	Creator  string             `json:"creator" bson:"creator"`
	Operator string             `json:"operator" bson:"operator"`

	Name   string   `json:"name" bson:"name"`
	Remark string   `json:"remark" bson:"remark"`
	SN     string   `json:"sn" bson:"sn"`
	Image  string   `json:"image" bson:"image"`
	Type   uint8    `json:"type" bson:"type"`
	Status uint8    `json:"status" bson:"status"`
	Style  string   `json:"style" bson:"style"`
	Target string   `json:"target" bson:"target"`
	Scene  string   `json:"scene" bson:"scene"`
	Tags   []string `json:"tags" bson:"tags"`
}

func CreateCertificate(info *Certificate) error {
	_, err := insertOne(TableCertificate, info)
	if err != nil {
		return err
	}
	return nil
}

func GetCertificateNextID() uint64 {
	num, _ := getSequenceNext(TableCertificate)
	return num
}

func GetCertificate(uid string) (*Certificate, error) {
	result, err := findOne(TableCertificate, uid)
	if err != nil {
		return nil, err
	}
	model := new(Certificate)
	err1 := result.Decode(model)
	if err1 != nil {
		return nil, err1
	}
	return model, nil
}

func GetCertificatesByScene(scene string) ([]*Certificate, error) {
	msg := bson.M{"scene": scene, TimeDeleted: 0}
	cursor, err1 := findMany(TableCertificate, msg, 0)
	if err1 != nil {
		return nil, err1
	}
	defer cursor.Close(context.Background())
	var items = make([]*Certificate, 0, 50)
	for cursor.Next(context.Background()) {
		var node = new(Certificate)
		if err := cursor.Decode(node); err != nil {
			return nil, err
		} else {
			items = append(items, node)
		}
	}
	return items, nil
}

func GetCertificatesByTarget(uid string) ([]*Certificate, error) {
	msg := bson.M{"target": uid, TimeDeleted: 0}
	cursor, err1 := findMany(TableCertificate, msg, 0)
	if err1 != nil {
		return nil, err1
	}
	defer cursor.Close(context.Background())
	var items = make([]*Certificate, 0, 50)
	for cursor.Next(context.Background()) {
		var node = new(Certificate)
		if err := cursor.Decode(node); err != nil {
			return nil, err
		} else {
			items = append(items, node)
		}
	}
	return items, nil
}

func GetCertificatesByStyle(style string) ([]*Certificate, error) {
	msg := bson.M{"style": style, TimeDeleted: 0}
	cursor, err1 := findMany(TableCertificate, msg, 0)
	if err1 != nil {
		return nil, err1
	}
	defer cursor.Close(context.Background())
	var items = make([]*Certificate, 0, 50)
	for cursor.Next(context.Background()) {
		var node = new(Certificate)
		if err := cursor.Decode(node); err != nil {
			return nil, err
		} else {
			items = append(items, node)
		}
	}
	return items, nil
}

func GetCertificatesCountByStyle(uid string) (uint32, error) {
	msg := bson.M{"style": uid, TimeDeleted: 0}
	num, err1 := getCountByFilter(TableCertificate, msg)
	if err1 != nil {
		return 0, err1
	}
	return uint32(num), nil
}

func GetAllCertificates() ([]*Certificate, error) {
	cursor, err1 := findAllEnable(TableCertificate, 0)
	if err1 != nil {
		return nil, err1
	}
	defer cursor.Close(context.Background())
	var items = make([]*Certificate, 0, 100)
	for cursor.Next(context.Background()) {
		var node = new(Certificate)
		if err := cursor.Decode(node); err != nil {
			return nil, err
		} else {
			items = append(items, node)
		}
	}
	return items, nil
}

func UpdateCertificateBase(uid, name, remark, operator string, tags []string) error {
	msg := bson.M{"name": name, "remark": remark, "tags": tags,
		"operator": operator, TimeUpdated: time.Now()}
	_, err := updateOne(TableCertificate, uid, msg)
	return err
}

func UpdateCertificateCover(uid, cover, operator string) error {
	msg := bson.M{"cover": cover, "operator": operator, TimeUpdated: time.Now().Unix()}
	_, err := updateOne(TableCertificate, uid, msg)
	return err
}

func RemoveCertificate(uid, operator string) error {
	_, err := removeOne(TableCertificate, uid, operator)
	return err
}
