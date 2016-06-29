package reviews

import (
	"log"

	"github.com/micro/go-micro"
	proto "github.com/pivotal-sg/ochoku/reviews/proto"
	"github.com/pivotal-sg/ochoku/reviews/storage"
)

var ServiceName = "pivotal.io.ochoku.reviews"
var Version = "0.1.0"

// NewService returns a ReviewService with server
func NewService(storageFile string) micro.Service {
	service := micro.NewService(
		micro.Name(ServiceName),
		micro.Version(Version),
	)

	service.Init()

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
