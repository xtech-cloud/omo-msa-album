package cache

import (
	"errors"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"omo.msa.album/config"
	"omo.msa.album/proxy"
	"omo.msa.album/proxy/nosql"
	"omo.msa.album/tool"
	"time"
)

type CollAlbumInfo struct {
	Status uint8
	Type   uint8
	baseInfo
	Remark string
	Cover  string
	Quote  string //
	Size   uint64 //体积大小（字节）
	Limit  uint16 //照片最大数量
	Star   uint32 //点赞数
	Style  uint32
	Date   proxy.DurationInfo //日期
	Owner  string             //组织，机构
	Tags   []string
	Assets []string
}

func (mine *cacheContext) CreateCollAlbum(name, remark, user, group string, tp uint8, begin, stop int64) (*CollAlbumInfo, error) {
	db := new(nosql.Collective)
	db.UID = primitive.NewObjectID()
	db.ID = nosql.GetCollectiveNextID()
	db.CreatedTime = time.Now()
	db.Created = time.Now().Unix()
	db.Creator = user
	db.Name = name
	db.Remark = remark
	db.Cover = ""
	db.Size = 0
	db.Type = tp
	db.Group = group
	db.Date = proxy.DurationInfo{Start: begin, Stop: stop}
	db.MaxCount = config.Schema.Album.Group.MaxCount
	db.Assets = make([]string, 0, 1)
	db.Status = uint8(AlbumPrivate)
	db.Tags = make([]string, 0, 1)
	err := nosql.CreateCollective(db)
	if err != nil {
		return nil, err
	}
	info := new(CollAlbumInfo)
	info.initInfo(db)
	return info, nil
}

func (mine *cacheContext) GetCollAlbum(uid string) (*CollAlbumInfo, error) {
	if len(uid) < 2 {
		return nil, errors.New("the collective album uid is empty")
	}
	db, err := nosql.GetCollective(uid)
	if err != nil {
		return nil, err
	}
	info := new(CollAlbumInfo)
	info.initInfo(db)
	return info, nil
}

func (mine *cacheContext) GetCollAlbumByName(owner, name string) (*CollAlbumInfo, error) {
	if len(name) < 2 {
		return nil, errors.New("the collective album name is empty")
	}
	db, err := nosql.GetCollectiveByName(owner, name)
	if err != nil {
		return nil, err
	}
	info := new(CollAlbumInfo)
	info.initInfo(db)
	return info, nil
}

func (mine *cacheContext) HadCollAlbum(owner, name string) bool {
	info, _ := mine.GetCollAlbumByName(owner, name)
	if info == nil {
		return false
	}
	return true
}

func (mine *cacheContext) GetCollAlbums(user string) []*CollAlbumInfo {
	list := make([]*CollAlbumInfo, 0, 20)
	array, err := nosql.GetCollectivesByCreator(user)
	if err != nil {
		return list
	}
	for _, item := range array {
		info := new(CollAlbumInfo)
		info.initInfo(item)
		list = append(list, info)
	}
	return list
}

func (mine *cacheContext) GetCollAlbumsByGroup(uid string) []*CollAlbumInfo {
	array, err := nosql.GetCollectivesByGroup(uid)
	if err != nil {
		return make([]*CollAlbumInfo, 0, 0)
	}
	list := make([]*CollAlbumInfo, 0, len(array))
	for _, item := range array {
		info := new(CollAlbumInfo)
		info.initInfo(item)
		list = append(list, info)
	}
	return list
}

func (mine *cacheContext) GetCollAlbumsByGroups(array []string) []*CollAlbumInfo {
	albums := make([]*CollAlbumInfo, 0, 10)
	for _, group := range array {
		album := mine.GetCollAlbumsByGroup(group)
		if album != nil {
			albums = append(albums, album...)
		}
	}
	return albums
}

func (mine *cacheContext) GetCollAlbumList(array []string) []*CollAlbumInfo {
	if array == nil || len(array) < 1 {
		return make([]*CollAlbumInfo, 0, 0)
	}
	list := make([]*CollAlbumInfo, 0, len(array))
	for i := 0; i < len(array); i += 1 {
		db, err := nosql.GetCollective(array[i])
		if err == nil {
			info := new(CollAlbumInfo)
			info.initInfo(db)
			list = append(list, info)
		}
	}
	return list
}

func (mine *CollAlbumInfo) initInfo(db *nosql.Collective) {
	mine.Name = db.Name
	mine.UID = db.UID.Hex()
	mine.ID = db.ID
	mine.Remark = db.Remark
	mine.Created = db.Created
	mine.Updated = db.Updated
	mine.Creator = db.Creator
	mine.Cover = db.Cover
	mine.Limit = db.MaxCount
	mine.Size = db.Size
	mine.Owner = db.Group
	mine.Status = db.Status
	mine.Tags = db.Tags
	mine.Type = db.Type
	mine.Date = db.Date

	mine.Assets = db.Assets
	if mine.Assets == nil {
		mine.Assets = make([]string, 0, 1)
	}
}

func (mine *CollAlbumInfo) UpdateBase(name, remark, operator string) error {
	err := nosql.UpdateCollectiveBase(mine.UID, name, remark, operator)
	if err == nil {
		mine.Name = name
		mine.Remark = remark
		mine.Operator = operator
	}
	return err
}

func (mine *CollAlbumInfo) UpdateCover(cover, operator string) error {
	err := nosql.UpdateCollectiveCover(mine.UID, cover, operator)
	if err == nil {
		mine.Cover = cover
		mine.Operator = operator
	}
	return err
}

func (mine *CollAlbumInfo) UpdateStatus(operator string, st uint8) error {
	err := nosql.UpdateCollectiveStatus(mine.UID, operator, st)
	if err == nil {
		mine.Status = st
		mine.Operator = operator
	}
	return err
}

func (mine *CollAlbumInfo) Remove(operator string) error {
	if len(mine.Assets) > 0 {
		return errors.New("the album is not empty")
	}
	return nosql.RemoveCollective(mine.UID, operator)
}

/*
获取当前相册体积大小(kb)
*/
func (mine *CollAlbumInfo) GetSize() uint64 {
	return mine.Size / 1024
}

func (mine *CollAlbumInfo) UpdateSize(size uint64, operator string) error {
	err := nosql.UpdateAlbumSize(mine.UID, operator, size)
	if err == nil {
		mine.Size = size
	}
	return err
}

func (mine *CollAlbumInfo) appendSize(size uint64) error {
	sum := mine.Size + size
	err := nosql.UpdateAlbumSize(mine.UID, "", sum)
	if err == nil {
		mine.Size = sum
	}
	return err
}

func (mine *CollAlbumInfo) subtractSize(size uint64) error {
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

func (mine *CollAlbumInfo) AppendAssets(assets []string, operator string) error {
	if assets == nil || len(assets) < 1 {
		return errors.New("the assets is nil when append")
	}
	//all := make([]string, 0, len(mine.Assets)+len(assets))
	//all = append(all, mine.Assets...)
	for _, asset := range assets {
		if !mine.hadAsset(asset) {
			//all = append(all, asset)
			er := nosql.AppendCollectiveAsset(mine.UID, asset)
			if er != nil {
				return er
			}
		}
	}
	return nil
}

func (mine *CollAlbumInfo) UpdateAssets(assets []string, operator string) error {
	if assets == nil {
		assets = make([]string, 0, 1)
	}
	err := nosql.UpdateCollectiveAssets(mine.UID, operator, assets)
	if err == nil {
		mine.Assets = assets
		mine.Operator = operator
	}
	return err
}

func (mine *CollAlbumInfo) SubtractAssets(assets []string, operator string) error {
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

func (mine *CollAlbumInfo) hadAsset(asset string) bool {
	for _, node := range mine.Assets {
		if node == asset {
			return true
		}
	}
	return false
}

/**
从相册中删除一个照片
*/
func (mine *CollAlbumInfo) SubtractAsset(asset, operator string) error {
	if !mine.hadAsset(asset) {
		return nil
	}
	err := nosql.SubtractCollectiveAsset(mine.UID, asset)
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
