BINDIR = .

all: reviews-svc reviews-client

reviews-svc:
	go build -o $(BINDIR)/reviews-svc cli/service/reviews.go

reviews-client:
	go build -o $(BINDIR)/reviews-client cli/client/client.go
