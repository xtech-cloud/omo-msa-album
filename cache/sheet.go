package cache

import (
	"errors"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"omo.msa.album/proxy"
	"omo.msa.album/proxy/nosql"
	"time"
)

type SheetInfo struct {
	Status uint8 //激活状态，只有紧急的可以被激活，而且同时只有一个
	Type   uint8 //普通播放列表， 紧急播放列表
	Size   uint32
	baseInfo
	Remark string
	Owner  string //组织，机构
	Cover  string

	Aspect string
	Target string //播放的区域
	Tags   []string
	Pages  []*proxy.SheetPage
}

func (mine *cacheContext) CreateSheet(name, remark, user, owner, aspect string, tags []string) (*SheetInfo, error) {
	db := new(nosql.Sheet)
	db.UID = primitive.NewObjectID()
	db.ID = nosql.GetSheetNextID()
	db.Created = time.Now().Unix()
	db.Creator = user
	db.Name = name
	db.Remark = remark
	db.Cover = ""
	db.Owner = owner
	db.Tags = tags
	db.Size = 0
	db.Aspect = aspect
	if db.Tags == nil {
		db.Tags = make([]string, 0, 1)
	}
	db.Pages = make([]*proxy.SheetPage, 0, 1)
	err := nosql.CreateSheet(db)
	if err != nil {
		return nil, err
	}
	info := new(SheetInfo)
	info.initInfo(db)
	return info, nil
}

func (mine *cacheContext) GetSheet(uid string) (*SheetInfo, error) {
	if len(uid) < 2 {
		return nil, errors.New("the collective album uid is empty")
	}
	db, err := nosql.GetSheet(uid)
	if err != nil {
		return nil, err
	}
	info := new(SheetInfo)
	info.initInfo(db)
	return info, nil
}

func (mine *cacheContext) GetSheetsByOwner(uid string) []*SheetInfo {
	array, err := nosql.GetSheetsByOwner(uid)
	if err != nil {
		return make([]*SheetInfo, 0, 0)
	}
	list := make([]*SheetInfo, 0, len(array))
	for _, item := range array {
		info := new(SheetInfo)
		info.initInfo(item)
		list = append(list, info)
	}
	return list
}

func (mine *cacheContext) GetSheetsByTarget(uid string) []*SheetInfo {
	array, err := nosql.GetSheetsByTarget(uid)
	if err != nil {
		return make([]*SheetInfo, 0, 0)
	}
	list := make([]*SheetInfo, 0, len(array))
	for _, item := range array {
		info := new(SheetInfo)
		info.initInfo(item)
		list = append(list, info)
	}
	return list
}

func (mine *cacheContext) GetSheetsByPage(uid string) []*SheetInfo {
	array, err := nosql.GetSheetsByPage(uid)
	if err != nil {
		return make([]*SheetInfo, 0, 0)
	}
	list := make([]*SheetInfo, 0, len(array))
	for _, item := range array {
		info := new(SheetInfo)
		info.initInfo(item)
		list = append(list, info)
	}
	return list
}

func (mine *SheetInfo) initInfo(db *nosql.Sheet) {
	mine.Name = db.Name
	mine.UID = db.UID.Hex()
	mine.ID = db.ID
	mine.Remark = db.Remark
	mine.Created = db.Created
	mine.Updated = db.Updated
	mine.Creator = db.Creator
	mine.Cover = db.Cover
	mine.Owner = db.Owner
	mine.Target = db.Target
	mine.Aspect = db.Aspect
	mine.Tags = db.Tags

	mine.Pages = db.Pages
	if mine.Pages == nil {
		mine.Pages = make([]*proxy.SheetPage, 0, 1)
	}
}

func (mine *SheetInfo) UpdateBase(name, remark, operator string) error {
	err := nosql.UpdateSheetBase(mine.UID, name, remark, operator)
	if err == nil {
		mine.Name = name
		mine.Remark = remark
		mine.Operator = operator
		mine.Updated = time.Now().Unix()
	}
	return err
}

func (mine *SheetInfo) UpdateCover(cover, operator string) error {
	err := nosql.UpdateSheetCover(mine.UID, cover, operator)
	if err == nil {
		mine.Cover = cover
		mine.Operator = operator
		mine.Updated = time.Now().Unix()
	}
	return err
}

func (mine *SheetInfo) UpdateTarget(tar, operator string) error {
	err := nosql.UpdateSheetTarget(mine.UID, tar, operator)
	if err == nil {
		mine.Target = tar
		mine.Operator = operator
		mine.Updated = time.Now().Unix()
	}
	return err
}

func (mine *SheetInfo) GetAssetCount() uint32 {
	num := 0
	for _, page := range mine.Pages {
		tmp, _ := cacheCtx.GetPage(page.UID)
		if tmp != nil {
			for _, content := range tmp.Contents {
				num += len(content.List)
			}
		}
	}
	return uint32(num)
}

func (mine *SheetInfo) Remove(operator string) error {
	//if len(mine.Pages) > 0 {
	//	return errors.New("the sheet is not empty")
	//}
	return nosql.RemoveSheet(mine.UID, operator)
}

func (mine *SheetInfo) AppendPages(list []*proxy.SheetPage, operator string) error {
	if list == nil || len(list) < 1 {
		return errors.New("the pages is nil when append")
	}
	all := make([]*proxy.SheetPage, 0, len(mine.Pages)+len(list))
	all = append(all, mine.Pages...)
	for _, asset := range list {
		if !mine.hadPage(asset.UID) {
			all = append(all, asset)
		}
	}
	return mine.UpdatePages(all, operator)
}

func (mine *SheetInfo) UpdatePages(assets []*proxy.SheetPage, operator string) error {
	if assets == nil {
		return nil
	}
	size := uint32(0)
	for _, asset := range assets {
		size += asset.Weight
	}
	err := nosql.UpdateSheetPages(mine.UID, operator, assets)
	if err == nil {
		mine.Pages = assets
		mine.Operator = operator
		mine.Updated = time.Now().Unix()
		mine.Size = size
	}
	return err
}

func (mine *SheetInfo) SubtractAssets(assets []string, operator string) error {
	if assets == nil || len(assets) < 1 {
		return nil
	}
	all := make([]*proxy.SheetPage, 0, len(mine.Pages))
	all = append(all, mine.Pages...)
	for _, asset := range assets {
		for i := 0; i < len(all); i += 1 {
			if all[i].UID == asset {
				if i == len(all)-1 {
					all = append(all[:i])
				} else {
					all = append(all[:i], all[i+1:]...)
				}
				break
			}
		}
	}
	return mine.UpdatePages(all, operator)
}

func (mine *SheetInfo) hadPage(uid string) bool {
	for _, node := range mine.Pages {
		if node.UID == uid {
			return true
		}
	}
	return false
}

/**
从相册中删除一个照片
*/
func (mine *SheetInfo) SubtractAsset(page, operator string) error {
	if !mine.hadPage(page) {
		return nil
	}
	err := nosql.SubtractSheetPage(mine.UID, page)
	if err == nil {
		for i := 0; i < len(mine.Pages); i += 1 {
			if mine.Pages[i].UID == page {
				if i == len(mine.Pages)-1 {
					mine.Pages = append(mine.Pages[:i])
				} else {
					mine.Pages = append(mine.Pages[:i], mine.Pages[i+1:]...)
				}
				break
			}
		}
		mine.Operator = operator
	}
	return err
}
