package grpc

import (
	"context"
	"errors"
	"fmt"
	pb "github.com/xtech-cloud/omo-msp-album/proto/album"
	pbstatus "github.com/xtech-cloud/omo-msp-status/proto/status"
	"omo.msa.album/cache"
)

type AlbumService struct{}

func switchAlbum(info *cache.AlbumInfo) *pb.AlbumInfo {
	tmp := new(pb.AlbumInfo)
	tmp.Uid = info.UID
	tmp.Id = info.ID
	tmp.Created = info.CreateTime.Unix()
	tmp.Updated = info.UpdateTime.Unix()
	tmp.Operator = info.Operator
	tmp.Creator = info.Creator
	tmp.Name = info.Name
	tmp.Remark = info.Remark
	tmp.Cover = info.Cover
	tmp.Status = uint32(info.Status)
	tmp.Tags = info.Tags
	tmp.Type = uint32(info.Kind)
	tmp.Style = uint32(info.Style)
	tmp.Size = uint64(info.Size)
	tmp.Limit = uint32(info.MaxCount)
	tmp.Star = uint32(info.StarCount)
	tmp.Assets = info.Assets
	tmp.Targets = info.Targets
	return tmp
}

func (mine *AlbumService) AddOne(ctx context.Context, in *pb.ReqAlbumAdd, out *pb.ReplyAlbumInfo) error {
	path := "album.add"
	inLog(path, in)
	if len(in.Name) < 1 {
		out.Status = outError(path, "the name is empty ", pbstatus.ResultStatus_Empty)
		return nil
	}

	info, err := cache.Context().CreateAlbum(in.Name, in.Remark, in.Operator, in.Location, uint8(in.Type), uint16(in.Style), in.Targets)
	if err != nil {
		out.Status = outError(path, err.Error(), pbstatus.ResultStatus_DBException)
		return nil
	}
	out.Info = switchAlbum(info)
	out.Status = outLog(path, out)
	return nil
}

func (mine *AlbumService) GetOne(ctx context.Context, in *pb.RequestInfo, out *pb.ReplyAlbumInfo) error {
	path := "album.getOne"
	inLog(path, in)
	if len(in.Uid) < 1 {
		out.Status = outError(path, "the uid is empty ", pbstatus.ResultStatus_Empty)
		return nil
	}
	info, er := cache.Context().GetAlbum(in.Uid)
	if er != nil {
		out.Status = outError(path, "the album not found ", pbstatus.ResultStatus_NotExisted)
		return nil
	}
	out.Info = switchAlbum(info)
	out.Status = outLog(path, out)
	return nil
}

func (mine *AlbumService) GetStatistic(ctx context.Context, in *pb.RequestFilter, out *pb.ReplyStatistic) error {
	path := "album.getStatistic"
	inLog(path, in)
	if len(in.Field) < 1 {
		out.Status = outError(path, "the user is empty ", pbstatus.ResultStatus_Empty)
		return nil
	}

	out.Status = outLog(path, out)
	return nil
}

func (mine *AlbumService) RemoveOne(ctx context.Context, in *pb.RequestInfo, out *pb.ReplyInfo) error {
	path := "album.remove"
	inLog(path, in)
	if len(in.Uid) < 1 {
		out.Status = outError(path, "the uid is empty ", pbstatus.ResultStatus_Empty)
		return nil
	}
	info, er := cache.Context().GetAlbum(in.Uid)
	if er != nil {
		out.Status = outError(path, "the album not found ", pbstatus.ResultStatus_NotExisted)
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

func (mine *AlbumService) Search(ctx context.Context, in *pb.RequestInfo, out *pb.ReplyAlbumList) error {
	path := "album.search"
	inLog(path, in)

	out.Status = outLog(path, fmt.Sprintf("the length = %d", len(out.List)))
	return nil
}

func (mine *AlbumService) GetListByFilter(ctx context.Context, in *pb.RequestFilter, out *pb.ReplyAlbumList) error {
	path := "album.getListByFilter"
	inLog(path, in)
	var list []*cache.AlbumInfo
	var err error
	if in.Field == "" {
		list = cache.Context().GetAlbumsByUser(in.Owner)
	} else if in.Field == "status" {
		st := parseStringToInt(in.Value)
		list = cache.Context().GetAlbums(in.Owner, cache.AlbumStatus(st))
	} else if in.Field == "targets" {
		list = cache.Context().GetAlbumsByTargets(in.Values)
	} else if in.Field == "array" {
		list = cache.Context().GetAlbumList(in.Values)
	} else {
		err = errors.New("the key not defined")
	}
	if err != nil {
		out.Status = outError(path, err.Error(), pbstatus.ResultStatus_DBException)
		return nil
	}
	out.List = make([]*pb.AlbumInfo, 0, len(list))
	for _, value := range list {
		out.List = append(out.List, switchAlbum(value))
	}

	out.Status = outLog(path, fmt.Sprintf("the length = %d", len(out.List)))
	return nil
}

func (mine *AlbumService) UpdateBase(ctx context.Context, in *pb.ReqAlbumUpdate, out *pb.ReplyInfo) error {
	path := "album.updateBase"
	inLog(path, in)
	if len(in.Uid) < 1 {
		out.Status = outError(path, "the uid is empty ", pbstatus.ResultStatus_Empty)
		return nil
	}
	info, er := cache.Context().GetAlbum(in.Uid)
	if er != nil {
		out.Status = outError(path, er.Error(), pbstatus.ResultStatus_NotExisted)
		return nil
	}
	var err error
	err = info.UpdateBase(in.Name, in.Remark, in.Operator, in.Passwords)
	if err != nil {
		out.Status = outError(path, err.Error(), pbstatus.ResultStatus_DBException)
		return nil
	}

	out.Status = outLog(path, out)
	return nil
}

func (mine *AlbumService) UpdateByFilter(ctx context.Context, in *pb.RequestUpdate, out *pb.ReplyInfo) error {
	path := "album.updateByFilter"
	inLog(path, in)
	if len(in.Uid) < 1 {
		out.Status = outError(path, "the uid is empty ", pbstatus.ResultStatus_Empty)
		return nil
	}
	info, er := cache.Context().GetAlbum(in.Uid)
	if er != nil {
		out.Status = outError(path, "the album not found ", pbstatus.ResultStatus_NotExisted)
		return nil
	}
	var err error
	if in.Field == "cover" {
		err = info.UpdateCover(in.Value, in.Operator)
	} else if in.Field == "targets" {
		err = info.UpdateTargets(in.Operator, in.Values)
	} else if in.Field == "size" {
		size := parseStringToInt(in.Value)
		err = info.UpdateSize(uint64(size))
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

func (mine *AlbumService) UpdateStatus(ctx context.Context, in *pb.RequestIntFlag, out *pb.ReplyInfo) error {
	path := "album.updateStatus"
	inLog(path, in)
	if len(in.Uid) < 1 {
		out.Status = outError(path, "the uid is empty ", pbstatus.ResultStatus_Empty)
		return nil
	}
	info, er := cache.Context().GetAlbum(in.Uid)
	if er != nil {
		out.Status = outError(path, er.Error(), pbstatus.ResultStatus_NotExisted)
		return nil
	}
	err := info.UpdateStatus(in.Operator, cache.AlbumStatus(in.Flag))
	if err != nil {
		out.Status = outError(path, err.Error(), pbstatus.ResultStatus_DBException)
		return nil
	}
	out.Status = outLog(path, out)
	return nil
}

func (mine *AlbumService) AppendAsset(ctx context.Context, in *pb.RequestList, out *pb.ReplyList) error {
	path := "album.appendAsset"
	inLog(path, in)
	if len(in.Uid) < 1 {
		out.Status = outError(path, "the uid is empty ", pbstatus.ResultStatus_Empty)
		return nil
	}
	info, er := cache.Context().GetAlbum(in.Uid)
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

func (mine *AlbumService) SubtractAsset(ctx context.Context, in *pb.RequestList, out *pb.ReplyList) error {
	path := "album.subtractAsset"
	inLog(path, in)
	if len(in.Uid) < 1 {
		out.Status = outError(path, "the uid is empty ", pbstatus.ResultStatus_Empty)
		return nil
	}
	info, er := cache.Context().GetAlbum(in.Uid)
	if er != nil {
		out.Status = outError(path, er.Error(), pbstatus.ResultStatus_NotExisted)
		return nil
	}

	err := info.RemoveAssets(in.List, in.Operator)
	if err != nil {
		out.Status = outError(path, err.Error(), pbstatus.ResultStatus_DBException)
		return nil
	}
	out.Uid = in.Uid
	out.List = info.Assets
	out.Status = outLog(path, out)
	return nil
}
