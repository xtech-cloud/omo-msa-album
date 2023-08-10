package cache

import (
	"errors"
	pb "github.com/xtech-cloud/omo-msp-album/proto/album"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"omo.msa.album/proxy"
	"omo.msa.album/proxy/nosql"
	"time"
)

// 文件夹或者包
type FolderInfo struct {
	baseInfo
	Remark   string
	Owner    string
	Parent   string
	Cover    string
	Access   uint8
	Tags     []string
	Contents []*proxy.FolderContent
}

func (mine *cacheContext) CreateFolder(name, remark, user, owner, cover, parent string, tags []string, list []*pb.PairInt) (*FolderInfo, error) {
	db := new(nosql.Folder)
	db.UID = primitive.NewObjectID()
	db.ID = nosql.GetFolderNextID()
	db.Created = time.Now().Unix()
	db.Creator = user
	db.Name = name
	db.Remark = remark
	db.Owner = owner
	db.Parent = parent
	db.Cover = cover
	db.Access = 0
	db.Tags = tags
	if db.Tags == nil {
		db.Tags = make([]string, 0, 1)
	}

	db.Contents = make([]*proxy.FolderContent, 0, len(list))
	for _, pair := range list {
		db.Contents = append(db.Contents, &proxy.FolderContent{UID: pair.Value, Type: pair.Key})
	}

	err := nosql.CreateFolder(db)
	if err != nil {
		return nil, err
	}
	info := new(FolderInfo)
	info.initInfo(db)
	return info, nil
}

func (mine *cacheContext) GetFolder(uid string) (*FolderInfo, error) {
	if len(uid) < 2 {
		return nil, errors.New("the album uid is empty")
	}
	db, err := nosql.GetFolder(uid)
	if err != nil {
		return nil, err
	}
	info := new(FolderInfo)
	info.initInfo(db)
	return info, nil
}

func (mine *cacheContext) GetFoldersByOwner(uid string) []*FolderInfo {
	list := make([]*FolderInfo, 0, 20)
	if len(uid) < 2 {
		return nil
	}
	array, err := nosql.GetFoldersByOwner(uid)
	if err != nil {
		return list
	}
	for _, item := range array {
		info := new(FolderInfo)
		info.initInfo(item)
		list = append(list, info)
	}
	return list
}

func (mine *cacheContext) GetFoldersByParent(uid string) []*FolderInfo {
	list := make([]*FolderInfo, 0, 20)
	if len(uid) < 2 {
		return nil
	}
	array, err := nosql.GetFoldersByParent(uid)
	if err != nil {
		return list
	}
	for _, item := range array {
		info := new(FolderInfo)
		info.initInfo(item)
		list = append(list, info)
	}
	return list
}

func (mine *FolderInfo) initInfo(db *nosql.Folder) {
	mine.Name = db.Name
	mine.UID = db.UID.Hex()
	mine.ID = db.ID
	mine.Created = db.Created
	mine.Creator = db.Creator
	mine.Operator = db.Operator
	mine.Updated = db.Updated

	mine.Remark = db.Remark
	mine.Owner = db.Owner
	mine.Parent = db.Parent
	mine.Cover = db.Cover
	mine.Access = db.Access

	mine.Tags = db.Tags
	mine.Contents = db.Contents
}

func (mine *FolderInfo) UpdateBase(name, remark, operator, parent string) error {
	err := nosql.UpdateFolderBase(mine.UID, name, remark, operator, parent)
	if err == nil {
		mine.Name = name
		mine.Remark = remark
		mine.Operator = operator
		mine.Parent = parent
		mine.Updated = time.Now().Unix()
	}
	return err
}

func (mine *FolderInfo) UpdateCover(cover, operator string) error {
	err := nosql.UpdateFolderCover(mine.UID, cover, operator)
	if err == nil {
		mine.Cover = cover
		mine.Operator = operator
		mine.Updated = time.Now().Unix()
	}
	return err
}

func (mine *FolderInfo) Remove(operator string) error {
	return nosql.RemoveFolder(mine.UID, operator)
}

func (mine *FolderInfo) UpdateContents(operator string, list []*proxy.FolderContent) error {
	if list == nil {
		list = make([]*proxy.FolderContent, 0, 10)
	}
	err := nosql.UpdateFolderContents(mine.UID, operator, list)
	if err == nil {
		mine.Contents = list
		mine.Operator = operator
	}
	return err
}

func (mine *FolderInfo) AppendAssets(assets []*proxy.FolderContent, operator string) error {
	if assets == nil || len(assets) < 1 {
		return errors.New("the assets is nil when append")
	}
	all := make([]*proxy.FolderContent, 0, len(mine.Contents)+len(assets))
	all = append(all, mine.Contents...)
	for _, asset := range assets {
		if !mine.hadContent(asset.UID) {
			all = append(all, asset)
		}
	}
	return mine.UpdateContents(operator, all)
}

func (mine *FolderInfo) RemoveAssets(list []string, operator string) error {
	if list == nil || len(list) < 1 {
		return nil
	}
	all := make([]*proxy.FolderContent, 0, len(mine.Contents))
	all = append(all, mine.Contents...)
	for _, asset := range list {
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
	return mine.UpdateContents(operator, all)
}

func (mine *FolderInfo) hadContent(asset string) bool {
	for _, node := range mine.Contents {
		if node.UID == asset {
			return true
		}
	}
	return false
}

/**
从目录中删除一个内容
*/
func (mine *FolderInfo) SubtractAsset(asset, operator string) error {
	if !mine.hadContent(asset) {
		return nil
	}
	err := nosql.SubtractFolderContent(mine.UID, asset)
	if err == nil {
		for i := 0; i < len(mine.Contents); i += 1 {
			if mine.Contents[i].UID == asset {
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
