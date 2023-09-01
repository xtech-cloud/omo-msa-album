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

type PageService struct{}

func switchPage(info *cache.PageInfo) *pb.PageInfo {
	tmp := new(pb.PageInfo)
	tmp.Uid = info.UID
	tmp.Id = info.ID
	tmp.Created = info.Created
	tmp.Updated = info.Updated
	tmp.Operator = info.Operator
	tmp.Creator = info.Creator
	tmp.Name = info.Name
	tmp.Remark = info.Remark
	tmp.Composition = info.Composition
	tmp.Type = uint32(info.Type)
	tmp.Lifecycle = info.Lifecycle
	tmp.Owner = info.Owner
	tmp.Tags = info.Tags
	tmp.Contents = switchPageContents(info.Contents)
	return tmp
}

func switchPageContents(list []*proxy.PageContents) []*pb.PageContent {
	slots := make([]*pb.PageContent, 0, len(list))
	for _, slot := range list {
		item := &pb.PageContent{Type: slot.Type, Slot: slot.Slot, Way: uint32(slot.Way),
			Interval: slot.Interval, List: slot.List}
		slots = append(slots, item)
	}
	return slots
}

func (mine *PageService) AddOne(ctx context.Context, in *pb.ReqPageAdd, out *pb.ReplyPageInfo) error {
	path := "page.addOne"
	inLog(path, in)
	info, err := cache.Context().CreatePage(in.Name, in.Remark, in.Operator, in.Owner, in.Composition, in.Type, in.Lifecycle, in.Tags)
	if err != nil {
		out.Status = outError(path, err.Error(), pbstatus.ResultStatus_DBException)
		return nil
	}
	out.Info = switchPage(info)
	out.Status = outLog(path, out)
	return nil
}

func (mine *PageService) GetOne(ctx context.Context, in *pb.RequestInfo, out *pb.ReplyPageInfo) error {
	path := "page.getOne"
	inLog(path, in)
	if len(in.Uid) < 1 {
		out.Status = outError(path, "the uid is empty ", pbstatus.ResultStatus_Empty)
		return nil
	}
	info, er := cache.Context().GetPage(in.Uid)
	if er != nil {
		out.Status = outError(path, "the page not found ", pbstatus.ResultStatus_NotExisted)
		return nil
	}
	out.Info = switchPage(info)
	out.Status = outLog(path, out)
	return nil
}

func (mine *PageService) GetStatistic(ctx context.Context, in *pb.RequestFilter, out *pb.ReplyStatistic) error {
	path := "page.getStatistic"
	inLog(path, in)
	if len(in.Field) < 1 {
		out.Status = outError(path, "the user is empty ", pbstatus.ResultStatus_Empty)
		return nil
	}

	out.Status = outLog(path, out)
	return nil
}

func (mine *PageService) RemoveOne(ctx context.Context, in *pb.RequestInfo, out *pb.ReplyInfo) error {
	path := "page.remove"
	inLog(path, in)
	if len(in.Uid) < 1 {
		out.Status = outError(path, "the uid is empty ", pbstatus.ResultStatus_Empty)
		return nil
	}
	info, er := cache.Context().GetPage(in.Uid)
	if er != nil {
		out.Status = outError(path, "the page not found ", pbstatus.ResultStatus_NotExisted)
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

func (mine *PageService) Search(ctx context.Context, in *pb.RequestInfo, out *pb.ReplyPageList) error {
	path := "page.search"
	inLog(path, in)

	out.Status = outLog(path, fmt.Sprintf("the length = %d", len(out.List)))
	return nil
}

func (mine *PageService) GetListBy(ctx context.Context, in *pb.RequestFilter, out *pb.ReplyPageList) error {
	path := "page.getListByFilter"
	inLog(path, in)
	var list []*cache.PageInfo
	var err error
	if in.Field == "" {
		list = cache.Context().GetPagesByOwner(in.Owner)
	} else if in.Field == "sheet" {
		list = cache.Context().GetPagesBySheet(in.Value)
	} else if in.Field == "type" {
		tp := parseStringToInt(in.Value)
		list = cache.Context().GetPagesByType(in.Owner, uint32(tp))
	} else if in.Field == "status" {
		tp := parseStringToInt(in.Value)
		list = cache.Context().GetPagesByStatus(in.Owner, uint32(tp))
	} else if in.Field == "list" {
		list = cache.Context().GetPagesByList(in.Values)
	} else {
		err = errors.New("the key not defined")
	}
	if err != nil {
		out.Status = outError(path, err.Error(), pbstatus.ResultStatus_DBException)
		return nil
	}
	out.List = make([]*pb.PageInfo, 0, len(list))
	for _, value := range list {
		out.List = append(out.List, switchPage(value))
	}

	out.Status = outLog(path, fmt.Sprintf("the length = %d", len(out.List)))
	return nil
}

func (mine *PageService) UpdateBase(ctx context.Context, in *pb.ReqPageUpdate, out *pb.ReplyInfo) error {
	path := "page.updateBase"
	inLog(path, in)
	if len(in.Uid) < 1 {
		out.Status = outError(path, "the uid is empty ", pbstatus.ResultStatus_Empty)
		return nil
	}
	info, er := cache.Context().GetPage(in.Uid)
	if er != nil {
		out.Status = outError(path, er.Error(), pbstatus.ResultStatus_NotExisted)
		return nil
	}
	var err error
	err = info.UpdateBase(in.Name, in.Remark, in.Operator, in.Lifecycle)
	if err != nil {
		out.Status = outError(path, err.Error(), pbstatus.ResultStatus_DBException)
		return nil
	}

	out.Status = outLog(path, out)
	return nil
}

func (mine *PageService) UpdateBy(ctx context.Context, in *pb.RequestUpdate, out *pb.ReplyInfo) error {
	path := "page.updateByFilter"
	inLog(path, in)
	if len(in.Uid) < 1 {
		out.Status = outError(path, "the uid is empty ", pbstatus.ResultStatus_Empty)
		return nil
	}
	//info, er := cache.Context().GetPage(in.Uid)
	//if er != nil {
	//	out.Status = outError(path, "the page not found ", pbstatus.ResultStatus_NotExisted)
	//	return nil
	//}
	//var err error
	//if in.Field == "cover" {
	//	info.up
	//} else {
	//	err = errors.New("the field not defined")
	//}
	//if err != nil {
	//	out.Status = outError(path, err.Error(), pbstatus.ResultStatus_DBException)
	//	return nil
	//}
	out.Status = outLog(path, out)
	return nil
}

func (mine *PageService) SetContent(ctx context.Context, in *pb.ReqPageContent, out *pb.ReplyPageContents) error {
	path := "page.setContent"
	inLog(path, in)
	info, er := cache.Context().GetPage(in.Uid)
	if er != nil {
		out.Status = outError(path, er.Error(), pbstatus.ResultStatus_NotExisted)
		return nil
	}
	tmp := &proxy.PageContents{Slot: in.Slot, Way: uint8(in.Fill), Interval: in.Interval, Type: in.Type, List: in.Assets}
	err := info.SetContent(tmp, in.Operator)
	if err != nil {
		out.Status = outError(path, err.Error(), pbstatus.ResultStatus_DBException)
		return nil
	}
	out.List = switchPageContents(info.Contents)
	out.Status = outLog(path, out)
	return nil
}
