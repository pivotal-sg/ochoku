_BINDIR = $(BINDIR)

all: reviews-svc reviews-client

reviews-svc:
	go build -o $(_BINDIR)/reviews-svc cli/service/reviews.go

reviews-client:
	go build -o $(_BINDIR)/reviews-client cli/client/client.go
