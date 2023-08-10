package nosql

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/labstack/gommon/log"
	"github.com/tidwall/gjson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readconcern"
	"io/ioutil"
	"mime/multipart"
	"os"
	"time"
)

var noSql *mongo.Database
var dbClient *mongo.Client

const (
	RoleAll    uint8 = 0
	RoleSelf   uint8 = 1
	RoleMaster uint8 = 2
)

func initMongoDB(ip string, port string, db string) error {
	//mongodb://myuser:mypass@localhost:40001
	addr := "mongodb://" + ip + ":" + port
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)

	opt := options.Client().ApplyURI(addr)
	opt.SetLocalThreshold(3 * time.Second)     //只使用与mongo操作耗时小于3秒的
	opt.SetMaxConnIdleTime(5 * time.Second)    //指定连接可以保持空闲的最大毫秒数
	opt.SetMaxPoolSize(200)                    //使用最大的连接数
	opt.SetReadConcern(readconcern.Majority()) //指定查询应返回实例的最新数据确认为，已写入副本集中的大多数成员
	var err error
	dbClient, err = mongo.Connect(ctx, opt)
	if err != nil {
		return err
	}
	noSql = dbClient.Database(db)

	tables, _ := noSql.ListCollectionNames(ctx, nil)
	for i := 0; i < len(tables); i++ {
		log.Info("no sql table name = " + tables[i])
	}
	return nil
}

func initMysql() error {
	/*uri := core.DBConf.User + ":" + core.DBConf.Password + "@tcp(" + core.DBConf.URL+":"+core.DBConf.Port + ")/" + core.DBConf.Name
	db, err := gorm.Open(core.DBConf.Type, uri)
	if err != nil {
		panic("failed to connect database!!!" + uri)
		return err
	}
	dbSql = db
	db.DB().SetMaxIdleConns(10)
	db.DB().SetMaxOpenConns(100)

	dbSql.LogMode(true)

	warn("connect database success!!!")
	initTeacherTable()*/
	return nil
}

func InitDB(ip string, port string, db string, kind string) error {
	if kind == "mongodb" {
		return initMongoDB(ip, port, db)
	} else {
		return initMysql()
	}
}

func CheckTimes() {
	dbs := make([]*Album, 0, 100)
	dbs = GetAll(TableAlbum, dbs)
	for _, db := range dbs {
		UpdateItemTime(TableAlbum, db.UID.Hex(), db.CreatedTime, db.UpdatedTime, db.DeleteTime)
	}
	dbs1 := make([]*Collective, 0, 100)
	dbs1 = GetAll(TableCollective, dbs1)
	for _, db := range dbs1 {
		UpdateItemTime(TableCollective, db.UID.Hex(), db.CreatedTime, db.UpdatedTime, db.DeleteTime)
	}
	dbs3 := make([]*Exhibit, 0, 100)
	dbs3 = GetAll(TableExhibit, dbs3)
	for _, db := range dbs3 {
		UpdateItemTime(TableExhibit, db.UID.Hex(), db.CreatedTime, db.UpdatedTime, db.DeleteTime)
	}
	dbs2 := make([]*PhotoFrame, 0, 100)
	dbs2 = GetAll(TablePhotoFrame, dbs2)
	for _, db := range dbs2 {
		UpdateItemTime(TablePhotoFrame, db.UID.Hex(), db.CreatedTime, db.UpdatedTime, db.DeleteTime)
	}
	dbs4 := make([]*Panorama, 0, 100)
	dbs4 = GetAll(TablePanorama, dbs4)
	for _, db := range dbs4 {
		UpdateItemTime(TablePanorama, db.UID.Hex(), db.CreatedTime, db.UpdatedTime, db.DeleteTime)
	}
	dbs5 := make([]*Photocopy, 0, 100)
	dbs5 = GetAll(TablePhotocopy, dbs5)
	for _, db := range dbs5 {
		UpdateItemTime(TablePhotocopy, db.UID.Hex(), db.CreatedTime, db.UpdatedTime, db.DeleteTime)
	}
}

func tableExist(collection string) bool {
	c := noSql.Collection(collection)
	if c == nil {
		return false
	} else {
		return true
	}
}

func checkConnected() bool {
	err := dbClient.Ping(context.TODO(), nil)
	if err != nil {
		return false
	}
	return true
}

func analyticDataStructure(table string, data []gjson.Result) error {
	return nil
}

func writeFile(path string, table string, list interface{}) error {
	f, err := os.OpenFile(path+table+".json", os.O_WRONLY|os.O_CREATE, 0666)
	defer f.Close()
	if err != nil {
		os.Remove(path + table + ".json")
		return errors.New("open the database failed")
	}
	bytes, _ := json.Marshal(list)
	_, err2 := f.Write(bytes)
	if err2 != nil {
		os.Remove(path + table + ".json")
		return errors.New("write the database failed")
	}
	return nil
}

func readFile(path string, table string) error {
	f, err := os.OpenFile(path+table+".json", os.O_RDWR, 0666)
	defer f.Close()
	if err != nil {
		return errors.New("open the database failed")
	}
	body, err := ioutil.ReadAll(f)
	if err != nil {
		return errors.New("read the file failed")
	}

	dataJson := string(body)
	result := gjson.Parse(dataJson)
	data := result.Array()

	return analyticDataStructure(table, data)
}

func ImportDatabase(table string, file multipart.File) error {
	body, err := ioutil.ReadAll(file)
	if err != nil {
		return errors.New("read the file failed")
	}

	dataJson := string(body)
	result := gjson.Parse(dataJson)
	data := result.Array()

	return analyticDataStructure(table, data)
}
