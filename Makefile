# This is how we want to name the binary output
export BINDIR = $(CURDIR)/bin


# These are the values we want to pass for Version and BuildTime
VERSION=1.0.0
BUILD_TIME=`date +%FT%T%z`
NOVENDOR = $(shell glide novendor)

# Setup the -ldflags option for go build here, interpolate the variable values
# LDFLAGS=-ldflags "-X github.com/ariejan/roll/core.Version=${VERSION} -X github.com/ariejan/roll/core.BuildTime=${BUILD_TIME}"
# go build -o ${BINARY} main.go

SUBDIRS = reviews imageservice

.PHONY: all $(SUBDIRS)
all: reviews imageservice
	

.PHONY: imageservice
imageservice:
	make -e -C imageservice

.PHONY: reviews
reviews:
	make -e -C reviews

.PHONY: test
test:
	go test $(NOVENDOR)

.PHONY: clean
clean:
	rm -rf bin || 0;
