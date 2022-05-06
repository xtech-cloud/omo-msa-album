package grpc

import (
	"context"
	"errors"
	"fmt"
	pb "github.com/xtech-cloud/omo-msp-album/proto/album"
	pbstatus "github.com/xtech-cloud/omo-msp-status/proto/status"
	"omo.msa.album/cache"
)

type PanoramaService struct{}

func switchPanorama(info *cache.PanoramaInfo) *pb.PanoramaInfo {
	tmp := new(pb.PanoramaInfo)
	tmp.Uid = info.UID
	tmp.Id = info.ID
	tmp.Created = info.CreateTime.Unix()
	tmp.Updated = info.UpdateTime.Unix()
	tmp.Operator = info.Operator
	tmp.Creator = info.Creator
	tmp.Name = info.Name
	tmp.Remark = info.Remark
	tmp.Content = info.Content
	tmp.Owner = info.Owner
	return tmp
}

func (mine *PanoramaService) AddOne(ctx context.Context, in *pb.ReqPanoramaAdd, out *pb.ReplyPanoramaInfo) error {
	path := "panorama.add"
	inLog(path, in)
	if len(in.Name) < 1 {
		out.Status = outError(path, "the name is empty ", pbstatus.ResultStatus_Empty)
		return nil
	}

	info, err := cache.Context().CreatePanorama(in.Name, in.Remark, in.Content, in.Owner, in.Operator)
	if err != nil {
		out.Status = outError(path, err.Error(), pbstatus.ResultStatus_DBException)
		return nil
	}
	out.Info = switchPanorama(info)
	out.Status = outLog(path, out)
	return nil
}

func (mine *PanoramaService) GetOne(ctx context.Context, in *pb.RequestInfo, out *pb.ReplyPanoramaInfo) error {
	path := "panorama.getOne"
	inLog(path, in)
	if len(in.Uid) < 1 {
		out.Status = outError(path, "the uid is empty ", pbstatus.ResultStatus_Empty)
		return nil
	}
	info, er := cache.Context().GetPanorama(in.Uid)
	if er != nil {
		out.Status = outError(path, "the panorama not found ", pbstatus.ResultStatus_NotExisted)
		return nil
	}
	out.Info = switchPanorama(info)
	out.Status = outLog(path, out)
	return nil
}

func (mine *PanoramaService) GetStatistic(ctx context.Context, in *pb.RequestFilter, out *pb.ReplyStatistic) error {
	path := "panorama.getStatistic"
	inLog(path, in)
	if len(in.Field) < 1 {
		out.Status = outError(path, "the user is empty ", pbstatus.ResultStatus_Empty)
		return nil
	}

	out.Status = outLog(path, out)
	return nil
}

func (mine *PanoramaService) RemoveOne(ctx context.Context, in *pb.RequestInfo, out *pb.ReplyInfo) error {
	path := "panorama.remove"
	inLog(path, in)
	if len(in.Uid) < 1 {
		out.Status = outError(path, "the uid is empty ", pbstatus.ResultStatus_Empty)
		return nil
	}
	er := cache.Context().RemovePanorama(in.Uid, in.Operator)
	if er != nil {
		out.Status = outError(path, "the panorama not found ", pbstatus.ResultStatus_NotExisted)
		return nil
	}

	out.Uid = in.Uid
	out.Status = outLog(path, out)
	return nil
}

func (mine *PanoramaService) Search(ctx context.Context, in *pb.RequestInfo, out *pb.ReplyPanoramaList) error {
	path := "panorama.search"
	inLog(path, in)

	out.Status = outLog(path, fmt.Sprintf("the length = %d", len(out.List)))
	return nil
}

func (mine *PanoramaService) GetListByFilter(ctx context.Context, in *pb.RequestFilter, out *pb.ReplyPanoramaList) error {
	path := "panorama.getListByFilter"
	inLog(path, in)
	var list []*cache.PanoramaInfo
	var err error
	if in.Field == "" {
		list, err = cache.Context().GetPanoramas(in.Owner)
	} else {
		err = errors.New("the key not defined")
	}
	if err != nil {
		out.Status = outError(path, err.Error(), pbstatus.ResultStatus_DBException)
		return nil
	}
	out.List = make([]*pb.PanoramaInfo, 0, len(list))
	for _, value := range list {
		out.List = append(out.List, switchPanorama(value))
	}

	out.Status = outLog(path, fmt.Sprintf("the length = %d", len(out.List)))
	return nil
}

func (mine *PanoramaService) UpdateBase(ctx context.Context, in *pb.ReqPanoramaUpdate, out *pb.ReplyInfo) error {
	path := "panorama.updateBase"
	inLog(path, in)
	if len(in.Uid) < 1 {
		out.Status = outError(path, "the uid is empty ", pbstatus.ResultStatus_Empty)
		return nil
	}
	info, er := cache.Context().GetPanorama(in.Uid)
	if er != nil {
		out.Status = outError(path, er.Error(), pbstatus.ResultStatus_NotExisted)
		return nil
	}
	var err error
	err = info.UpdateContent(in.Content, in.Operator)
	if err != nil {
		out.Status = outError(path, err.Error(), pbstatus.ResultStatus_DBException)
		return nil
	}

	out.Status = outLog(path, out)
	return nil
}

func (mine *PanoramaService) UpdateByFilter(ctx context.Context, in *pb.RequestUpdate, out *pb.ReplyInfo) error {
	path := "panorama.updateByFilter"
	inLog(path, in)
	if len(in.Uid) < 1 {
		out.Status = outError(path, "the uid is empty ", pbstatus.ResultStatus_Empty)
		return nil
	}
	_, er := cache.Context().GetPanorama(in.Uid)
	if er != nil {
		out.Status = outError(path, "the panorama not found ", pbstatus.ResultStatus_NotExisted)
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
