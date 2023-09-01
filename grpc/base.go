package grpc

import (
	"encoding/json"
	"github.com/micro/go-micro/v2/logger"
	pb "github.com/xtech-cloud/omo-msp-album/proto/album"
	pbstatus "github.com/xtech-cloud/omo-msp-status/proto/status"
	"strconv"
	"strings"
)

func inLog(name, data interface{}) {
	bytes, _ := json.Marshal(data)
	msg := ByteString(bytes)
	logger.Infof("[in.%s]:data = %s", name, msg)
}

func ByteString(p []byte) string {
	for i := 0; i < len(p); i++ {
		if p[i] == 0 {
			return string(p[0:i])
		}
	}
	return string(p)
}

func outError(name, msg string, code pbstatus.ResultStatus) *pb.ReplyStatus {
	logger.Warnf("[error.%s]:code = %d, msg = %s", name, code, msg)
	tmp := &pb.ReplyStatus{
		Code:  uint32(code),
		Error: msg,
	}
	return tmp
}

func outLog(name, data interface{}) *pb.ReplyStatus {
	bytes, _ := json.Marshal(data)
	msg := ByteString(bytes)
	logger.Infof("[out.%s]:data = %s", name, msg)
	tmp := &pb.ReplyStatus{
		Code:  0,
		Error: "",
	}
	return tmp
}

func parseString(src string, sep string) (string, int) {
	arr := strings.Split(src, sep)
	if len(arr) < 2 {
		return "", -1
	}
	st, er := strconv.ParseInt(arr[1], 10, 32)
	if er != nil {
		return "", -1
	}
	return arr[0], int(st)
}

func parseStringToInt(src string) int {
	if src == "" {
		return -1
	}
	st, er := strconv.ParseInt(src, 10, 32)
	if er != nil {
		return -1
	}
	return int(st)
}

func parseStringToInt64(src string) int64 {
	if src == "" {
		return -1
	}
	st, er := strconv.ParseInt(src, 10, 64)
	if er != nil {
		return -1
	}
	return st
}
