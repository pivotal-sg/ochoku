package reviews_test

import (
	"reflect"
	"testing"
	"time"

	"os"

	"github.com/micro/go-micro"
	"github.com/micro/go-micro/server"
	"github.com/pivotal-sg/ochoku/reviews"
	proto "github.com/pivotal-sg/ochoku/reviews/proto"
	"golang.org/x/net/context"
)

func runService(service micro.Service, end chan struct{}, ch chan error) {
	select {
	case ch <- service.Run():
		close(ch)
		return
	case <-end:
		ch <- nil
		close(end)
		close(ch)
		return
	}
}

func cleanUp(dbFileName string) {
	os.Remove(dbFileName)
}

var testDbName string = "test.Db"

func TestReviewAChocolateWithClient(t *testing.T) {
	defer cleanUp(testDbName)
	var ch chan error = make(chan error)
	var end chan struct{} = make(chan struct{})

	server.Init(
		server.Name("pivotal.io.ochoku.srv"),
		server.Address("localhost:9001"),
	)

	service := reviews.NewService(testDbName)
	client := reviews.NewClient()

	go runService(service, end, ch)

	expected := &proto.ReviewDetails{
		Name:     "Choco Name",
		Reviewer: "James",
		Review:   "It was brown in colour",
		Rating:   2,
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	rsp, err := client.Review(ctx,
		&proto.ReviewRequest{
			Name:     "Choco Name",
			Reviewer: "James",
			Review:   "It was brown in colour",
			Rating:   2,
		})

	end <- struct{}{}

	if err != nil {
		t.Errorf("Client Error should have been 'nil', was '%v'", err)
	}

	select {
	case err := <-ch:
		t.Errorf("Client Error should have been 'nil', was '%v'", err)
	case <-time.After(100 * time.Millisecond):
		break
	}

	if !reflect.DeepEqual(expected, rsp) {
		t.Errorf("Expected '%v', got '%v'", expected, rsp)
	}

}
