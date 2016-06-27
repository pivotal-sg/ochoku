package reviews_test

import (
	"reflect"
	"testing"

	"github.com/pivotal-sg/ochoku/reviews"
	proto "github.com/pivotal-sg/ochoku/reviews/proto"
	"golang.org/x/net/context"
)

var reviewServiceObject reviews.ReviewService = reviews.ReviewService{}

func TestReviewChocolate(t *testing.T) {
	reviewRequest := proto.ReviewRequest{
		Reviewer: "James",
		Name:     "Hershy's Dark Cardboard",
		Review:   "I ate the wrapper as well, and it tasted better than the chocolate",
		Rating:   -5,
	}
	response, err := reviewServiceObject.Review(context.TODO(), &reviewRequest)

	if err != nil {
		t.Errorf("expected the response to not have an error, it was %v", err)
	}

	if !response.Success {
		t.Errorf("expected the response to be successful, it wasn't")
	}
}

func TestValidInputs(t *testing.T) {

	var testValues = []struct {
		Input    *proto.ReviewRequest
		Expected *proto.StatusResponse
	}{
		{&proto.ReviewRequest{
			Reviewer: "James",
			Name:     "",
			Review:   "I ate the wrapper as well, and it tasted better than the chocolate",
			Rating:   -5},
			&proto.StatusResponse{
				Message: "Name is Missing",
				Success: false},
		},
		{&proto.ReviewRequest{
			Reviewer: "",
			Name:     "Valrona",
			Review:   "I ate the wrapper as well, and it tasted worse than the chocolate",
			Rating:   1},
			&proto.StatusResponse{
				Message: "Reviewer is Missing",
				Success: false,
			},
		},
	}

	for _, testValue := range testValues {
		response, err := reviewServiceObject.Review(context.TODO(), testValue.Input)

		if err != nil {
			t.Errorf("expected the response for '%v' to not have an error, it was %v", testValue.Input, err)
		}

		if !reflect.DeepEqual(testValue.Expected, response) {
			t.Errorf("Expected response to equal '%v', was '%v'", testValue.Expected, response)
		}
	}
}