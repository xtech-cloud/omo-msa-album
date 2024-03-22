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

type StyleService struct{}

func switchStyle(info *cache.CertificateStyleInfo) *pb.StyleInfo {
	tmp := new(pb.StyleInfo)
	tmp.Uid = info.UID
	tmp.Id = info.ID
	tmp.Created = info.Created
	tmp.Updated = info.Updated
	tmp.Operator = info.Operator
	tmp.Creator = info.Creator
	tmp.Name = info.Name
	tmp.Remark = info.Remark
	tmp.Cover = info.Cover
	tmp.Prefix = info.Prefix
	tmp.Background = info.Background
	tmp.Type = uint32(info.Type)
	tmp.Tags = info.Tags
	tmp.Scenes = info.Scenes
	tmp.Width = uint32(info.Width)
	tmp.Height = uint32(info.Height)
	tmp.Slots = switchStyleSlots(info.Slots)
	tmp.Relates = switchStyleRelates(info.Relates)
	return tmp
}

func (mine *StyleService) AddOne(ctx context.Context, in *pb.ReqStyleAdd, out *pb.ReplyStyleInfo) error {
	path := "style.add"
	inLog(path, in)
	if len(in.Name) < 1 {
		out.Status = outError(path, "the name is empty ", pbstatus.ResultStatus_Empty)
		return nil
	}
	slots := make([]proxy.StyleSlot, 0, len(in.Slots))
	for _, slot := range in.Slots {
		slots = append(slots, proxy.StyleSlot{
			Key:    slot.Key,
			X:      slot.X,
			Y:      slot.Y,
			Width:  slot.Width,
			Height: slot.Height,
			Bold:   slot.Bold,
			Size:   slot.Size,
		})
	}
	info, err := cache.Context().CreateStyle(in, slots)
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
	info, er := cache.Context().GetStyle(in.Uid)
	if er != nil {
		out.Status = outError(path, er.Error(), pbstatus.ResultStatus_NotExisted)
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
	info, er := cache.Context().GetStyle(in.Uid)
	if er != nil {
		out.Status = outError(path, er.Error(), pbstatus.ResultStatus_NotExisted)
		return nil
	}
	var err error
	err = info.UpdateBase(in.Remark, in)
	if err != nil {
		out.Status = outError(path, err.Error(), pbstatus.ResultStatus_DBException)
		return nil
	}

	out.Status = outLog(path, out)
	return nil
}

func (mine *StyleService) GetStatistic(ctx context.Context, in *pb.RequestFilter, out *pb.ReplyStatistic) error {
	path := "style.getStatistic"
	inLog(path, in)
	if len(in.Field) < 1 {
		out.Status = outError(path, "the field is empty ", pbstatus.ResultStatus_Empty)
		return nil
	}

	if in.Field == "sn" {
		info, err := cache.Context().GetStyle(in.Value)
		if err != nil {
			out.Status = outError(path, err.Error(), pbstatus.ResultStatus_NotExisted)
			return nil
		}
		out.Key = info.GetSN(in.Operator)
	} else if in.Field == "count" {
		out.Count = cache.Context().GetCertificatesCountByStyle(in.Value)
	} else if in.Field == "scene_count" {
		out.Count = cache.Context().GetCertificatesCountBySceneStyle(in.Owner, in.Value)
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
	info, er := cache.Context().GetStyle(in.Uid)
	if er != nil {
		out.Status = outError(path, er.Error(), pbstatus.ResultStatus_NotExisted)
		return nil
	}
	er = info.Remove(in.Operator)
	if er != nil {
		out.Status = outError(path, er.Error(), pbstatus.ResultStatus_DBException)
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
	var list []*cache.CertificateStyleInfo
	var err error

	if in.Field == "" {
		out.Total, out.Pages, list = cache.Context().GetStyles(in.Page, in.Number)
	} else if in.Field == "scene" {
		list = cache.Context().GetStylesByScene(in.Value)
	} else if in.Field == "array" {
		list = cache.Context().GetStylesByArray(in.Values)
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
	info, err := cache.Context().GetStyle(in.Uid)
	if err != nil {
		out.Status = outError(path, err.Error(), pbstatus.ResultStatus_NotExisted)
		return nil
	}
	if in.Field == "cover" {
		err = info.UpdateCover(in.Value, in.Operator)
	} else if in.Field == "relate" {
		tp := 0
		if len(in.Values) > 0 {
			tp = parseStringToInt(in.Values[0])
		}
		err = info.AppendRelate(in.Value, in.Operator, uint32(tp))
	}
	if err != nil {
		out.Status = outError(path, err.Error(), pbstatus.ResultStatus_DBException)
		return nil
	}
	out.Status = outLog(path, out)
	return nil
}
