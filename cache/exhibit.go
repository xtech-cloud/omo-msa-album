package cache

import (
	"errors"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"omo.msa.album/proxy/nosql"
	"time"
)

type ExhibitInfo struct {
	Status uint8
	Type   uint8
	baseInfo
	Remark string
	Cover  string
	Owner  string
	Tags   []string
	Assets []string
}

func (mine *cacheContext) CreateExhibit(name, remark, cover, owner, operator string) (*ExhibitInfo, error) {
	if owner == "" {
		owner = DefaultOwner
	}
	db := new(nosql.Exhibit)
	db.UID = primitive.NewObjectID()
	db.ID = nosql.GetExhibitNextID()
	db.CreatedTime = time.Now()
	db.Creator = operator
	db.Name = name
	db.Remark = remark
	db.Cover = cover
	db.Owner = owner
	db.Status = 0
	db.Type = 0
	db.Assets = make([]string, 0, 1)
	db.Tags = make([]string, 0, 1)
	err := nosql.CreateExhibit(db)
	if err != nil {
		return nil, err
	}
	info := new(ExhibitInfo)
	info.initInfo(db)
	return info, nil
}

func (mine *cacheContext) HadExhibitByName(name string) bool {
	db, err := nosql.GetExhibitByName(name)
	if err != nil {
		return false
	}
	if db == nil {
		return false
	} else {
		return true
	}
}

func (mine *cacheContext) GetExhibit(uid string) (*ExhibitInfo, error) {
	db, err := nosql.GetExhibit(uid)
	if err != nil {
		return nil, err
	}
	info := new(ExhibitInfo)
	info.initInfo(db)
	return info, nil
}

func (mine *cacheContext) RemoveExhibit(uid, operator string) error {
	if uid == "" {
		return errors.New("the Exhibit uid is empty")
	}
	err := nosql.RemoveExhibit(uid, operator)
	return err
}

func (mine *cacheContext) GetExhibits(array []string) []*ExhibitInfo {
	if array == nil {
		return make([]*ExhibitInfo, 0, 1)
	}
	list := make([]*ExhibitInfo, 0, len(array))
	for _, item := range array {
		info, _ := mine.GetExhibit(item)
		if info != nil {
			list = append(list, info)
		}
	}
	return list
}

func (mine *cacheContext) GetAllExhibits(owner string) ([]*ExhibitInfo, error) {
	if owner == "" {
		owner = DefaultOwner
	}
	array, err := nosql.GetAllExhibitsByOwner(owner)
	if err != nil {
		return nil, err
	}
	list := make([]*ExhibitInfo, 0, len(array))
	for _, item := range array {
		info := new(ExhibitInfo)
		info.initInfo(item)
		list = append(list, info)
	}
	return list, nil
}

func (mine *ExhibitInfo) initInfo(db *nosql.Exhibit) {
	mine.UID = db.UID.Hex()
	mine.Name = db.Name
	mine.CreateTime = db.CreatedTime
	mine.Remark = db.Remark
	mine.Cover = db.Cover
	mine.UpdateTime = db.UpdatedTime
	mine.Creator = db.Creator
	mine.Operator = db.Operator
	mine.Type = db.Type
	mine.Status = db.Status
	mine.Owner = db.Owner
	mine.Assets = db.Assets
	mine.Tags = db.Tags
}

func (mine *ExhibitInfo) UpdateBase(name, remark, cover, operator string) error {
	err := nosql.UpdateExhibitBase(mine.UID, name, remark, cover, operator)
	if err == nil {
		mine.Name = name
		mine.Remark = remark
		mine.Cover = cover
		mine.Operator = operator
		mine.UpdateTime = time.Now()
	}
	return err
}

func (mine *ExhibitInfo) HadAsset(asset string) bool {
	for i := 0; i < len(mine.Assets); i += 1 {
		if mine.Assets[i] == asset {
			return true
		}
	}
	return false
}

func (mine *ExhibitInfo) UpdateStatus(operator string, status uint8) error {
	err := nosql.UpdateExhibitStatus(mine.UID, operator, uint8(status))
	if err == nil {
		mine.Status = status
		mine.Operator = operator
	}
	return err
}

func (mine *ExhibitInfo) UpdateAssets(assets []string, operator string) error {
	if len(assets) < 1 {
		return errors.New("the assets is empty")
	}
	err := nosql.UpdateExhibitAssets(mine.UID, operator, assets)
	if err == nil {
		mine.Assets = assets
		mine.Operator = operator
		mine.UpdateTime = time.Now()
	}
	return err
}

func (mine *ExhibitInfo) AppendAsset(asset, operator string) error {
	if len(asset) < 1 {
		return errors.New("the asset uid is empty")
	}
	if mine.HadAsset(asset) {
		return nil
	}
	array := make([]string, 0, len(mine.Assets)+1)
	array = append(array, mine.Assets[0:]...)
	array = append(array, asset)
	return mine.UpdateAssets(array, operator)
}

func (mine *ExhibitInfo) SubtractAsset(asset, operator string) error {
	if len(asset) < 1 {
		return errors.New("the asset uid is empty")
	}
	if !mine.HadAsset(asset) {
		return nil
	}
	array := make([]string, 0, len(mine.Assets))
	for _, uid := range mine.Assets {
		if uid != asset {
			array = append(array, uid)
		}
	}
	err := nosql.UpdateExhibitAssets(mine.UID, operator, array)
	if err == nil {
		mine.Assets = array
		mine.Operator = operator
		mine.UpdateTime = time.Now()
	}
	return err
}
