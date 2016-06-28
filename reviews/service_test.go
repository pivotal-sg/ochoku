package reviews_test

import (
	"reflect"
	"testing"
	"time"

	"github.com/pivotal-sg/ochoku/reviews"
	proto "github.com/pivotal-sg/ochoku/reviews/proto"
	"golang.org/x/net/context"
)

func TestReviewAChocolateWithClient(t *testing.T) {
	service := reviews.NewService()
	client := reviews.NewClient()
	expected := &proto.ReviewDetails{
		Name:     "Choco Name",
		Reviewer: "James",
		Review:   "It was brown in colour",
		Rating:   2,
	}

	var errc chan error = make(chan error)

	go func() {
		errc <- service.Run()
	}()
	rsp, err := client.Review(context.TODO(),
		&proto.ReviewRequest{
			Name:     "Choco Name",
			Reviewer: "James",
			Review:   "It was brown in colour",
			Rating:   2,
		})

	if err != nil {
		t.Errorf("Client Error should have been 'nil', was '%v'", err)
	}

	select {
	case err := <-errc:
		t.Errorf("Client Error should have been 'nil', was '%v'", err)
	case <-time.After(100 * time.Millisecond):
		break
	}

	if !reflect.DeepEqual(expected, rsp) {
		t.Errorf("Expected '%v', got '%v'", expected, rsp)
	}

}
