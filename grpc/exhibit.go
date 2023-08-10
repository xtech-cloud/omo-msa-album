package grpc

import (
	"context"
	"errors"
	"fmt"
	pb "github.com/xtech-cloud/omo-msp-album/proto/album"
	pbstatus "github.com/xtech-cloud/omo-msp-status/proto/status"
	"omo.msa.album/cache"
)

type ExhibitService struct{}

func switchExhibit(info *cache.ExhibitInfo) *pb.ExhibitInfo {
	tmp := new(pb.ExhibitInfo)
	tmp.Uid = info.UID
	tmp.Id = info.ID
	tmp.Created = info.Created
	tmp.Updated = info.Updated
	tmp.Operator = info.Operator
	tmp.Creator = info.Creator
	tmp.Name = info.Name
	tmp.Remark = info.Remark
	tmp.Cover = info.Cover
	tmp.Status = uint32(info.Status)
	tmp.Type = uint32(info.Type)
	tmp.Owner = info.Owner
	tmp.Assets = info.Assets
	tmp.Tags = info.Tags
	return tmp
}

func (mine *ExhibitService) AddOne(ctx context.Context, in *pb.ReqExhibitAdd, out *pb.ReplyExhibitInfo) error {
	path := "exhibit.add"
	inLog(path, in)
	if len(in.Name) < 1 {
		out.Status = outError(path, "the name is empty ", pbstatus.ResultStatus_Empty)
		return nil
	}

	if cache.Context().HadExhibitByName(in.Name) {
		out.Status = outError(path, "the name is repeated", pbstatus.ResultStatus_Repeated)
		return nil
	}

	info, err := cache.Context().CreateExhibit(in.Name, in.Remark, in.Cover, in.Owner, in.Operator)
	if err != nil {
		out.Status = outError(path, err.Error(), pbstatus.ResultStatus_DBException)
		return nil
	}
	out.Info = switchExhibit(info)
	out.Status = outLog(path, out)
	return nil
}

func (mine *ExhibitService) GetOne(ctx context.Context, in *pb.RequestInfo, out *pb.ReplyExhibitInfo) error {
	path := "exhibit.getOne"
	inLog(path, in)
	if len(in.Uid) < 1 {
		out.Status = outError(path, "the uid is empty ", pbstatus.ResultStatus_Empty)
		return nil
	}
	info, er := cache.Context().GetExhibit(in.Uid)
	if er != nil {
		out.Status = outError(path, "the exhibit not found ", pbstatus.ResultStatus_NotExisted)
		return nil
	}
	out.Info = switchExhibit(info)
	out.Status = outLog(path, out)
	return nil
}

func (mine *ExhibitService) GetStatistic(ctx context.Context, in *pb.RequestFilter, out *pb.ReplyStatistic) error {
	path := "exhibit.getStatistic"
	inLog(path, in)
	if len(in.Field) < 1 {
		out.Status = outError(path, "the user is empty ", pbstatus.ResultStatus_Empty)
		return nil
	}

	out.Status = outLog(path, out)
	return nil
}

func (mine *ExhibitService) RemoveOne(ctx context.Context, in *pb.RequestInfo, out *pb.ReplyInfo) error {
	path := "exhibit.remove"
	inLog(path, in)
	if len(in.Uid) < 1 {
		out.Status = outError(path, "the uid is empty ", pbstatus.ResultStatus_Empty)
		return nil
	}
	er := cache.Context().RemoveExhibit(in.Uid, in.Operator)
	if er != nil {
		out.Status = outError(path, "the exhibit not found ", pbstatus.ResultStatus_NotExisted)
		return nil
	}

	out.Uid = in.Uid
	out.Status = outLog(path, out)
	return nil
}

func (mine *ExhibitService) Search(ctx context.Context, in *pb.RequestInfo, out *pb.ReplyExhibitList) error {
	path := "exhibit.search"
	inLog(path, in)

	out.Status = outLog(path, fmt.Sprintf("the length = %d", len(out.List)))
	return nil
}

func (mine *ExhibitService) GetListByFilter(ctx context.Context, in *pb.RequestFilter, out *pb.ReplyExhibitList) error {
	path := "exhibit.getListByFilter"
	inLog(path, in)
	var list []*cache.ExhibitInfo
	var err error
	if in.Field == "" {
		list, err = cache.Context().GetAllExhibits(in.Owner)
	} else if in.Field == "array" {
		list = cache.Context().GetExhibits(in.Values)
	} else {
		err = errors.New("the key not defined")
	}
	if err != nil {
		out.Status = outError(path, err.Error(), pbstatus.ResultStatus_DBException)
		return nil
	}
	out.List = make([]*pb.ExhibitInfo, 0, len(list))
	for _, value := range list {
		out.List = append(out.List, switchExhibit(value))
	}

	out.Status = outLog(path, fmt.Sprintf("the length = %d", len(out.List)))
	return nil
}

func (mine *ExhibitService) UpdateBase(ctx context.Context, in *pb.ReqExhibitUpdate, out *pb.ReplyInfo) error {
	path := "exhibit.updateBase"
	inLog(path, in)
	if len(in.Uid) < 1 {
		out.Status = outError(path, "the uid is empty ", pbstatus.ResultStatus_Empty)
		return nil
	}
	info, er := cache.Context().GetExhibit(in.Uid)
	if er != nil {
		out.Status = outError(path, er.Error(), pbstatus.ResultStatus_NotExisted)
		return nil
	}
	if in.Name != info.Name && cache.Context().HadExhibitByName(in.Name) {
		out.Status = outError(path, "the name is repeated", pbstatus.ResultStatus_Repeated)
		return nil
	}
	var err error
	err = info.UpdateBase(in.Name, in.Remark, in.Cover, in.Operator)
	if err != nil {
		out.Status = outError(path, err.Error(), pbstatus.ResultStatus_DBException)
		return nil
	}

	out.Status = outLog(path, out)
	return nil
}

func (mine *ExhibitService) UpdateByFilter(ctx context.Context, in *pb.RequestUpdate, out *pb.ReplyInfo) error {
	path := "exhibit.updateByFilter"
	inLog(path, in)
	if len(in.Uid) < 1 {
		out.Status = outError(path, "the uid is empty ", pbstatus.ResultStatus_Empty)
		return nil
	}
	info, er := cache.Context().GetExhibit(in.Uid)
	if er != nil {
		out.Status = outError(path, "the exhibit not found ", pbstatus.ResultStatus_NotExisted)
		return nil
	}
	var err error
	if in.Field == "assets" {
		err = info.UpdateAssets(in.Values, in.Operator)
	} else {
		err = errors.New("the key not defined")
	}
	if err != nil {
		out.Status = outError(path, err.Error(), pbstatus.ResultStatus_DBException)
		return nil
	}
	out.Status = outLog(path, out)
	return nil
}

func (mine *ExhibitService) UpdateStatus(ctx context.Context, in *pb.RequestIntFlag, out *pb.ReplyInfo) error {
	path := "exhibit.updateStatus"
	inLog(path, in)
	if len(in.Uid) < 1 {
		out.Status = outError(path, "the uid is empty ", pbstatus.ResultStatus_Empty)
		return nil
	}
	info, er := cache.Context().GetExhibit(in.Uid)
	if er != nil {
		out.Status = outError(path, er.Error(), pbstatus.ResultStatus_NotExisted)
		return nil
	}
	err := info.UpdateStatus(in.Operator, uint8(in.Flag))
	if err != nil {
		out.Status = outError(path, err.Error(), pbstatus.ResultStatus_DBException)
		return nil
	}
	out.Status = outLog(path, out)
	return nil
}

func (mine *ExhibitService) AppendAsset(ctx context.Context, in *pb.RequestList, out *pb.ReplyList) error {
	path := "exhibit.appendAsset"
	inLog(path, in)
	if len(in.Uid) < 1 {
		out.Status = outError(path, "the uid is empty ", pbstatus.ResultStatus_Empty)
		return nil
	}
	info, er := cache.Context().GetExhibit(in.Uid)
	if er != nil {
		out.Status = outError(path, er.Error(), pbstatus.ResultStatus_NotExisted)
		return nil
	}

	err := info.AppendAsset(in.List[0], in.Operator)
	if err != nil {
		out.Status = outError(path, err.Error(), pbstatus.ResultStatus_DBException)
		return nil
	}
	out.Uid = in.Uid
	out.List = info.Assets
	out.Status = outLog(path, out)
	return nil
}

func (mine *ExhibitService) SubtractAsset(ctx context.Context, in *pb.RequestList, out *pb.ReplyList) error {
	path := "exhibit.subtractAsset"
	inLog(path, in)
	if len(in.Uid) < 1 {
		out.Status = outError(path, "the uid is empty ", pbstatus.ResultStatus_Empty)
		return nil
	}
	info, er := cache.Context().GetExhibit(in.Uid)
	if er != nil {
		out.Status = outError(path, er.Error(), pbstatus.ResultStatus_NotExisted)
		return nil
	}

	err := info.SubtractAsset(in.List[0], in.Operator)
	if err != nil {
		out.Status = outError(path, err.Error(), pbstatus.ResultStatus_DBException)
		return nil
	}
	out.Uid = in.Uid
	out.List = info.Assets
	out.Status = outLog(path, out)
	return nil
}
