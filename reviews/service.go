package reviews

import (
	"log"

	"github.com/micro/go-micro"
	"github.com/micro/go-micro/client"
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
	storage, err := storage.New(storageFile)
	if err != nil {
		log.Fatalf("Couldn't create bolt storage backend: \n%v\n", err)
	}

	a := auth.NewAuth(
		auth.Id(ReviewConfig.Get("reviews", "auth", "client_id").String("")),
		auth.Secret(ReviewConfig.Get("reviews", "auth", "secret").String("")),
	)

	service := micro.NewService(
		micro.WrapHandler(auth.HandlerWrapper(a)),
		micro.Name(ServiceName),
		micro.Version(Version),
		micro.Cmd(c),
	)

	service.Init()

	reviewService := ReviewService{Store: storage}

	proto.RegisterReviewerHandler(service.Server(), &reviewService)

	return service
}

func NewClient(c cmd.Cmd) (client.Client, auth.Auth) {
	c.Init()
	ReviewConfig := Config(ConfigFileName)

	a := auth.NewAuth(
		auth.Id(ReviewConfig.Get("reviews", "auth", "client_id").String("")),
		auth.Secret(ReviewConfig.Get("reviews", "auth", "secret").String("")),
	)

	service := micro.NewService(
		micro.Name(ServiceName+".client"),
		micro.Version(Version),
		micro.WrapClient(auth.ClientWrapper(a)),
	)

	client := service.Client()
	return client, a
}
