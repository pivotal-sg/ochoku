package reviews_test

import (
	"errors"
	"reflect"
	"testing"

	"github.com/pivotal-sg/ochoku/reviews"
	proto "github.com/pivotal-sg/ochoku/reviews/proto"
	"golang.org/x/net/context"
)

// MockStorage for testing
type MockStorage struct {
	store        []*proto.ReviewDetails
	insertCalled bool
	getCalled    bool
	listCalled   bool
}

func (s *MockStorage) Get(string) (rd proto.ReviewDetails, err error) {
	s.getCalled = true
	return
}

func (s *MockStorage) List() ([]*proto.ReviewDetails, error) {
	s.listCalled = true
	return s.store, nil
}

func (s *MockStorage) Insert(review proto.ReviewDetails) error {
	s.insertCalled = true
	s.store = append(s.store, &review)
	return nil
}

func (s *MockStorage) reset() {
	s.store = make([]*proto.ReviewDetails, 0, 0)
	s.insertCalled = false
	s.getCalled = false
	s.listCalled = false
}

var storage MockStorage = MockStorage{store: make([]*proto.ReviewDetails, 0, 0)}
var reviewServiceObject reviews.ReviewService = reviews.ReviewService{Store: &storage}

func TestReviewChocolateCallsInsert(t *testing.T) {
	storage.reset()
	ctx := context.TODO()

	var response *proto.StatusResponse = &proto.StatusResponse{}

	reviewRequest := proto.ReviewRequest{
		Reviewer: "James",
		Name:     "Hershy's Dark Cardboard",
		Review:   "I ate the wrapper as well, and it tasted better than the chocolate",
		Rating:   -5,
	}
	err := reviewServiceObject.Review(ctx, &reviewRequest, response)

	if err != nil {
		t.Errorf("expected the response to not have an error, it was %v", err)
	}

	if !response.Success {
		t.Errorf("expected the response to be successful, it wasn't")
	}

	if !storage.insertCalled {
		t.Errorf("expected the insert to have been called; it wasn't")
	}
}

func TestValidInputs(t *testing.T) {
	storage.reset()
	ctx := context.TODO()

	var response *proto.StatusResponse = &proto.StatusResponse{}

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
		storage.reset()
		err := reviewServiceObject.Review(ctx, testValue.Input, response)

		if err != nil {
			t.Errorf("expected the response for '%v' to not have an error, it was %v", testValue.Input, err)
		}

		if !reflect.DeepEqual(testValue.Expected, response) {
			t.Errorf("Expected response to equal '%v', was '%v'", testValue.Expected, response)
		}

		if storage.insertCalled {
			t.Errorf("Expected insert to not have been called it was for '%v'", testValue.Input)
		}
	}
}

func TestNilRequest(t *testing.T) {
	storage.reset()
	var response *proto.StatusResponse
	err := reviewServiceObject.Review(context.TODO(), nil, response)

	expectedError := errors.New("ReviewRequest was nil, must be valid reference")

	if response != nil {
		t.Errorf("Expected response to be nil; was '%v'", response)
	}

	if !reflect.DeepEqual(expectedError, err) {
		t.Errorf("Expected error to be '%v'; was '%v'", expectedError, err)
	}
}

func TestGetAllReviews(t *testing.T) {
	storage.reset()
	ctx := context.TODO()
	var response *proto.ReviewList = &proto.ReviewList{}

	review := proto.ReviewDetails{
		Reviewer: "James",
		Name:     "Hershy's Dark",
		Review:   "I ate the wrapper as well, and it tasted better than the chocolate",
		Rating:   -5}

	storage.Insert(review)

	expectedResponse := &proto.ReviewList{
		Reviews: []*proto.ReviewDetails{&review},
		Count:   1,
	}

	err := reviewServiceObject.AllReviews(ctx, &proto.Empty{}, response)

	if err != nil {
		t.Errorf("Expected error to be '%v'; was '%v'", nil, err)
	}

	if !reflect.DeepEqual(expectedResponse, response) {
		t.Errorf("Expected response to be '%v'; was '%v'", expectedResponse, response)
	}
}

func TestAllReviewsMissingStorerError(t *testing.T) {
	var reviewServiceObject reviews.ReviewService = reviews.ReviewService{}
	var response *proto.ReviewList
	err := reviewServiceObject.AllReviews(context.TODO(), &proto.Empty{}, response)
	if err == nil {
		t.Errorf("Expected error to be 'Storer not set in context'; was '%v'", err)
	}

}
