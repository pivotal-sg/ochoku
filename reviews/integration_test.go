package reviews_test

import (
	"testing"
	"time"

	"golang.org/x/net/context"

	"github.com/pivotal-sg/ochoku/reviews"
	proto "github.com/pivotal-sg/ochoku/reviews/proto"
)

func TestIntegration(t *testing.T) {
	go func() {
		service := reviews.NewService("reviews_test.db")
		service.Run()
	}()

	time.Sleep(5 * time.Second)
	t.Run("storage=file", func(t *testing.T) {
		client := reviews.NewClient()

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
}
