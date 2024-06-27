package grpc

import (
	"context"
	"errors"
	"fmt"
	pb "github.com/xtech-cloud/omo-msp-album/proto/album"
	pbstatus "github.com/xtech-cloud/omo-msp-status/proto/status"
	"omo.msa.album/cache"
)

type CollectiveService struct{}

func switchCollective(info *cache.CollAlbumInfo) *pb.CollectiveInfo {
	tmp := new(pb.CollectiveInfo)
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
	tmp.Tags = info.Tags
	tmp.Assets = info.Assets
	tmp.Duration = &pb.DurationInfo{Begin: info.Date.Start, End: info.Date.Stop}
	return tmp
}

func (mine *CollectiveService) AddOne(ctx context.Context, in *pb.ReqCollectiveAdd, out *pb.ReplyCollectiveInfo) error {
	path := "collective.add"
	inLog(path, in)
	if len(in.Name) < 1 {
		out.Status = outError(path, "the name is empty ", pbstatus.ResultStatus_Empty)
		return nil
	}

	//tmp, _ := cache.Context().GetCollAlbumByName(in.Owner, in.Name)
	//if tmp != nil {
	//	out.Status = outError(path, "the name is repeated", pbstatus.ResultStatus_Repeated)
	//	return nil
	//}
	info, err := cache.Context().CreateCollAlbum(in.Name, in.Remark, in.Operator, in.Owner, uint8(in.Type), in.Duration.Begin, in.Duration.End)
	if err != nil {
		out.Status = outError(path, err.Error(), pbstatus.ResultStatus_DBException)
		return nil
	}
	out.Info = switchCollective(info)
	out.Status = outLog(path, out)
	return nil
}

func (mine *CollectiveService) GetOne(ctx context.Context, in *pb.RequestInfo, out *pb.ReplyCollectiveInfo) error {
	path := "collective.getOne"
	inLog(path, in)
	if len(in.Uid) < 1 {
		out.Status = outError(path, "the uid is empty ", pbstatus.ResultStatus_Empty)
		return nil
	}
	info, er := cache.Context().GetCollAlbum(in.Uid)
	if er != nil {
		out.Status = outError(path, "the collective not found ", pbstatus.ResultStatus_NotExisted)
		return nil
	}
	out.Info = switchCollective(info)
	out.Status = outLog(path, out)
	return nil
}

func (mine *CollectiveService) GetStatistic(ctx context.Context, in *pb.RequestFilter, out *pb.ReplyStatistic) error {
	path := "collective.getStatistic"
	inLog(path, in)
	if len(in.Field) < 1 {
		out.Status = outError(path, "the user is empty ", pbstatus.ResultStatus_Empty)
		return nil
	}

	out.Status = outLog(path, out)
	return nil
}

func (mine *CollectiveService) RemoveOne(ctx context.Context, in *pb.RequestInfo, out *pb.ReplyInfo) error {
	path := "collective.remove"
	inLog(path, in)
	if len(in.Uid) < 1 {
		out.Status = outError(path, "the uid is empty ", pbstatus.ResultStatus_Empty)
		return nil
	}
	info, er := cache.Context().GetCollAlbum(in.Uid)
	if er != nil {
		out.Status = outError(path, "the collective not found ", pbstatus.ResultStatus_NotExisted)
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

func (mine *CollectiveService) Search(ctx context.Context, in *pb.RequestInfo, out *pb.ReplyCollectiveList) error {
	path := "collective.search"
	inLog(path, in)

	out.Status = outLog(path, fmt.Sprintf("the length = %d", len(out.List)))
	return nil
}

func (mine *CollectiveService) GetListByFilter(ctx context.Context, in *pb.RequestFilter, out *pb.ReplyCollectiveList) error {
	path := "collective.getListByFilter"
	inLog(path, in)
	var list []*cache.CollAlbumInfo
	var err error
	if in.Field == "" {
		list = cache.Context().GetCollAlbumsByGroup(in.Owner)
	} else if in.Field == "type" {

	} else if in.Field == "user" {
		list = cache.Context().GetCollAlbums(in.Value)
	} else if in.Field == "array" {
		list = cache.Context().GetCollAlbumList(in.Values)
	} else {
		err = errors.New("the key not defined")
	}
	if err != nil {
		out.Status = outError(path, err.Error(), pbstatus.ResultStatus_DBException)
		return nil
	}
	out.List = make([]*pb.CollectiveInfo, 0, len(list))
	for _, value := range list {
		out.List = append(out.List, switchCollective(value))
	}

	out.Status = outLog(path, fmt.Sprintf("the length = %d", len(out.List)))
	return nil
}

func (mine *CollectiveService) UpdateBase(ctx context.Context, in *pb.ReqCollectiveUpdate, out *pb.ReplyInfo) error {
	path := "collective.updateBase"
	inLog(path, in)
	if len(in.Uid) < 1 {
		out.Status = outError(path, "the uid is empty ", pbstatus.ResultStatus_Empty)
		return nil
	}
	info, er := cache.Context().GetCollAlbum(in.Uid)
	if er != nil {
		out.Status = outError(path, er.Error(), pbstatus.ResultStatus_NotExisted)
		return nil
	}
	var err error
	if in.Name != info.Name && cache.Context().HadCollAlbum(info.Owner, in.Name) {
		out.Status = outError(path, "the collective name repeated ", pbstatus.ResultStatus_Repeated)
		return nil
	}
	err = info.UpdateBase(in.Name, in.Remark, in.Operator)
	if err != nil {
		out.Status = outError(path, err.Error(), pbstatus.ResultStatus_DBException)
		return nil
	}

	out.Status = outLog(path, out)
	return nil
}

func (mine *CollectiveService) UpdateByFilter(ctx context.Context, in *pb.RequestUpdate, out *pb.ReplyInfo) error {
	path := "collective.updateByFilter"
	inLog(path, in)
	if len(in.Uid) < 1 {
		out.Status = outError(path, "the uid is empty ", pbstatus.ResultStatus_Empty)
		return nil
	}
	info, er := cache.Context().GetCollAlbum(in.Uid)
	if er != nil {
		out.Status = outError(path, "the collective not found ", pbstatus.ResultStatus_NotExisted)
		return nil
	}
	var err error
	if in.Field == "cover" {
		err = info.UpdateCover(in.Value, in.Operator)
	} else if in.Field == "size" {
		size := parseStringToInt64(in.Value)
		if size > -1 {
			err = info.UpdateSize(uint64(size), in.Operator)
		}

	} else if in.Field == "assets" {
		err = info.UpdateAssets(in.Values, in.Operator)
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

func (mine *CollectiveService) UpdateStatus(ctx context.Context, in *pb.RequestIntFlag, out *pb.ReplyInfo) error {
	path := "collective.updateStatus"
	inLog(path, in)
	if len(in.Uid) < 1 {
		out.Status = outError(path, "the uid is empty ", pbstatus.ResultStatus_Empty)
		return nil
	}
	info, er := cache.Context().GetCollAlbum(in.Uid)
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

func (mine *CollectiveService) AppendAsset(ctx context.Context, in *pb.RequestList, out *pb.ReplyList) error {
	path := "collective.appendAsset"
	inLog(path, in)
	if len(in.Uid) < 1 {
		out.Status = outError(path, "the uid is empty ", pbstatus.ResultStatus_Empty)
		return nil
	}
	info, er := cache.Context().GetCollAlbum(in.Uid)
	if er != nil {
		out.Status = outError(path, er.Error(), pbstatus.ResultStatus_NotExisted)
		return nil
	}

	err := info.AppendAssets(in.List, in.Operator)
	if err != nil {
		out.Status = outError(path, err.Error(), pbstatus.ResultStatus_DBException)
		return nil
	}
	out.Uid = in.Uid
	out.List = info.Assets
	out.Status = outLog(path, out)
	return nil
}

func (mine *CollectiveService) SubtractAsset(ctx context.Context, in *pb.RequestList, out *pb.ReplyList) error {
	path := "collective.subtractAsset"
	inLog(path, in)
	if len(in.Uid) < 1 {
		out.Status = outError(path, "the uid is empty ", pbstatus.ResultStatus_Empty)
		return nil
	}
	info, er := cache.Context().GetCollAlbum(in.Uid)
	if er != nil {
		out.Status = outError(path, er.Error(), pbstatus.ResultStatus_NotExisted)
		return nil
	}

	err := info.SubtractAssets(in.List, in.Operator)
	if err != nil {
		out.Status = outError(path, err.Error(), pbstatus.ResultStatus_DBException)
		return nil
	}
	out.Uid = in.Uid
	out.List = info.Assets
	out.Status = outLog(path, out)
	return nil
}
