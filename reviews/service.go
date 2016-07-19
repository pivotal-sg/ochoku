package reviews

import (
	"fmt"
	"log"

	"github.com/micro/cli"
	"github.com/micro/go-micro"
	"github.com/micro/go-platform/auth"
	proto "github.com/pivotal-sg/ochoku/reviews/proto"
	"github.com/pivotal-sg/ochoku/reviews/storage"
)

const ServerName = "pivotal.io.ochoku"
const ServiceName = "reviews"
const Version = "0.1.0"

// NewService returns a ReviewService with server
func NewService(storageFile string, a auth.Auth) micro.Service {
	service := micro.NewService(
		micro.WrapHandler(auth.HandlerWrapper(a)),
		micro.Name(ServiceName),
		micro.Version(Version),
		micro.Flags(
			cli.StringFlag{
				Name:        "config",
				EnvVar:      "MICRO_REVIEWS_CONFIG_FILE",
				Usage:       "Path to the file backed configuration for the reviews service",
				Value:       DefaultConfigFileName,
				Destination: &configFileName,
			}),
	)

	service.Init(micro.Action(func(c *cli.Context) { fmt.Println(c.String("config")) }))

	storage, err := storage.New(storageFile)

	if err != nil {
		log.Fatalf("Couldn't create bolt storage backend: \n%v\n", err)
	}
	reviewService := ReviewService{Store: storage}

	proto.RegisterReviewerHandler(service.Server(), &reviewService)

	return service
}

func NewClient() proto.ReviewerClient {
	service := micro.NewService(
		micro.Name(ServiceName+".client"),
		micro.Version(Version),
	)
	client := proto.NewReviewerClient(ServiceName, service.Client())
	return client
}
