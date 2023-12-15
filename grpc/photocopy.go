package grpc

import (
	"context"
	"errors"
	"fmt"
	pb "github.com/xtech-cloud/omo-msp-album/proto/album"
	pbstatus "github.com/xtech-cloud/omo-msp-status/proto/status"
	"omo.msa.album/cache"
	"omo.msa.album/proxy"
)

type PhotocopyService struct{}

func switchPhotocopy(info *cache.PhotocopyInfo) *pb.PhotocopyInfo {
	tmp := new(pb.PhotocopyInfo)
	tmp.Uid = info.UID
	tmp.Id = info.ID
	tmp.Created = info.Created
	tmp.Updated = info.Updated
	tmp.Operator = info.Operator
	tmp.Creator = info.Creator
	tmp.Name = info.Name
	tmp.Remark = info.Remark
	tmp.Owner = info.Owner
	tmp.Mother = info.Mother
	tmp.Template = info.Template
	tmp.Count = info.Count
	tmp.Tags = info.Tags
	tmp.Slots = switchSlots(info.Slots)
	return tmp
}

func switchSlots(list []proxy.PhotocopySlot) []*pb.PhotocopySlot {
	arr := make([]*pb.PhotocopySlot, 0, len(list))
	for _, slot := range list {
		tmp := new(pb.PhotocopySlot)
		tmp.Index = uint32(slot.Index)
		tmp.Remark = slot.Remark
		tmp.Page = uint32(slot.Page)
		tmp.Role = uint32(slot.Role)
		tmp.Background = slot.Background
		tmp.Asset = slot.Asset
		arr = append(arr, tmp)
	}
	return arr
}

func switchStyleSlots(list []proxy.StyleSlot) []*pb.StyleSlot {
	arr := make([]*pb.StyleSlot, 0, len(list))
	for _, slot := range list {
		tmp := new(pb.StyleSlot)
		tmp.Key = slot.Key
		tmp.X = slot.X
		tmp.Y = slot.Y
		tmp.Width = slot.Width
		tmp.Height = slot.Height
		tmp.Bold = slot.Bold
		tmp.Size = slot.Size
		arr = append(arr, tmp)
	}
	return arr
}

func (mine *PhotocopyService) AddOne(ctx context.Context, in *pb.ReqPhotocopyAdd, out *pb.ReplyPhotocopyInfo) error {
	path := "photocopy.add"
	inLog(path, in)
	if len(in.Name) < 1 {
		out.Status = outError(path, "the name is empty ", pbstatus.ResultStatus_Empty)
		return nil
	}

	info, err := cache.Context().CreatePhotocopy(in.Name, in.Remark, in.Operator, in.Template, in.Owner)
	if err != nil {
		out.Status = outError(path, err.Error(), pbstatus.ResultStatus_DBException)
		return nil
	}
	out.Info = switchPhotocopy(info)
	out.Status = outLog(path, out)
	return nil
}

func (mine *PhotocopyService) GetOne(ctx context.Context, in *pb.RequestInfo, out *pb.ReplyPhotocopyInfo) error {
	path := "photocopy.getOne"
	inLog(path, in)
	if len(in.Uid) < 1 {
		out.Status = outError(path, "the uid is empty ", pbstatus.ResultStatus_Empty)
		return nil
	}
	var info *cache.PhotocopyInfo
	var er error
	if in.Flag == "" {
		info, er = cache.Context().GetPhotocopy(in.Uid)
	} else if in.Flag == "clone" {
		info, er = cache.Context().ClonePhotocopy(in.Uid, in.User)
	}

	if er != nil {
		out.Status = outError(path, "the photocopy not found ", pbstatus.ResultStatus_NotExisted)
		return nil
	}
	out.Info = switchPhotocopy(info)
	out.Status = outLog(path, out)
	return nil
}

func (mine *PhotocopyService) GetStatistic(ctx context.Context, in *pb.RequestFilter, out *pb.ReplyStatistic) error {
	path := "photocopy.getStatistic"
	inLog(path, in)
	if len(in.Field) < 1 {
		out.Status = outError(path, "the user is empty ", pbstatus.ResultStatus_Empty)
		return nil
	}

	out.Status = outLog(path, out)
	return nil
}

func (mine *PhotocopyService) RemoveOne(ctx context.Context, in *pb.RequestInfo, out *pb.ReplyInfo) error {
	path := "photocopy.remove"
	inLog(path, in)
	if len(in.Uid) < 1 {
		out.Status = outError(path, "the uid is empty ", pbstatus.ResultStatus_Empty)
		return nil
	}
	info, er := cache.Context().GetPhotocopy(in.Uid)
	if er != nil {
		out.Status = outError(path, "the photocopy not found ", pbstatus.ResultStatus_NotExisted)
		return nil
	}
	er = info.Remove(in.Operator)
	if er != nil {
		out.Status = outError(path, "the photocopy not found ", pbstatus.ResultStatus_NotExisted)
		return nil
	}

	out.Uid = in.Uid
	out.Status = outLog(path, out)
	return nil
}

func (mine *PhotocopyService) Search(ctx context.Context, in *pb.RequestInfo, out *pb.ReplyPhotocopyList) error {
	path := "photocopy.search"
	inLog(path, in)

	out.Status = outLog(path, fmt.Sprintf("the length = %d", len(out.List)))
	return nil
}

func (mine *PhotocopyService) GetListByFilter(ctx context.Context, in *pb.RequestFilter, out *pb.ReplyPhotocopyList) error {
	path := "photocopy.getListByFilter"
	inLog(path, in)
	var list []*cache.PhotocopyInfo
	var err error
	if in.Field == "" {
		list, err = cache.Context().GetPhotocopiesByMaster(in.Owner)
	} else {
		err = errors.New("the key not defined")
	}
	if err != nil {
		out.Status = outError(path, err.Error(), pbstatus.ResultStatus_DBException)
		return nil
	}
	out.List = make([]*pb.PhotocopyInfo, 0, len(list))
	for _, value := range list {
		out.List = append(out.List, switchPhotocopy(value))
	}

	out.Status = outLog(path, fmt.Sprintf("the length = %d", len(out.List)))
	return nil
}

func (mine *PhotocopyService) UpdateBase(ctx context.Context, in *pb.ReqPhotocopyUpdate, out *pb.ReplyInfo) error {
	path := "photocopy.updateBase"
	inLog(path, in)
	if len(in.Uid) < 1 {
		out.Status = outError(path, "the uid is empty ", pbstatus.ResultStatus_Empty)
		return nil
	}
	info, er := cache.Context().GetPhotocopy(in.Uid)
	if er != nil {
		out.Status = outError(path, er.Error(), pbstatus.ResultStatus_NotExisted)
		return nil
	}
	var err error
	err = info.UpdateBase(in.Name, in.Remark, in.Operator)
	if err != nil {
		out.Status = outError(path, err.Error(), pbstatus.ResultStatus_DBException)
		return nil
	}

	out.Status = outLog(path, out)
	return nil
}

func (mine *PhotocopyService) UpdateByFilter(ctx context.Context, in *pb.RequestUpdate, out *pb.ReplyInfo) error {
	path := "photocopy.updateByFilter"
	inLog(path, in)
	if len(in.Uid) < 1 {
		out.Status = outError(path, "the uid is empty ", pbstatus.ResultStatus_Empty)
		return nil
	}
	_, er := cache.Context().GetPhotocopy(in.Uid)
	if er != nil {
		out.Status = outError(path, "the photocopy not found ", pbstatus.ResultStatus_NotExisted)
		return nil
	}
	var err error

	if err != nil {
		out.Status = outError(path, err.Error(), pbstatus.ResultStatus_DBException)
		return nil
	}
	out.Status = outLog(path, out)
	return nil
}

func (mine *PhotocopyService) AppendSlot(ctx context.Context, in *pb.RequestList, out *pb.ReplyList) error {
	path := "photocopy.appendSlot"
	inLog(path, in)
	if len(in.Uid) < 1 {
		out.Status = outError(path, "the uid is empty ", pbstatus.ResultStatus_Empty)
		return nil
	}
	info, er := cache.Context().GetPhotocopy(in.Uid)
	if er != nil {
		out.Status = outError(path, er.Error(), pbstatus.ResultStatus_NotExisted)
		return nil
	}

	out.Uid = info.UID
	out.Status = outLog(path, out)
	return nil
}

func (mine *PhotocopyService) SubtractSlot(ctx context.Context, in *pb.RequestList, out *pb.ReplyList) error {
	path := "photocopy.subtractSlot"
	inLog(path, in)
	if len(in.Uid) < 1 {
		out.Status = outError(path, "the uid is empty ", pbstatus.ResultStatus_Empty)
		return nil
	}
	info, er := cache.Context().GetPhotocopy(in.Uid)
	if er != nil {
		out.Status = outError(path, er.Error(), pbstatus.ResultStatus_NotExisted)
		return nil
	}

	out.Uid = info.UID
	out.Status = outLog(path, out)
	return nil
}
