package main

import (
	"fmt"
	"log"
	"os"

	"github.com/pivotal-sg/ochoku/reviews"
	proto "github.com/pivotal-sg/ochoku/reviews/proto"
	"golang.org/x/net/context"
)

func main() {
	client := reviews.NewClient()

	resp, err := client.AllReviews(context.TODO(), &proto.Empty{})
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}

	fmt.Printf("%v", resp)
}
