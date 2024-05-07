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

type SheetService struct{}

func switchSheet(info *cache.SheetInfo) *pb.SheetInfo {
	tmp := new(pb.SheetInfo)
	tmp.Uid = info.UID
	tmp.Id = info.ID
	tmp.Created = info.Created
	tmp.Updated = info.Updated
	tmp.Operator = info.Operator
	tmp.Creator = info.Creator
	tmp.Name = info.Name
	tmp.Remark = info.Remark
	tmp.Owner = info.Owner
	tmp.Cover = info.Cover
	tmp.Area = info.Target
	tmp.Aspect = info.Aspect

	tmp.Tags = info.Tags
	tmp.Pages = switchSheetPages(info.Pages)
	return tmp
}

func switchSheetPages(list []*proxy.SheetPage) []*pb.SheetPage {
	slots := make([]*pb.SheetPage, 0, len(list))
	for _, slot := range list {
		item := &pb.SheetPage{Weight: slot.Weight, Page: slot.UID}
		slots = append(slots, item)
	}
	return slots
}

func (mine *SheetService) AddOne(ctx context.Context, in *pb.ReqSheetAdd, out *pb.ReplySheetInfo) error {
	path := "sheet.addOne"
	inLog(path, in)
	if len(in.Name) < 1 {
		out.Status = outError(path, "the name is empty ", pbstatus.ResultStatus_Empty)
		return nil
	}

	info, err := cache.Context().CreateSheet(in.Name, in.Remark, in.Operator, in.Owner, in.Aspect, in.Tags)
	if err != nil {
		out.Status = outError(path, err.Error(), pbstatus.ResultStatus_DBException)
		return nil
	}
	out.Info = switchSheet(info)
	out.Status = outLog(path, out)
	return nil
}

func (mine *SheetService) GetOne(ctx context.Context, in *pb.RequestInfo, out *pb.ReplySheetInfo) error {
	path := "sheet.getOne"
	inLog(path, in)
	if len(in.Uid) < 1 {
		out.Status = outError(path, "the uid is empty ", pbstatus.ResultStatus_Empty)
		return nil
	}
	info, er := cache.Context().GetSheet(in.Uid)
	if er != nil {
		out.Status = outError(path, "the sheet not found ", pbstatus.ResultStatus_NotExisted)
		return nil
	}
	out.Info = switchSheet(info)
	out.Status = outLog(path, out)
	return nil
}

func (mine *SheetService) GetStatistic(ctx context.Context, in *pb.RequestFilter, out *pb.ReplyStatistic) error {
	path := "sheet.getStatistic"
	inLog(path, in)
	if len(in.Field) < 1 {
		out.Status = outError(path, "the user is empty ", pbstatus.ResultStatus_Empty)
		return nil
	}
	if in.Field == "count" {
		info, _ := cache.Context().GetSheet(in.Value)
		if info != nil {
			out.Count = info.GetAssetCount()
		}
	}

	out.Status = outLog(path, out)
	return nil
}

func (mine *SheetService) RemoveOne(ctx context.Context, in *pb.RequestInfo, out *pb.ReplyInfo) error {
	path := "sheet.remove"
	inLog(path, in)
	if len(in.Uid) < 1 {
		out.Status = outError(path, "the uid is empty ", pbstatus.ResultStatus_Empty)
		return nil
	}
	info, er := cache.Context().GetSheet(in.Uid)
	if er != nil {
		out.Status = outError(path, "the sheet not found ", pbstatus.ResultStatus_NotExisted)
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

func (mine *SheetService) Search(ctx context.Context, in *pb.RequestInfo, out *pb.ReplySheetList) error {
	path := "sheet.search"
	inLog(path, in)

	out.Status = outLog(path, fmt.Sprintf("the length = %d", len(out.List)))
	return nil
}

func (mine *SheetService) GetListBy(ctx context.Context, in *pb.RequestFilter, out *pb.ReplySheetList) error {
	path := "sheet.getListByFilter"
	inLog(path, in)
	var list []*cache.SheetInfo
	var err error
	if in.Field == "" {
		list = cache.Context().GetSheetsByOwner(in.Owner)
	} else if in.Field == "target" {
		list = cache.Context().GetSheetsByTarget(in.Value)
	} else {
		err = errors.New("the key not defined")
	}
	if err != nil {
		out.Status = outError(path, err.Error(), pbstatus.ResultStatus_DBException)
		return nil
	}
	out.List = make([]*pb.SheetInfo, 0, len(list))
	for _, value := range list {
		out.List = append(out.List, switchSheet(value))
	}

	out.Status = outLog(path, fmt.Sprintf("the length = %d", len(out.List)))
	return nil
}

func (mine *SheetService) UpdateBase(ctx context.Context, in *pb.ReqSheetUpdate, out *pb.ReplyInfo) error {
	path := "sheet.updateBase"
	inLog(path, in)
	if len(in.Uid) < 1 {
		out.Status = outError(path, "the uid is empty ", pbstatus.ResultStatus_Empty)
		return nil
	}
	info, er := cache.Context().GetSheet(in.Uid)
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

func (mine *SheetService) UpdateBy(ctx context.Context, in *pb.RequestUpdate, out *pb.ReplyInfo) error {
	path := "sheet.updateByFilter"
	inLog(path, in)
	if len(in.Uid) < 1 {
		out.Status = outError(path, "the uid is empty ", pbstatus.ResultStatus_Empty)
		return nil
	}
	info, er := cache.Context().GetSheet(in.Uid)
	if er != nil {
		out.Status = outError(path, "the sheet not found ", pbstatus.ResultStatus_NotExisted)
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

func (mine *SheetService) UpdatePages(ctx context.Context, in *pb.ReqSheetPages, out *pb.ReplySheetPages) error {
	path := "sheet.updatePages"
	inLog(path, in)
	info, er := cache.Context().GetSheet(in.Uid)
	if er != nil {
		out.Status = outError(path, er.Error(), pbstatus.ResultStatus_NotExisted)
		return nil
	}
	slots := make([]*proxy.SheetPage, 0, 3)
	for _, slot := range in.List {
		slots = append(slots, &proxy.SheetPage{
			UID: slot.Page, Weight: slot.Weight,
		})
	}
	err := info.UpdatePages(slots, in.Operator)
	if err != nil {
		out.Status = outError(path, err.Error(), pbstatus.ResultStatus_DBException)
		return nil
	}
	out.List = switchSheetPages(info.Pages)
	out.Status = outLog(path, out)
	return nil
}
