package main

import (
	"fmt"
	"log"
	"os"

	"github.com/micro/cli"
	"github.com/micro/go-micro/cmd"
	"github.com/pivotal-sg/ochoku/reviews"
	proto "github.com/pivotal-sg/ochoku/reviews/proto"
	"golang.org/x/net/context"
)

// create a new command object to use in the service
func newCmd() cmd.Cmd {
	command := cmd.DefaultCmd
	app := command.App()
	app.Flags = append(app.Flags,
		cli.StringFlag{
			Name:        "config",
			EnvVar:      "MICRO_REVIEWS_CONFIG_FILE",
			Usage:       "Path to the file backed configuration for the reviews service",
			Value:       reviews.DefaultConfigFileName,
			Destination: &reviews.ConfigFileName,
		})
	return command
}

func main() {
	client := reviews.NewClient(newCmd())

	allReviews, err := client.AllReviews(context.TODO(), &proto.Empty{})
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}

	fmt.Printf("%v", allReviews)
}
