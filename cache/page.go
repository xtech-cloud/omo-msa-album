package cache

import (
	"errors"
	"fmt"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"omo.msa.album/proxy"
	"omo.msa.album/proxy/nosql"
	"time"
)

type PageInfo struct {
	Status uint8
	baseInfo
	Remark      string
	Owner       string //组织，机构
	Composition string
	Type        uint8  //类型
	Lifecycle   uint32 //播放时长
	Tags        []string
	Contents    []*proxy.PageContents
}

func (mine *cacheContext) CreatePage(name, remark, user, owner, composition string, tp, life uint32, tags []string) (*PageInfo, error) {
	db := new(nosql.Page)
	db.UID = primitive.NewObjectID()
	db.ID = nosql.GetPageNextID()
	db.Created = time.Now().Unix()
	db.Creator = user
	db.Name = name
	db.Remark = remark
	db.Owner = owner
	db.Type = uint8(tp)
	db.Lifecycle = life
	db.Composition = composition
	db.Tags = tags
	if db.Tags == nil {
		db.Tags = make([]string, 0, 1)
	}
	err := nosql.CreatePage(db)
	if err != nil {
		return nil, err
	}
	info := new(PageInfo)
	info.initInfo(db)
	return info, nil
}

func (mine *cacheContext) GetPage(uid string) (*PageInfo, error) {
	if len(uid) < 2 {
		return nil, errors.New("the collective album uid is empty")
	}
	db, err := nosql.GetPage(uid)
	if err != nil {
		return nil, err
	}
	info := new(PageInfo)
	info.initInfo(db)
	return info, nil
}

func (mine *cacheContext) GetPagesByOwner(uid string) []*PageInfo {
	array, err := nosql.GetPagesByOwner(uid)
	if err != nil {
		return make([]*PageInfo, 0, 0)
	}
	list := make([]*PageInfo, 0, len(array))
	for _, item := range array {
		info := new(PageInfo)
		info.initInfo(item)
		list = append(list, info)
	}
	return list
}

func (mine *cacheContext) GetPagesByType(owner string, tp uint32) []*PageInfo {
	array, err := nosql.GetPagesByType(owner, tp)
	if err != nil {
		return make([]*PageInfo, 0, 0)
	}
	list := make([]*PageInfo, 0, len(array))
	for _, item := range array {
		info := new(PageInfo)
		info.initInfo(item)
		list = append(list, info)
	}
	return list
}

func (mine *cacheContext) GetPagesByStatus(owner string, st uint32) []*PageInfo {
	array, err := nosql.GetPagesByStatus(owner, st)
	if err != nil {
		return make([]*PageInfo, 0, 0)
	}
	list := make([]*PageInfo, 0, len(array))
	for _, item := range array {
		info := new(PageInfo)
		info.initInfo(item)
		list = append(list, info)
	}
	return list
}

func (mine *cacheContext) GetPagesBySheet(uid string) []*PageInfo {
	sheet, err := mine.GetSheet(uid)
	if err != nil {
		return nil
	}
	list := make([]*PageInfo, 0, len(sheet.Pages))
	for _, page := range sheet.Pages {
		tmp, _ := cacheCtx.GetPage(page.UID)
		if tmp != nil && tmp.baseInfo.Deleted == 0 {
			list = append(list, tmp)
		}
	}
	return list
}

func (mine *cacheContext) GetPagesByList(arr []string) []*PageInfo {
	list := make([]*PageInfo, 0, len(arr))
	for _, page := range arr {
		tmp, _ := cacheCtx.GetPage(page)
		if tmp != nil {
			list = append(list, tmp)
		}
	}
	return list
}

func (mine *PageInfo) initInfo(db *nosql.Page) {
	mine.Name = db.Name
	mine.UID = db.UID.Hex()
	mine.ID = db.ID

	mine.Created = db.Created
	mine.Updated = db.Updated
	mine.Creator = db.Creator
	mine.Deleted = db.Deleted
	mine.Operator = db.Operator
	mine.Remark = db.Remark
	mine.Owner = db.Owner
	mine.Type = db.Type
	mine.Composition = db.Composition
	mine.Lifecycle = db.Lifecycle
	mine.Tags = db.Tags

	mine.Contents = db.Contents
	if mine.Contents == nil {
		mine.Contents = make([]*proxy.PageContents, 0, 1)
	}
}

func (mine *PageInfo) UpdateBase(name, remark, operator string, life uint32) error {
	err := nosql.UpdatePageBase(mine.UID, name, remark, operator, life)
	if err == nil {
		mine.Name = name
		mine.Remark = remark
		mine.Operator = operator
		mine.Lifecycle = life
		mine.Updated = time.Now().Unix()
	}
	return err
}

func (mine *PageInfo) Remove(operator string) error {
	sheets := cacheCtx.GetSheetsByPage(mine.UID)
	if len(sheets) > 0 {
		return errors.New(fmt.Sprintf("the page had used by sheets count = %d", len(sheets)))
	}
	return nosql.RemovePage(mine.UID, operator)
}

func (mine *PageInfo) SetContent(info *proxy.PageContents, operator string) error {
	if info == nil {
		return errors.New("the info is nil when setContent")
	}
	all := make([]*proxy.PageContents, 0, len(mine.Contents))
	all = append(all, mine.Contents...)
	for i, item := range all {
		if item.Slot == info.Slot {
			if i == len(all)-1 {
				all = append(all[:i])
			} else {
				all = append(all[:i], all[i+1:]...)
			}
			break
		}
	}
	all = append(all, info)
	return mine.UpdateContents(all, operator)
}

func (mine *PageInfo) AppendContents(list []*proxy.PageContents, operator string) error {
	if list == nil || len(list) < 1 {
		return errors.New("the assets is nil when append")
	}
	all := make([]*proxy.PageContents, 0, len(list)+len(mine.Contents))
	all = append(all, mine.Contents...)
	for _, item := range list {
		if !mine.hadContent(item.Slot) {
			all = append(all, item)
		}
	}
	return mine.UpdateContents(all, operator)
}

func (mine *PageInfo) UpdateContents(list []*proxy.PageContents, operator string) error {
	if list == nil {
		return nil
	}
	err := nosql.UpdatePageContents(mine.UID, operator, list)
	if err == nil {
		mine.Contents = list
		mine.Operator = operator
		mine.Updated = time.Now().Unix()
	}
	return err
}

func (mine *PageInfo) SubtractContents(list []*proxy.PageContents, operator string) error {
	if list == nil || len(list) < 1 {
		return nil
	}
	all := make([]*proxy.PageContents, 0, len(mine.Contents))
	all = append(all, mine.Contents...)
	for _, asset := range list {
		for i := 0; i < len(all); i += 1 {
			if all[i].Slot == asset.Slot {
				if i == len(all)-1 {
					all = append(all[:i])
				} else {
					all = append(all[:i], all[i+1:]...)
				}
				break
			}
		}
	}
	return mine.UpdateContents(all, operator)
}

func (mine *PageInfo) hadContent(slot uint32) bool {
	for _, node := range mine.Contents {
		if node.Slot == slot {
			return true
		}
	}
	return false
}

//删除一个位置的内容
func (mine *PageInfo) SubtractContent(operator string, slot uint32) error {
	if !mine.hadContent(slot) {
		return nil
	}
	err := nosql.SubtractPageContent(mine.UID, slot)
	if err == nil {
		for i := 0; i < len(mine.Contents); i += 1 {
			if mine.Contents[i].Slot == slot {
				if i == len(mine.Contents)-1 {
					mine.Contents = append(mine.Contents[:i])
				} else {
					mine.Contents = append(mine.Contents[:i], mine.Contents[i+1:]...)
				}
				break
			}
		}
		mine.Operator = operator
	}
	return err
}
