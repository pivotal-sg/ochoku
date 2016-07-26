package reviews

import (
	"fmt"
	"log"

	"github.com/micro/go-micro"
	"github.com/micro/go-micro/cmd"
	"github.com/micro/go-platform/auth"
	proto "github.com/pivotal-sg/ochoku/reviews/proto"
	"github.com/pivotal-sg/ochoku/reviews/storage"
)

const ServerName = "pivotal.io.ochoku"
const ServiceName = "reviews"
const Version = "0.1.0"

// NewService returns a ReviewService with server, it takes int the command
// parser and sets up the storage and auth
func NewService(c cmd.Cmd) micro.Service {
	c.Init()
	ReviewConfig := Config(ConfigFileName)

	storageFile := ReviewConfig.Get("reviews", "storage", "filename").String("reviews.db")
	fmt.Println("in NewService, storageFile", storageFile)
	storage, err := storage.New(storageFile)
	if err != nil {
		log.Fatalf("Couldn't create bolt storage backend: \n%v\n", err)
	}

	opts := []micro.Option{
		micro.Name(ServiceName),
		micro.Version(Version),
		micro.Cmd(c),
	}

	if ReviewConfig.Get("reviews", "authFlag").Bool(false) {
		a := auth.NewAuth(
			auth.Id(ReviewConfig.Get("reviews", "auth", "client_id").String("")),
			auth.Secret(ReviewConfig.Get("reviews", "auth", "secret").String("")),
		)
		opts = append(opts, micro.WrapHandler(auth.HandlerWrapper(a)))
	}

	service := micro.NewService(opts...)

	service.Init()

	reviewService := ReviewService{Store: storage}

	proto.RegisterReviewerHandler(service.Server(), &reviewService)

	return service
}

func NewClient(c cmd.Cmd) proto.ReviewerClient {
	c.Init()
	ReviewConfig := Config(ConfigFileName)

	opts := []micro.Option{
		micro.Name(ServiceName + ".client"),
		micro.Version(Version),
	}

	if ReviewConfig.Get("reviews", "authFlag").Bool(false) {
		a := auth.NewAuth(
			auth.Id(ReviewConfig.Get("reviews", "auth", "client_id").String("")),
			auth.Secret(ReviewConfig.Get("reviews", "auth", "secret").String("")),
		)
		opts = append(opts, micro.WrapClient(auth.ClientWrapper(a)))
	}
	service := micro.NewService(opts...)
	return proto.NewReviewerClient(ServiceName, service.Client())
}
