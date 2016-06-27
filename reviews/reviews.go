package reviews

import (
	"github.com/micro/go-micro/client"
	proto "github.com/pivotal-sg/ochoku/reviews/proto"
	"golang.org/x/net/context"
)

type ReviewService struct{}

func (*ReviewService) Review(c context.Context, reviewRequest *proto.ReviewRequest, opts ...client.CallOption) (*proto.StatusResponse, error) {

	if reviewRequest.Reviewer == "" {
		return &proto.StatusResponse{
			Message: "Reviewer is Missing",
			Success: false,
		}, nil
	}

	if reviewRequest.Name == "" {
		return &proto.StatusResponse{
			Message: "Name is Missing",
			Success: false,
		}, nil
	}

	return &proto.StatusResponse{
		Message: "All Good!",
		Success: true,
	}, nil
}
