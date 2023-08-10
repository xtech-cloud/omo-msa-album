package cache

import (
	"errors"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"omo.msa.album/config"
	"omo.msa.album/proxy/nosql"
	"omo.msa.album/tool"
	"time"
)

const (
	AlbumUnknown AlbumStatus = 0
	AlbumPrivate AlbumStatus = 1
	AlbumPublic  AlbumStatus = 2
	AlbumAbandon AlbumStatus = 3
)

type AlbumStatus uint8

type AlbumType uint8

// 个人相册
type AlbumInfo struct {
	Status AlbumStatus
	Kind   AlbumType
	baseInfo
	Remark   string
	Cover    string
	Style    uint16
	Location string

	//访问密码
	Passwords string
	StarCount uint32
	// 照片最大数量
	MaxCount uint16
	// 相册体积
	Size    uint64
	Assets  []string
	Targets []string
	Tags    []string
}

func (mine *cacheContext) CreateAlbum(name, remark, user, loc string, kind uint8, style uint16, targets []string) (*AlbumInfo, error) {
	db := new(nosql.Album)
	db.UID = primitive.NewObjectID()
	db.ID = nosql.GetAlbumNextID()
	db.Created = time.Now().Unix()
	db.Creator = user
	db.Name = name
	db.Remark = remark
	db.Kind = kind
	db.Style = style
	db.Location = loc
	db.Passwords = ""
	db.Cover = ""
	db.Star = 0
	db.Targets = targets
	db.MaxCount = config.Schema.Album.Person.MaxCount
	db.Size = 0
	db.Assets = make([]string, 0, 1)
	if db.Targets == nil {
		db.Targets = make([]string, 0, 0)
	}
	db.Status = uint8(AlbumPrivate)
	db.Tags = make([]string, 0, 1)
	err := nosql.CreateAlbum(db)
	if err != nil {
		return nil, err
	}
	info := new(AlbumInfo)
	info.initInfo(db)
	return info, nil
}

func (mine *cacheContext) GetAlbum(uid string) (*AlbumInfo, error) {
	if len(uid) < 2 {
		return nil, errors.New("the album uid is empty")
	}
	db, err := nosql.GetAlbum(uid)
	if err != nil {
		return nil, err
	}
	info := new(AlbumInfo)
	info.initInfo(db)
	return info, nil
}

func (mine *cacheContext) GetAlbumsByUser(uid string) []*AlbumInfo {
	list := make([]*AlbumInfo, 0, 20)
	if len(uid) < 2 {
		return nil
	}
	array, err := nosql.GetAlbumsByCreator(uid)
	if err != nil {
		return list
	}
	for _, item := range array {
		info := new(AlbumInfo)
		info.initInfo(item)
		list = append(list, info)
	}
	return list
}

func (mine *cacheContext) GetAlbumsByTargets(targets []string) []*AlbumInfo {
	list := make([]*AlbumInfo, 0, 10)
	for _, item := range targets {
		array, err := nosql.GetAlbumsByTarget(item)
		if err != nil {
			return list
		}
		for _, item := range array {
			info := new(AlbumInfo)
			info.initInfo(item)
			list = append(list, info)
		}
	}
	return list
}

func (mine *cacheContext) GetAlbums(user string, status AlbumStatus) []*AlbumInfo {
	list := make([]*AlbumInfo, 0, 20)
	if status == AlbumPublic || status == AlbumPrivate {
		array, err := nosql.GetAlbumsByStatus(user, uint8(status))
		if err != nil {
			return list
		}
		for _, item := range array {
			info := new(AlbumInfo)
			info.initInfo(item)
			list = append(list, info)
		}
	} else {
		array, err := nosql.GetAlbumsByCreator(user)
		if err != nil {
			return list
		}
		for _, item := range array {
			info := new(AlbumInfo)
			info.initInfo(item)
			list = append(list, info)
		}
	}

	return list
}

func (mine *cacheContext) GetAlbumList(array []string) []*AlbumInfo {
	list := make([]*AlbumInfo, 0, len(array))
	for i := 0; i < len(array); i += 1 {
		db, err := nosql.GetAlbum(array[i])
		if err == nil {
			info := new(AlbumInfo)
			info.initInfo(db)
			list = append(list, info)
		}
	}
	return list
}

func (mine *AlbumInfo) initInfo(db *nosql.Album) {
	mine.Name = db.Name
	mine.UID = db.UID.Hex()
	mine.ID = db.ID
	mine.Remark = db.Remark
	mine.Kind = AlbumType(db.Kind)
	mine.Status = AlbumStatus(db.Status)
	mine.Created = db.Created
	mine.Updated = db.Updated
	mine.Creator = db.Creator
	mine.Operator = db.Operator
	mine.Cover = db.Cover
	mine.Style = db.Style
	mine.Size = db.Size
	mine.MaxCount = db.MaxCount
	mine.Location = db.Location
	mine.Passwords = db.Passwords
	mine.Tags = db.Tags
	mine.StarCount = db.Star
	mine.Targets = db.Targets
	if mine.Targets == nil {
		mine.Targets = make([]string, 0, 1)
	}
	mine.Assets = db.Assets
	if mine.Assets == nil {
		mine.Assets = make([]string, 0, 1)
	}
}

func (mine *AlbumInfo) Update(name, remark, operator, psw, loc string, style uint16, targets []string) error {
	err := nosql.UpdateAlbumBase(mine.UID, name, remark, operator, psw, loc, style)
	if err == nil {
		mine.Name = name
		mine.Remark = remark
		mine.Operator = operator
		mine.Passwords = psw
		mine.Location = loc
		mine.Style = style
		_ = mine.UpdateTargets(operator, targets)
	}
	return err
}

func (mine *AlbumInfo) UpdateBase(name, remark, operator, psw string) error {
	err := nosql.UpdateAlbumBase2(mine.UID, name, remark, operator, psw)
	if err == nil {
		mine.Name = name
		mine.Remark = remark
		mine.Operator = operator
		mine.Passwords = psw
	}
	return err
}

func (mine *AlbumInfo) Remove(operator string) error {
	return nosql.RemoveAlbum(mine.UID, operator)
}

func (mine *AlbumInfo) UpdateCover(cover, operator string) error {
	err := nosql.UpdateAlbumCover(mine.UID, cover, operator)
	if err == nil {
		mine.Cover = cover
		mine.Operator = operator
	}
	return err
}

func (mine *AlbumInfo) UpdateStyle(operator string, style uint16) error {
	err := nosql.UpdateAlbumStyle(mine.UID, operator, style)
	if err == nil {
		mine.Style = style
		mine.Operator = operator
	}
	return err
}

func (mine *AlbumInfo) UpdateStatus(operator string, status AlbumStatus) error {
	err := nosql.UpdateAlbumStatus(mine.UID, operator, uint8(status))
	if err == nil {
		mine.Status = status
		mine.Operator = operator
	}
	return err
}

func (mine *AlbumInfo) UpdateTargets(operator string, targets []string) error {
	list := targets
	if list == nil {
		list = make([]string, 0, 1)
	}
	err := nosql.UpdateAlbumTargets(mine.UID, operator, list)
	if err == nil {
		mine.Targets = list
		mine.Operator = operator
	}
	return err
}

func (mine *AlbumInfo) UpdateSize(size uint64) error {
	err := nosql.UpdateAlbumSize(mine.UID, "", size)
	if err == nil {
		mine.Size = size
	}
	return err
}

/*
获取当前相册体积大小(kb)
*/
func (mine *AlbumInfo) GetSize() uint64 {
	return mine.Size / 1024
}

func (mine *AlbumInfo) appendSize(size uint64) error {
	sum := mine.Size + size
	err := nosql.UpdateAlbumSize(mine.UID, "", sum)
	if err == nil {
		mine.Size = sum
	}
	return err
}

func (mine *AlbumInfo) subtractSize(size uint64) error {
	d := mine.Size - size
	if d < 0 {
		d = 0
	}
	err := nosql.UpdateAlbumSize(mine.UID, "", d)
	if err == nil {
		mine.Size = d
	}
	return err
}

func (mine *AlbumInfo) AppendAssets(assets []string, operator string) error {
	if assets == nil || len(assets) < 1 {
		return errors.New("the assets is nil when append")
	}
	all := make([]string, 0, len(mine.Assets)+len(assets))
	all = append(all, mine.Assets...)
	for _, asset := range assets {
		if !mine.hadAsset(asset) {
			all = append(all, asset)
		}
	}
	return mine.UpdateAssets(all, operator)
}

func (mine *AlbumInfo) UpdateAssets(assets []string, operator string) error {
	if assets == nil {
		assets = make([]string, 0, 1)
	}
	err := nosql.UpdateAlbumAssets(mine.UID, operator, assets)
	if err == nil {
		mine.Assets = assets
		mine.Operator = operator
	}
	return err
}

func (mine *AlbumInfo) RemoveAssets(assets []string, operator string) error {
	if assets == nil || len(assets) < 1 {
		return nil
	}
	all := make([]string, 0, len(mine.Assets))
	all = append(all, mine.Assets...)
	for _, asset := range assets {
		array, ok := tool.RemoveItem(all, asset)
		if ok {
			all = array
		}
	}
	return mine.UpdateAssets(all, operator)
}

func (mine *AlbumInfo) hadAsset(asset string) bool {
	for _, node := range mine.Assets {
		if node == asset {
			return true
		}
	}
	return false
}

/**
收藏第三方相册照片
*/
func (mine *AlbumInfo) SelectAssets(operator string, assets []string) error {
	if assets == nil || len(assets) < 1 {
		return errors.New("the assets length = 0")
	}
	all := make([]string, 0, len(mine.Assets)+len(assets))
	all = append(all, mine.Assets...)
	for _, asset := range assets {
		if !mine.hadAsset(asset) {
			all = append(all, asset)
		}
	}
	return mine.UpdateAssets(all, operator)
}

/**
从相册中删除一个照片
*/
func (mine *AlbumInfo) SubtractAsset(asset, operator string) error {
	if !mine.hadAsset(asset) {
		return nil
	}
	err := nosql.SubtractAlbumAsset(mine.UID, asset)
	if err == nil {
		for i := 0; i < len(mine.Assets); i += 1 {
			if mine.Assets[i] == asset {
				if i == len(mine.Assets)-1 {
					mine.Assets = append(mine.Assets[:i])
				} else {
					mine.Assets = append(mine.Assets[:i], mine.Assets[i+1:]...)
				}
				break
			}
		}
		mine.Operator = operator
	}
	return err
}
