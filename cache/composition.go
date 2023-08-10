package cache

import (
	"errors"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"omo.msa.album/proxy"
	"omo.msa.album/proxy/nosql"
	"time"
)

const (
	Aspect169  AspectType = 1
	Aspect1610 AspectType = 2
)

type AspectType uint8

type CompositionInfo struct {
	baseInfo
	Remark string
	Owner  string     //组织，机构
	Aspect AspectType //屏幕纵横比
	Cover  string
	Tags   []string
	Slots  []*proxy.SlotInfo
}

func (mine *cacheContext) CreateComposition(name, remark, user, owner, cover string, aspect AspectType, tags []string) (*CompositionInfo, error) {
	db := new(nosql.Composition)
	db.UID = primitive.NewObjectID()
	db.ID = nosql.GetCompositionNextID()
	db.Created = time.Now().Unix()
	db.Creator = user
	db.Name = name
	db.Remark = remark
	db.Cover = cover
	db.Owner = owner
	db.Aspect = uint8(aspect)
	db.Slots = make([]*proxy.SlotInfo, 0, 1)
	if tags == nil {
		db.Tags = make([]string, 0, 1)
	} else {
		db.Tags = tags
	}

	err := nosql.CreateComposition(db)
	if err != nil {
		return nil, err
	}
	info := new(CompositionInfo)
	info.initInfo(db)
	return info, nil
}

func (mine *cacheContext) GetComposition(uid string) (*CompositionInfo, error) {
	if len(uid) < 2 {
		return nil, errors.New("the composition uid is empty")
	}
	db, err := nosql.GetComposition(uid)
	if err != nil {
		return nil, err
	}
	info := new(CompositionInfo)
	info.initInfo(db)
	return info, nil
}

func (mine *cacheContext) GetCompositionsByOwner(uid string) []*CompositionInfo {
	array, err := nosql.GetCompositionsByOwner(uid)
	if err != nil {
		return make([]*CompositionInfo, 0, 0)
	}
	list := make([]*CompositionInfo, 0, len(array))
	for _, item := range array {
		info := new(CompositionInfo)
		info.initInfo(item)
		list = append(list, info)
	}
	return list
}

func (mine *cacheContext) GetCompositionList(array []string) []*CompositionInfo {
	if array == nil || len(array) < 1 {
		return make([]*CompositionInfo, 0, 0)
	}
	list := make([]*CompositionInfo, 0, len(array))
	for i := 0; i < len(array); i += 1 {
		db, err := nosql.GetComposition(array[i])
		if err == nil {
			info := new(CompositionInfo)
			info.initInfo(db)
			list = append(list, info)
		}
	}
	return list
}

func (mine *CompositionInfo) initInfo(db *nosql.Composition) {
	mine.Name = db.Name
	mine.UID = db.UID.Hex()
	mine.ID = db.ID
	mine.Created = db.Created
	mine.Updated = db.Updated
	mine.Creator = db.Creator
	mine.Operator = db.Operator
	mine.Remark = db.Remark
	mine.Cover = db.Cover
	mine.Aspect = AspectType(db.Aspect)
	mine.Owner = db.Owner
	mine.Tags = db.Tags
	mine.Slots = db.Slots
}

func (mine *CompositionInfo) UpdateBase(name, remark, operator string, aspect AspectType) error {
	err := nosql.UpdateCompositionBase(mine.UID, name, remark, operator, uint8(aspect))
	if err == nil {
		mine.Name = name
		mine.Remark = remark
		mine.Aspect = aspect
		mine.Operator = operator
		mine.Updated = time.Now().Unix()
	}
	return err
}

func (mine *CompositionInfo) UpdateCover(cover, operator string) error {
	err := nosql.UpdateCompositionCover(mine.UID, cover, operator)
	if err == nil {
		mine.Cover = cover
		mine.Operator = operator
		mine.Updated = time.Now().Unix()
	}
	return err
}

func (mine *CompositionInfo) Remove(operator string) error {
	return nosql.RemoveComposition(mine.UID, operator)
}

func (mine *CompositionInfo) SetSlot(slot *proxy.SlotInfo, operator string) error {
	if slot == nil {
		return errors.New("the assets is nil when append")
	}
	all := make([]*proxy.SlotInfo, 0, len(mine.Slots)+1)
	all = append(all, mine.Slots...)
	for i, item := range all {
		if item.Index == slot.Index {
			if i == len(all)-1 {
				all = append(all[:i])
			} else {
				all = append(all[:i], all[i+1:]...)
			}
			break
		}
	}
	all = append(all, slot)
	return mine.UpdateSlots(all, operator)
}

func (mine *CompositionInfo) UpdateSlots(slots []*proxy.SlotInfo, operator string) error {
	if slots == nil {
		return nil
	}
	err := nosql.UpdateCompositionSlots(mine.UID, operator, slots)
	if err == nil {
		mine.Slots = slots
		mine.Operator = operator
		mine.Updated = time.Now().Unix()
	}
	return err
}
