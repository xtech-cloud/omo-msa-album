package grpc

import (
	"context"
	"errors"
	"fmt"
	pb "github.com/xtech-cloud/omo-msp-album/proto/album"
	pbstatus "github.com/xtech-cloud/omo-msp-status/proto/status"
	"omo.msa.album/cache"
)

type CertificateService struct{}

func switchCertificate(info *cache.CertificateInfo) *pb.CertificateInfo {
	tmp := new(pb.CertificateInfo)
	tmp.Uid = info.UID
	tmp.Id = info.ID
	tmp.Created = info.Created
	tmp.Updated = info.Updated
	tmp.Operator = info.Operator
	tmp.Creator = info.Creator
	tmp.Name = info.Name
	tmp.Remark = info.Remark
	tmp.Sn = info.SN
	tmp.Image = info.Image
	tmp.Style = info.Style
	tmp.Target = info.Target
	tmp.Scene = info.Scene
	tmp.Status = uint32(info.Status)
	tmp.Type = uint32(info.Type)
	tmp.Tags = info.Tags
	tmp.Assets = info.Assets
	if info.Contact != nil {
		tmp.Contact = &pb.ContactInfo{
			Name:    info.Contact.Name,
			Phone:   info.Contact.Phone,
			Address: info.Contact.Address,
			Remark:  info.Contact.Remark,
		}
	} else {
		tmp.Contact = &pb.ContactInfo{}
	}

	return tmp
}

func (mine *CertificateService) AddOne(ctx context.Context, in *pb.ReqCertificateAdd, out *pb.ReplyCertificateInfo) error {
	path := "certificate.add"
	inLog(path, in)
	if len(in.Name) < 1 {
		out.Status = outError(path, "the name is empty ", pbstatus.ResultStatus_Empty)
		return nil
	}

	info, err := cache.Context().CreateCertificate(in)
	if err != nil {
		out.Status = outError(path, err.Error(), pbstatus.ResultStatus_DBException)
		return nil
	}
	out.Info = switchCertificate(info)
	out.Status = outLog(path, out)
	return nil
}

func (mine *CertificateService) GetOne(ctx context.Context, in *pb.RequestInfo, out *pb.ReplyCertificateInfo) error {
	path := "certificate.getOne"
	inLog(path, in)
	if len(in.Uid) < 1 {
		out.Status = outError(path, "the uid is empty ", pbstatus.ResultStatus_Empty)
		return nil
	}
	info, er := cache.Context().GetCertificate(in.Uid)
	if er != nil {
		out.Status = outError(path, "the style not found ", pbstatus.ResultStatus_NotExisted)
		return nil
	}
	out.Info = switchCertificate(info)
	out.Status = outLog(path, out)
	return nil
}

func (mine *CertificateService) UpdateBase(ctx context.Context, in *pb.ReqCertificateUpdate, out *pb.ReplyInfo) error {
	path := "certificate.updateBase"
	inLog(path, in)
	if len(in.Uid) < 1 {
		out.Status = outError(path, "the uid is empty ", pbstatus.ResultStatus_Empty)
		return nil
	}
	info, er := cache.Context().GetCertificate(in.Uid)
	if er != nil {
		out.Status = outError(path, er.Error(), pbstatus.ResultStatus_NotExisted)
		return nil
	}
	var err error
	err = info.UpdateBase(in.Name, in.Remark, in.Operator, in.Tags)
	if err != nil {
		out.Status = outError(path, err.Error(), pbstatus.ResultStatus_DBException)
		return nil
	}

	out.Status = outLog(path, out)
	return nil
}

func (mine *CertificateService) GetStatistic(ctx context.Context, in *pb.RequestFilter, out *pb.ReplyStatistic) error {
	path := "certificate.getStatistic"
	inLog(path, in)
	if len(in.Field) < 1 {
		out.Status = outError(path, "the user is empty ", pbstatus.ResultStatus_Empty)
		return nil
	}

	out.Status = outLog(path, out)
	return nil
}

func (mine *CertificateService) RemoveOne(ctx context.Context, in *pb.RequestInfo, out *pb.ReplyInfo) error {
	path := "certificate.remove"
	inLog(path, in)
	if len(in.Uid) < 1 {
		out.Status = outError(path, "the uid is empty ", pbstatus.ResultStatus_Empty)
		return nil
	}
	info, er := cache.Context().GetStyle(in.Uid)
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

func (mine *CertificateService) Search(ctx context.Context, in *pb.RequestInfo, out *pb.ReplyCertificateList) error {
	path := "certificate.search"
	inLog(path, in)

	out.Status = outLog(path, fmt.Sprintf("the length = %d", len(out.List)))
	return nil
}

func (mine *CertificateService) GetListByFilter(ctx context.Context, in *pb.RequestFilter, out *pb.ReplyCertificateList) error {
	path := "certificate.getListByFilter"
	inLog(path, in)
	var list []*cache.CertificateInfo
	var err error
	if in.Field == "target" {
		list = cache.Context().GetCertificatesByTarget(in.Value)
	} else if in.Field == "scene" {
		list = cache.Context().GetCertificatesByScene(in.Value)
	} else if in.Field == "style" {
		list = cache.Context().GetCertificatesByStyle(in.Owner, in.Value)
	} else if in.Field == "contact" {
		list = cache.Context().GetCertificateByContact(in.Owner, in.Value)
	} else if in.Field == "array" {
		list = cache.Context().GetCertificateByArray(in.Values)
	} else if in.Field == "batch" {
		if len(in.Values) == 3 {
			num := parseStringToInt(in.Values[0])
			quote := in.Values[1]
			end := parseStringToInt64(in.Values[2])
			list, err = cache.Context().BatchCertificate(in.Value, in.Owner, quote, in.Operator, uint32(num), end)
		}
	} else {
		err = errors.New("the key not defined")
	}
	if err != nil {
		out.Status = outError(path, err.Error(), pbstatus.ResultStatus_DBException)
		return nil
	}
	out.List = make([]*pb.CertificateInfo, 0, len(list))
	for _, value := range list {
		out.List = append(out.List, switchCertificate(value))
	}

	out.Status = outLog(path, fmt.Sprintf("the length = %d", len(out.List)))
	return nil
}

func (mine *CertificateService) UpdateByFilter(ctx context.Context, in *pb.RequestUpdate, out *pb.ReplyInfo) error {
	path := "certificate.updateByFilter"
	inLog(path, in)
	if len(in.Uid) < 1 {
		out.Status = outError(path, "the uid is empty ", pbstatus.ResultStatus_Empty)
		return nil
	}
	info, er := cache.Context().GetCertificate(in.Uid)
	if er != nil {
		out.Status = outError(path, er.Error(), pbstatus.ResultStatus_NotExisted)
		return nil
	}
	if in.Field == "contact" {
		if len(in.Values) == 4 {
			er = info.UpdateContact(in.Values[0], in.Values[1], in.Values[2], in.Values[3], in.Operator)
		} else {
			er = errors.New("the param count error")
		}
	} else if in.Field == "target" {
		er = info.UpdateTarget(in.Value, in.Operator)
	}
	if er != nil {
		out.Status = outError(path, er.Error(), pbstatus.ResultStatus_DBException)
		return nil
	}
	out.Status = outLog(path, out)
	return nil
}

func (mine *CertificateService) UpdateStatus(ctx context.Context, in *pb.RequestIntFlag, out *pb.ReplyInfo) error {
	path := "certificate.updateStatus"
	inLog(path, in)
	if len(in.Uid) < 1 {
		out.Status = outError(path, "the uid is empty ", pbstatus.ResultStatus_Empty)
		return nil
	}
	info, er := cache.Context().GetCertificate(in.Uid)
	if er != nil {
		out.Status = outError(path, er.Error(), pbstatus.ResultStatus_NotExisted)
		return nil
	}
	er = info.UpdateStatus(in.Operator, uint8(in.Flag))
	if er != nil {
		out.Status = outError(path, er.Error(), pbstatus.ResultStatus_DBException)
		return nil
	}
	out.Status = outLog(path, out)
	return nil
}
