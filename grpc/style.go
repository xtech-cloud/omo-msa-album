package grpc

import (
	"context"
	"errors"
	"fmt"
	pb "github.com/xtech-cloud/omo-msp-album/proto/album"
	pbstatus "github.com/xtech-cloud/omo-msp-status/proto/status"
	"omo.msa.album/cache"
)

type StyleService struct{}

func switchStyle(info *cache.PhotoStyleInfo) *pb.StyleInfo {
	tmp := new(pb.StyleInfo)
	tmp.Uid = info.UID
	tmp.Id = info.ID
	tmp.Created = info.CreateTime.Unix()
	tmp.Updated = info.UpdateTime.Unix()
	tmp.Operator = info.Operator
	tmp.Creator = info.Creator
	tmp.Name = info.Name
	tmp.Remark = info.Remark
	tmp.Cover = info.Cover
	tmp.Size = uint32(info.Size)
	tmp.Price = uint32(info.Price)
	tmp.Type = uint32(info.Type)
	tmp.Tags = info.Tags
	tmp.Slots = switchSlots(info.Slots)
	return tmp
}

func (mine *StyleService) AddOne(ctx context.Context, in *pb.ReqStyleAdd, out *pb.ReplyStyleInfo) error {
	path := "style.add"
	inLog(path, in)
	if len(in.Name) < 1 {
		out.Status = outError(path, "the name is empty ", pbstatus.ResultStatus_Empty)
		return nil
	}

	info, err := cache.Context().CreatePhotoTemplate(in.Name, in.Remark, in.Operator, uint8(in.Type))
	if err != nil {
		out.Status = outError(path, err.Error(), pbstatus.ResultStatus_DBException)
		return nil
	}
	out.Info = switchStyle(info)
	out.Status = outLog(path, out)
	return nil
}

func (mine *StyleService) GetOne(ctx context.Context, in *pb.RequestInfo, out *pb.ReplyStyleInfo) error {
	path := "style.getOne"
	inLog(path, in)
	if len(in.Uid) < 1 {
		out.Status = outError(path, "the uid is empty ", pbstatus.ResultStatus_Empty)
		return nil
	}
	info, er := cache.Context().GetPhotoTemplate(in.Uid)
	if er != nil {
		out.Status = outError(path, "the style not found ", pbstatus.ResultStatus_NotExisted)
		return nil
	}
	out.Info = switchStyle(info)
	out.Status = outLog(path, out)
	return nil
}

func (mine *StyleService) UpdateBase(ctx context.Context, in *pb.ReqStyleUpdate, out *pb.ReplyInfo) error {
	path := "style.updateBase"
	inLog(path, in)
	if len(in.Uid) < 1 {
		out.Status = outError(path, "the uid is empty ", pbstatus.ResultStatus_Empty)
		return nil
	}
	info, er := cache.Context().GetPhotoTemplate(in.Uid)
	if er != nil {
		out.Status = outError(path, er.Error(), pbstatus.ResultStatus_NotExisted)
		return nil
	}
	var err error
	err = info.UpdateBase(in.Name, in.Remark, in.Operator, uint8(in.Type))
	if err != nil {
		out.Status = outError(path, err.Error(), pbstatus.ResultStatus_DBException)
		return nil
	}

	out.Status = outLog(path, out)
	return nil
}

func (mine *StyleService) GetStatistic(ctx context.Context, in *pb.RequestFilter, out *pb.ReplyStatistic) error {
	path := "Style.getStatistic"
	inLog(path, in)
	if len(in.Field) < 1 {
		out.Status = outError(path, "the user is empty ", pbstatus.ResultStatus_Empty)
		return nil
	}

	out.Status = outLog(path, out)
	return nil
}

func (mine *StyleService) RemoveOne(ctx context.Context, in *pb.RequestInfo, out *pb.ReplyInfo) error {
	path := "style.remove"
	inLog(path, in)
	if len(in.Uid) < 1 {
		out.Status = outError(path, "the uid is empty ", pbstatus.ResultStatus_Empty)
		return nil
	}
	info, er := cache.Context().GetPhotoTemplate(in.Uid)
	if er != nil {
		out.Status = outError(path, "the photo style not found ", pbstatus.ResultStatus_NotExisted)
		return nil
	}
	er = info.Remove(in.Operator)
	if er != nil {
		out.Status = outError(path, "the photo style not found ", pbstatus.ResultStatus_NotExisted)
		return nil
	}

	out.Uid = in.Uid
	out.Status = outLog(path, out)
	return nil
}

func (mine *StyleService) Search(ctx context.Context, in *pb.RequestInfo, out *pb.ReplyStyleList) error {
	path := "style.search"
	inLog(path, in)

	out.Status = outLog(path, fmt.Sprintf("the length = %d", len(out.List)))
	return nil
}

func (mine *StyleService) GetListByFilter(ctx context.Context, in *pb.RequestFilter, out *pb.ReplyStyleList) error {
	path := "style.getListByFilter"
	inLog(path, in)
	var list []*cache.PhotoStyleInfo
	var err error
	if in.Field == "" {
		//list,err = cache.Context().GetPhotoStyles(in.Page, in.Number)
	} else {
		err = errors.New("the key not defined")
	}
	if err != nil {
		out.Status = outError(path, err.Error(), pbstatus.ResultStatus_DBException)
		return nil
	}
	out.List = make([]*pb.StyleInfo, 0, len(list))
	for _, value := range list {
		out.List = append(out.List, switchStyle(value))
	}

	out.Status = outLog(path, fmt.Sprintf("the length = %d", len(out.List)))
	return nil
}

func (mine *StyleService) UpdateByFilter(ctx context.Context, in *pb.RequestUpdate, out *pb.ReplyInfo) error {
	path := "style.updateByFilter"
	inLog(path, in)
	if len(in.Uid) < 1 {
		out.Status = outError(path, "the uid is empty ", pbstatus.ResultStatus_Empty)
		return nil
	}

	out.Status = outLog(path, out)
	return nil
}

func (mine *StyleService) AppendSlot(ctx context.Context, in *pb.RequestList, out *pb.ReplyList) error {
	path := "style.appendSlot"
	inLog(path, in)
	if len(in.Uid) < 1 {
		out.Status = outError(path, "the uid is empty ", pbstatus.ResultStatus_Empty)
		return nil
	}

	out.Uid = in.Uid
	out.Status = outLog(path, out)
	return nil
}

func (mine *StyleService) SubtractSlot(ctx context.Context, in *pb.RequestList, out *pb.ReplyList) error {
	path := "style.subtractSlot"
	inLog(path, in)
	if len(in.Uid) < 1 {
		out.Status = outError(path, "the uid is empty ", pbstatus.ResultStatus_Empty)
		return nil
	}

	out.Uid = in.Uid
	out.Status = outLog(path, out)
	return nil
}
