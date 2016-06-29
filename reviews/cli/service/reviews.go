package main

import (
	"log"

	"github.com/micro/go-micro"
	"github.com/pivotal-sg/ochoku/reviews"
)

func run(service micro.Service, ch chan error) {
	ch <- service.Run()
	close(ch)
}

func main() {
	var ch chan error = make(chan error)

	service := reviews.NewService("reviews.db")
	go run(service, ch)

	if err := <-ch; err != nil {
		log.Fatalf("Error on running Reviews service: '%v'", err)
	}
}
