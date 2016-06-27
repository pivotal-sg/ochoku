package reviews_test

import (
	"errors"
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
		{
			&proto.ReviewRequest{
				Reviewer: "James",
				Name:     "",
				Review:   "I ate the wrapper as well, and it tasted better than the chocolate",
				Rating:   -5},
			&proto.StatusResponse{
				Message: `{"name":"missing"}`,
				Success: false},
		},
		{
			&proto.ReviewRequest{
				Reviewer: "",
				Name:     "Valrona",
				Review:   "I ate the wrapper as well, and it tasted worse than the chocolate",
				Rating:   1},
			&proto.StatusResponse{
				Message: `{"reviewer":"missing"}`,
				Success: false,
			},
		},
		{
			&proto.ReviewRequest{
				Reviewer: "",
				Name:     "",
				Review:   "I ate the wrapper as well, and it tasted better than the chocolate",
				Rating:   -5},
			&proto.StatusResponse{
				Message: `{"name":"missing","reviewer":"missing"}`,
				Success: false},
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

func TestNilRequest(t *testing.T) {
	response, err := reviewServiceObject.Review(context.TODO(), nil)

	expectedError := errors.New("ReviewRequest was nil, must be valid reference")

	if response != nil {
		t.Errorf("Expected response to be nil; was '%v'", response)
	}

	if !reflect.DeepEqual(expectedError, err) {
		t.Errorf("Expected error to be '%v'; was '%v'", expectedError, err)
	}
}
