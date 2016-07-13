package reviews_test

import (
	"os"
	"reflect"
	"testing"
	"time"

	"golang.org/x/net/context"

	"github.com/pivotal-sg/ochoku/reviews"
	proto "github.com/pivotal-sg/ochoku/reviews/proto"
)

func TestIntegration(t *testing.T) {
	os.Remove("reviews_test.db")
	go func() {
		service := reviews.NewService("reviews_test.db")
		service.Run()
	}()
	client := reviews.NewClient()

	time.Sleep(5000 * time.Millisecond)
	t.Run("storage=file", func(t *testing.T) {
		reviewRequest := &proto.ReviewRequest{
			Reviewer: "James",
			Name:     "Hershy's Dark Cardboard",
			Review:   "I ate the wrapper as well, and it tasted better than the chocolate",
			Rating:   -5,
		}

		_, err := client.Review(context.TODO(), reviewRequest)

		if err != nil {
			t.Errorf("Expected error to be nil, was  '%v'", err)
		}

	})

	t.Run("storage=file2", func(t *testing.T) {
		expected := &proto.ReviewList{
			Count: 1,
			Reviews: []*proto.ReviewDetails{&proto.ReviewDetails{
				Reviewer: "James",
				Name:     "Hershy's Dark Cardboard",
				Review:   "I ate the wrapper as well, and it tasted better than the chocolate",
				Rating:   -5,
			}}}
		allReviews, err := client.AllReviews(context.TODO(), &proto.Empty{})

		if err != nil {
			t.Errorf("Expected error to be nil, was  '%v'", err)
		}
		if !reflect.DeepEqual(expected, allReviews) {
			t.Errorf("Expected allReviews to be '%v', was  '%v'", expected, allReviews)
		}

	})
}
