package main

import (
	"log"

	"github.com/micro/go-micro/server"
	"github.com/pivotal-sg/ochoku/reviews"
)

func main() {
	server.Init(
		server.Name("pivotal.io.ochoku.srv"),
		server.Address("localhost:9001"),
	)

	service := reviews.NewService()

	if err := service.Run(); err != nil {
		log.Fatalf("Error on running Reviews service: '%v'", err)
	}
}
