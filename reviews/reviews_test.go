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

var reviewServiceObject reviews.ReviewService = reviews.ReviewService{}

func TestReviewChocolateCallsInsert(t *testing.T) {
	storage := MockStorage{store: make([]*proto.ReviewDetails, 0, 0)}

	reviewRequest := proto.ReviewRequest{
		Reviewer: "James",
		Name:     "Hershy's Dark Cardboard",
		Review:   "I ate the wrapper as well, and it tasted better than the chocolate",
		Rating:   -5,
	}
	response, err := reviewServiceObject.Review(context.WithValue(context.Background(), "storage", &storage), &reviewRequest)

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
	storage := &MockStorage{store: make([]*proto.ReviewDetails, 0, 0)}
	storageContext := context.WithValue(context.Background(), "storage", storage)

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
		response, err := reviewServiceObject.Review(storageContext, testValue.Input)

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
	response, err := reviewServiceObject.Review(context.TODO(), nil)

	expectedError := errors.New("ReviewRequest was nil, must be valid reference")

	if response != nil {
		t.Errorf("Expected response to be nil; was '%v'", response)
	}

	if !reflect.DeepEqual(expectedError, err) {
		t.Errorf("Expected error to be '%v'; was '%v'", expectedError, err)
	}
}

func TestGetAllReviews(t *testing.T) {
	storage := &MockStorage{store: make([]*proto.ReviewDetails, 0, 1)}
	storageContext := context.WithValue(context.Background(), "storage", storage)

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

	response, err := reviewServiceObject.AllReviews(storageContext, &proto.Empty{})

	if err != nil {
		t.Errorf("Expected error to be '%v'; was '%v'", nil, err)
	}

	if !reflect.DeepEqual(expectedResponse, response) {
		t.Errorf("Expected response to be '%v'; was '%v'", expectedResponse, response)
	}
}

func TestAllReviewsMissingStorerError(t *testing.T) {
	_, err := reviewServiceObject.AllReviews(
		context.WithValue(context.Background(), "storage", nil), &proto.Empty{})
	if err == nil {
		t.Errorf("Expected error to be 'Storer not set in context'; was '%v'", err)
	}

}
