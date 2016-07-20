package main

import (
	"log"

	"github.com/micro/cli"
	"github.com/micro/go-micro/cmd"
	"github.com/pivotal-sg/ochoku/reviews"
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
	service := reviews.NewService(newCmd())
	if err := service.Run(); err != nil {
		log.Fatalf("Error on running Reviews service: '%v'", err)
	}
}
