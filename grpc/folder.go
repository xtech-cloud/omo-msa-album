package grpc

//
//import (
//	"context"
//	"errors"
//	"fmt"
//	pb "github.com/xtech-cloud/omo-msp-album/proto/album"
//	pbstatus "github.com/xtech-cloud/omo-msp-status/proto/status"
//	"omo.msa.album/cache"
//	"omo.msa.album/proxy"
//)
//
//type FolderService struct{}
//
//func switchFolder(info *cache.FolderInfo) *pb.FolderInfo {
//	tmp := new(pb.FolderInfo)
//	tmp.Uid = info.UID
//	tmp.Id = info.ID
//	tmp.Created = info.Created
//	tmp.Updated = info.Updated
//	tmp.Operator = info.Operator
//	tmp.Creator = info.Creator
//	tmp.Name = info.Name
//	tmp.Remark = info.Remark
//	tmp.Cover = info.Cover
//	tmp.Access = uint32(info.Access)
//	tmp.Parent = info.Parent
//	tmp.Owner = info.Owner
//	tmp.Tags = info.Tags
//	tmp.Contents = switchFolderContents(info.Contents)
//	return tmp
//}
//
//func switchFolderContents(list []*proxy.FolderContent) []*pb.PairInt {
//	slots := make([]*pb.PairInt, 0, len(list))
//	for _, slot := range list {
//		item := &pb.PairInt{Key: slot.Type, Value: slot.UID}
//		slots = append(slots, item)
//	}
//	return slots
//}
//
//func (mine *FolderService) AddOne(ctx context.Context, in *pb.ReqFolderAdd, out *pb.ReplyFolderInfo) error {
//	path := "folder.addOne"
//	inLog(path, in)
//	if len(in.Name) < 1 {
//		out.Status = outError(path, "the name is empty ", pbstatus.ResultStatus_Empty)
//		return nil
//	}
//
//	info, err := cache.Context().CreateFolder(in.Name, in.Remark, in.Operator, in.Owner, in.Cover, in.Parent, in.Tags)
//	if err != nil {
//		out.Status = outError(path, err.Error(), pbstatus.ResultStatus_DBException)
//		return nil
//	}
//	out.Info = switchFolder(info)
//	out.Status = outLog(path, out)
//	return nil
//}
//
//func (mine *FolderService) GetOne(ctx context.Context, in *pb.RequestInfo, out *pb.ReplyFolderInfo) error {
//	path := "folder.getOne"
//	inLog(path, in)
//	if len(in.Uid) < 1 {
//		out.Status = outError(path, "the uid is empty ", pbstatus.ResultStatus_Empty)
//		return nil
//	}
//	info, er := cache.Context().GetFolder(in.Uid)
//	if er != nil {
//		out.Status = outError(path, "the folder not found ", pbstatus.ResultStatus_NotExisted)
//		return nil
//	}
//	out.Info = switchFolder(info)
//	out.Status = outLog(path, out)
//	return nil
//}
//
//func (mine *FolderService) GetStatistic(ctx context.Context, in *pb.RequestFilter, out *pb.ReplyStatistic) error {
//	path := "folder.getStatistic"
//	inLog(path, in)
//	if len(in.Field) < 1 {
//		out.Status = outError(path, "the user is empty ", pbstatus.ResultStatus_Empty)
//		return nil
//	}
//
//	out.Status = outLog(path, out)
//	return nil
//}
//
//func (mine *FolderService) RemoveOne(ctx context.Context, in *pb.RequestInfo, out *pb.ReplyInfo) error {
//	path := "folder.remove"
//	inLog(path, in)
//	if len(in.Uid) < 1 {
//		out.Status = outError(path, "the uid is empty ", pbstatus.ResultStatus_Empty)
//		return nil
//	}
//	info, er := cache.Context().GetFolder(in.Uid)
//	if er != nil {
//		out.Status = outError(path, "the folder not found ", pbstatus.ResultStatus_NotExisted)
//		return nil
//	}
//	err := info.Remove(in.Operator)
//	if err != nil {
//		out.Status = outError(path, err.Error(), pbstatus.ResultStatus_DBException)
//		return nil
//	}
//	out.Uid = in.Uid
//	out.Status = outLog(path, out)
//	return nil
//}
//
//func (mine *FolderService) Search(ctx context.Context, in *pb.RequestInfo, out *pb.ReplyFolderList) error {
//	path := "folder.search"
//	inLog(path, in)
//
//	out.Status = outLog(path, fmt.Sprintf("the length = %d", len(out.List)))
//	return nil
//}
//
//func (mine *FolderService) GetListBy(ctx context.Context, in *pb.RequestFilter, out *pb.ReplyFolderList) error {
//	path := "folder.getListByFilter"
//	inLog(path, in)
//	var list []*cache.FolderInfo
//	var err error
//	if in.Field == "" {
//		list = cache.Context().GetFoldersByOwner(in.Owner)
//	} else {
//		err = errors.New("the key not defined")
//	}
//	if err != nil {
//		out.Status = outError(path, err.Error(), pbstatus.ResultStatus_DBException)
//		return nil
//	}
//	out.List = make([]*pb.FolderInfo, 0, len(list))
//	for _, value := range list {
//		out.List = append(out.List, switchFolder(value))
//	}
//
//	out.Status = outLog(path, fmt.Sprintf("the length = %d", len(out.List)))
//	return nil
//}
//
//func (mine *FolderService) UpdateBase(ctx context.Context, in *pb.ReqFolderUpdate, out *pb.ReplyInfo) error {
//	path := "folder.updateBase"
//	inLog(path, in)
//	if len(in.Uid) < 1 {
//		out.Status = outError(path, "the uid is empty ", pbstatus.ResultStatus_Empty)
//		return nil
//	}
//	info, er := cache.Context().GetFolder(in.Uid)
//	if er != nil {
//		out.Status = outError(path, er.Error(), pbstatus.ResultStatus_NotExisted)
//		return nil
//	}
//	var err error
//	err = info.UpdateBase(in.Name, in.Remark, in.Operator, cache.AspectType(in.Inspect))
//	if err != nil {
//		out.Status = outError(path, err.Error(), pbstatus.ResultStatus_DBException)
//		return nil
//	}
//
//	out.Status = outLog(path, out)
//	return nil
//}
//
//func (mine *FolderService) UpdateBy(ctx context.Context, in *pb.RequestUpdate, out *pb.ReplyInfo) error {
//	path := "folder.updateByFilter"
//	inLog(path, in)
//	if len(in.Uid) < 1 {
//		out.Status = outError(path, "the uid is empty ", pbstatus.ResultStatus_Empty)
//		return nil
//	}
//	info, er := cache.Context().GetFolder(in.Uid)
//	if er != nil {
//		out.Status = outError(path, "the folder not found ", pbstatus.ResultStatus_NotExisted)
//		return nil
//	}
//	var err error
//	if in.Field == "cover" {
//		err = info.UpdateCover(in.Value, in.Operator)
//	} else {
//		err = errors.New("the field not defined")
//	}
//	if err != nil {
//		out.Status = outError(path, err.Error(), pbstatus.ResultStatus_DBException)
//		return nil
//	}
//	out.Status = outLog(path, out)
//	return nil
//}
//
//func (mine *FolderService) AppendContents(ctx context.Context, in *pb.ReqFolderAppend, out *pb.ReplyFolderContents) error {
//	path := "folder.appendAsset"
//	inLog(path, in)
//	if len(in.Uid) < 1 {
//		out.Status = outError(path, "the uid is empty ", pbstatus.ResultStatus_Empty)
//		return nil
//	}
//	info, er := cache.Context().GetFolder(in.Uid)
//	if er != nil {
//		out.Status = outError(path, er.Error(), pbstatus.ResultStatus_NotExisted)
//		return nil
//	}
//
//	err := info.AppendAssets(in.List, in.Operator)
//	if err != nil {
//		out.Status = outError(path, err.Error(), pbstatus.ResultStatus_DBException)
//		return nil
//	}
//	out.Uid = in.Uid
//	out.List = switchFolderContents(info.Contents)
//	out.Status = outLog(path, out)
//	return nil
//}
//
//func (mine *FolderService) SubtractContents(ctx context.Context, in *pb.RequestList, out *pb.ReplyFolderContents) error {
//	path := "folder.subtractAsset"
//	inLog(path, in)
//	if len(in.Uid) < 1 {
//		out.Status = outError(path, "the uid is empty ", pbstatus.ResultStatus_Empty)
//		return nil
//	}
//	info, er := cache.Context().GetFolder(in.Uid)
//	if er != nil {
//		out.Status = outError(path, er.Error(), pbstatus.ResultStatus_NotExisted)
//		return nil
//	}
//
//	err := info.RemoveAssets(in.List, in.Operator)
//	if err != nil {
//		out.Status = outError(path, err.Error(), pbstatus.ResultStatus_DBException)
//		return nil
//	}
//	out.List = switchFolderContents(info.Contents)
//	out.Status = outLog(path, out)
//	return nil
//}
