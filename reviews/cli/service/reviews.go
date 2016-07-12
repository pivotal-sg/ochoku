package main

import (
	"log"

	"github.com/pivotal-sg/ochoku/reviews"
)

func main() {
	service := reviews.NewService("reviews.db")
	if err := service.Run(); err != nil {
		log.Fatalf("Error on running Reviews service: '%v'", err)
	}
}
