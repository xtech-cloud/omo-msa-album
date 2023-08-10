package grpc

import (
	"context"
	"errors"
	"fmt"
	pb "github.com/xtech-cloud/omo-msp-album/proto/album"
	pbstatus "github.com/xtech-cloud/omo-msp-status/proto/status"
	"omo.msa.album/cache"
)

type FrameService struct{}

func switchFrame(info *cache.PhotoFrameInfo) *pb.FrameInfo {
	tmp := new(pb.FrameInfo)
	tmp.Uid = info.UID
	tmp.Id = info.ID
	tmp.Created = info.Created
	tmp.Updated = info.Updated
	tmp.Operator = info.Operator
	tmp.Creator = info.Creator
	tmp.Name = info.Name
	tmp.Remark = info.Remark
	tmp.Asset = info.Asset
	tmp.Width = info.Width
	tmp.Height = info.Height
	tmp.Owner = info.Owner
	return tmp
}

func (mine *FrameService) AddOne(ctx context.Context, in *pb.ReqFrameAdd, out *pb.ReplyFrameInfo) error {
	path := "frame.add"
	inLog(path, in)
	if len(in.Name) < 1 {
		out.Status = outError(path, "the name is empty ", pbstatus.ResultStatus_Empty)
		return nil
	}

	info, err := cache.Context().CreatePhotoFrame(in.Name, in.Remark, in.Asset, in.Owner, in.Operator, in.Width, in.Height)
	if err != nil {
		out.Status = outError(path, err.Error(), pbstatus.ResultStatus_DBException)
		return nil
	}
	out.Info = switchFrame(info)
	out.Status = outLog(path, out)
	return nil
}

func (mine *FrameService) GetOne(ctx context.Context, in *pb.RequestInfo, out *pb.ReplyFrameInfo) error {
	path := "frame.getOne"
	inLog(path, in)
	if len(in.Uid) < 1 {
		out.Status = outError(path, "the uid is empty ", pbstatus.ResultStatus_Empty)
		return nil
	}
	info, er := cache.Context().GetPhotoFrame(in.Uid)
	if er != nil {
		out.Status = outError(path, "the frame not found ", pbstatus.ResultStatus_NotExisted)
		return nil
	}
	out.Info = switchFrame(info)
	out.Status = outLog(path, out)
	return nil
}

func (mine *FrameService) UpdateBase(ctx context.Context, in *pb.ReqFrameUpdate, out *pb.ReplyInfo) error {
	path := "frame.updateBase"
	inLog(path, in)
	if len(in.Uid) < 1 {
		out.Status = outError(path, "the uid is empty ", pbstatus.ResultStatus_Empty)
		return nil
	}
	info, er := cache.Context().GetPhotoFrame(in.Uid)
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

func (mine *FrameService) GetStatistic(ctx context.Context, in *pb.RequestFilter, out *pb.ReplyStatistic) error {
	path := "frame.getStatistic"
	inLog(path, in)
	if len(in.Field) < 1 {
		out.Status = outError(path, "the user is empty ", pbstatus.ResultStatus_Empty)
		return nil
	}

	out.Status = outLog(path, out)
	return nil
}

func (mine *FrameService) RemoveOne(ctx context.Context, in *pb.RequestInfo, out *pb.ReplyInfo) error {
	path := "frame.remove"
	inLog(path, in)
	if len(in.Uid) < 1 {
		out.Status = outError(path, "the uid is empty ", pbstatus.ResultStatus_Empty)
		return nil
	}
	info, er := cache.Context().GetPhotoFrame(in.Uid)
	if er != nil {
		out.Status = outError(path, "the photo frame not found ", pbstatus.ResultStatus_NotExisted)
		return nil
	}
	er = info.Remove(in.Operator)
	if er != nil {
		out.Status = outError(path, "the photo frame not found ", pbstatus.ResultStatus_NotExisted)
		return nil
	}

	out.Uid = in.Uid
	out.Status = outLog(path, out)
	return nil
}

func (mine *FrameService) Search(ctx context.Context, in *pb.RequestInfo, out *pb.ReplyFrameList) error {
	path := "frame.search"
	inLog(path, in)

	out.Status = outLog(path, fmt.Sprintf("the length = %d", len(out.List)))
	return nil
}

func (mine *FrameService) GetListByFilter(ctx context.Context, in *pb.RequestFilter, out *pb.ReplyFrameList) error {
	path := "frame.getListByFilter"
	inLog(path, in)
	var list []*cache.PhotoFrameInfo
	var err error
	if in.Field == "" {
		list = cache.Context().GetPhotoFrames(in.Owner)
	} else {
		err = errors.New("the key not defined")
	}
	if err != nil {
		out.Status = outError(path, err.Error(), pbstatus.ResultStatus_DBException)
		return nil
	}
	out.List = make([]*pb.FrameInfo, 0, len(list))
	for _, value := range list {
		out.List = append(out.List, switchFrame(value))
	}

	out.Status = outLog(path, fmt.Sprintf("the length = %d", len(out.List)))
	return nil
}

func (mine *FrameService) UpdateByFilter(ctx context.Context, in *pb.RequestUpdate, out *pb.ReplyInfo) error {
	path := "frame.updateByFilter"
	inLog(path, in)
	if len(in.Uid) < 1 {
		out.Status = outError(path, "the uid is empty ", pbstatus.ResultStatus_Empty)
		return nil
	}

	out.Status = outLog(path, out)
	return nil
}
