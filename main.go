package main

import (
	"crypto/md5"
	"encoding/hex"
	"github.com/micro/go-micro/v2"
	"github.com/micro/go-micro/v2/logger"
	_ "github.com/micro/go-plugins/registry/consul/v2"
	_ "github.com/micro/go-plugins/registry/etcdv3/v2"
	proto "github.com/xtech-cloud/omo-msp-album/proto/album"
	"io"
	"omo.msa.album/cache"
	"omo.msa.album/config"
	"omo.msa.album/grpc"
	"os"
	"path/filepath"
	"time"
)

var (
	BuildVersion string
	BuildTime    string
	CommitID     string
)

func main() {
	config.Setup()
	err := cache.InitData()
	if err != nil {
		panic(err)
	}
	// New Service
	service := micro.NewService(
		micro.Name("omo.msa.album"),
		micro.Version(BuildVersion),
		micro.RegisterTTL(time.Second*time.Duration(config.Schema.Service.TTL)),
		micro.RegisterInterval(time.Second*time.Duration(config.Schema.Service.Interval)),
		micro.Address(config.Schema.Service.Address),
	)
	// Initialise service
	service.Init()
	// Register Handler
	_ = proto.RegisterAlbumServiceHandler(service.Server(), new(grpc.AlbumService))
	_ = proto.RegisterCollectiveServiceHandler(service.Server(), new(grpc.CollectiveService))
	_ = proto.RegisterExhibitServiceHandler(service.Server(), new(grpc.ExhibitService))
	_ = proto.RegisterPanoramaServiceHandler(service.Server(), new(grpc.PanoramaService))
	_ = proto.RegisterPhotocopyServiceHandler(service.Server(), new(grpc.PhotocopyService))
	_ = proto.RegisterStyleServiceHandler(service.Server(), new(grpc.StyleService))
	_ = proto.RegisterFrameServiceHandler(service.Server(), new(grpc.FrameService))
	_ = proto.RegisterCompositionServiceHandler(service.Server(), new(grpc.CompositionService))
	_ = proto.RegisterPageServiceHandler(service.Server(), new(grpc.PageService))
	_ = proto.RegisterSheetServiceHandler(service.Server(), new(grpc.SheetService))

	app, _ := filepath.Abs(os.Args[0])

	logger.Info("-------------------------------------------------------------")
	logger.Info("- Micro Service Agent -> Run")
	logger.Info("-------------------------------------------------------------")
	logger.Infof("- version      : %s", BuildVersion)
	logger.Infof("- application  : %s", app)
	logger.Infof("- md5          : %s", md5hex(app))
	logger.Infof("- build        : %s", BuildTime)
	logger.Infof("- commit       : %s", CommitID)
	logger.Info("-------------------------------------------------------------")
	// Run service
	if err := service.Run(); err != nil {
		logger.Fatal(err)
	}
}

func md5hex(_file string) string {
	h := md5.New()

	f, err := os.Open(_file)
	if err != nil {
		return ""
	}
	defer f.Close()
	io.Copy(h, f)
	return hex.EncodeToString(h.Sum(nil))
}
