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

type CompositionService struct{}

func switchComposition(info *cache.CompositionInfo) *pb.CompositionInfo {
	tmp := new(pb.CompositionInfo)
	tmp.Uid = info.UID
	tmp.Id = info.ID
	tmp.Created = info.Created
	tmp.Updated = info.Updated
	tmp.Operator = info.Operator
	tmp.Creator = info.Creator
	tmp.Name = info.Name
	tmp.Remark = info.Remark
	tmp.Cover = info.Cover
	tmp.Aspect = uint32(info.Aspect)
	tmp.Owner = info.Owner
	tmp.Tags = info.Tags
	tmp.Slots = switchCompositionSlots(info.Slots)
	return tmp
}

func switchCompositionSlots(list []*proxy.SlotInfo) []*pb.SlotInfo {
	slots := make([]*pb.SlotInfo, 0, len(list))
	for _, slot := range list {
		item := &pb.SlotInfo{Index: slot.Index, X: slot.Position.X, Y: slot.Position.Y,
			Type: uint32(slot.Type), Width: uint32(slot.Size.X), Height: uint32(slot.Size.Y)}
		slots = append(slots, item)
	}
	return slots
}

func (mine *CompositionService) AddOne(ctx context.Context, in *pb.ReqCompositionAdd, out *pb.ReplyCompositionInfo) error {
	path := "composition.addOne"
	inLog(path, in)
	if len(in.Name) < 1 {
		out.Status = outError(path, "the name is empty ", pbstatus.ResultStatus_Empty)
		return nil
	}

	info, err := cache.Context().CreateComposition(in.Name, in.Remark, in.Operator, in.Owner, in.Cover, cache.AspectType(in.Aspect), in.Tags)
	if err != nil {
		out.Status = outError(path, err.Error(), pbstatus.ResultStatus_DBException)
		return nil
	}
	out.Info = switchComposition(info)
	out.Status = outLog(path, out)
	return nil
}

func (mine *CompositionService) GetOne(ctx context.Context, in *pb.RequestInfo, out *pb.ReplyCompositionInfo) error {
	path := "composition.getOne"
	inLog(path, in)
	if len(in.Uid) < 1 {
		out.Status = outError(path, "the uid is empty ", pbstatus.ResultStatus_Empty)
		return nil
	}
	info, er := cache.Context().GetComposition(in.Uid)
	if er != nil {
		out.Status = outError(path, "the composition not found ", pbstatus.ResultStatus_NotExisted)
		return nil
	}
	out.Info = switchComposition(info)
	out.Status = outLog(path, out)
	return nil
}

func (mine *CompositionService) GetStatistic(ctx context.Context, in *pb.RequestFilter, out *pb.ReplyStatistic) error {
	path := "composition.getStatistic"
	inLog(path, in)
	if len(in.Field) < 1 {
		out.Status = outError(path, "the user is empty ", pbstatus.ResultStatus_Empty)
		return nil
	}

	out.Status = outLog(path, out)
	return nil
}

func (mine *CompositionService) RemoveOne(ctx context.Context, in *pb.RequestInfo, out *pb.ReplyInfo) error {
	path := "composition.remove"
	inLog(path, in)
	if len(in.Uid) < 1 {
		out.Status = outError(path, "the uid is empty ", pbstatus.ResultStatus_Empty)
		return nil
	}
	info, er := cache.Context().GetComposition(in.Uid)
	if er != nil {
		out.Status = outError(path, "the composition not found ", pbstatus.ResultStatus_NotExisted)
		return nil
	}
	err := info.Remove(in.Operator)
	if err != nil {
		out.Status = outError(path, err.Error(), pbstatus.ResultStatus_DBException)
		return nil
	}
	out.Uid = in.Uid
	out.Status = outLog(path, out)
	return nil
}

func (mine *CompositionService) Search(ctx context.Context, in *pb.RequestInfo, out *pb.ReplyCompositionList) error {
	path := "composition.search"
	inLog(path, in)

	out.Status = outLog(path, fmt.Sprintf("the length = %d", len(out.List)))
	return nil
}

func (mine *CompositionService) GetListBy(ctx context.Context, in *pb.RequestFilter, out *pb.ReplyCompositionList) error {
	path := "composition.getListByFilter"
	inLog(path, in)
	var list []*cache.CompositionInfo
	var err error
	if in.Field == "" {
		list = cache.Context().GetCompositionsByOwner(in.Owner)
	} else {
		err = errors.New("the key not defined")
	}
	if err != nil {
		out.Status = outError(path, err.Error(), pbstatus.ResultStatus_DBException)
		return nil
	}
	out.List = make([]*pb.CompositionInfo, 0, len(list))
	for _, value := range list {
		out.List = append(out.List, switchComposition(value))
	}

	out.Status = outLog(path, fmt.Sprintf("the length = %d", len(out.List)))
	return nil
}

func (mine *CompositionService) UpdateBase(ctx context.Context, in *pb.ReqCompositionUpdate, out *pb.ReplyInfo) error {
	path := "composition.updateBase"
	inLog(path, in)
	if len(in.Uid) < 1 {
		out.Status = outError(path, "the uid is empty ", pbstatus.ResultStatus_Empty)
		return nil
	}
	info, er := cache.Context().GetComposition(in.Uid)
	if er != nil {
		out.Status = outError(path, er.Error(), pbstatus.ResultStatus_NotExisted)
		return nil
	}
	var err error
	err = info.UpdateBase(in.Name, in.Remark, in.Operator, cache.AspectType(in.Aspect))
	if err != nil {
		out.Status = outError(path, err.Error(), pbstatus.ResultStatus_DBException)
		return nil
	}

	out.Status = outLog(path, out)
	return nil
}

func (mine *CompositionService) UpdateBy(ctx context.Context, in *pb.RequestUpdate, out *pb.ReplyInfo) error {
	path := "composition.updateByFilter"
	inLog(path, in)
	if len(in.Uid) < 1 {
		out.Status = outError(path, "the uid is empty ", pbstatus.ResultStatus_Empty)
		return nil
	}
	info, er := cache.Context().GetComposition(in.Uid)
	if er != nil {
		out.Status = outError(path, "the composition not found ", pbstatus.ResultStatus_NotExisted)
		return nil
	}
	var err error
	if in.Field == "cover" {
		err = info.UpdateCover(in.Value, in.Operator)
	} else {
		err = errors.New("the field not defined")
	}
	if err != nil {
		out.Status = outError(path, err.Error(), pbstatus.ResultStatus_DBException)
		return nil
	}
	out.Status = outLog(path, out)
	return nil
}

func (mine *CompositionService) SetSlots(ctx context.Context, in *pb.ReqCompositionSlots, out *pb.ReplyCompositionSlots) error {
	path := "composition.setSlots"
	inLog(path, in)
	info, er := cache.Context().GetComposition(in.Uid)
	if er != nil {
		out.Status = outError(path, er.Error(), pbstatus.ResultStatus_NotExisted)
		return nil
	}
	slots := make([]*proxy.SlotInfo, 0, 3)
	for _, slot := range in.Slots {
		slots = append(slots, &proxy.SlotInfo{
			Index: slot.Index, Type: uint8(slot.Type),
			Position: proxy.Vector2{X: int32(slot.X), Y: int32(slot.Y)},
			Size:     proxy.Vector2{X: int32(slot.Width), Y: int32(slot.Height)},
		})
	}
	err := info.UpdateSlots(slots, in.Operator)
	if err != nil {
		out.Status = outError(path, err.Error(), pbstatus.ResultStatus_DBException)
		return nil
	}
	out.List = switchCompositionSlots(info.Slots)
	out.Status = outLog(path, out)
	return nil
}
