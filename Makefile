# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
BINARY_NAME=pg-replication-lag
BINARY_LINUX=$(BINARY_NAME)_linux_amd64
    
all: test build
build:
	$(GOBUILD) -o $(BINARY_NAME) -v github.com/hilli/$(BINARY_NAME)

test: 
	$(GOTEST) -v ./...

clean: 
	$(GOCLEAN)
	rm -f $(BINARY_NAME)
	rm -f $(BINARY_LINUX)

run:
	$(GOBUILD) -o $(BINARY_NAME) -v ./...
	./$(BINARY_NAME)

deps:
	${GOGET} "gopkg.in/yaml.v2"


# Connect to test postgresql servers and expose pg ports on localhost
test-setup:	
	# Master
	ssh -f -o ExitOnForwardFailure=yes -L 5433:localhost:5432 db-master sleep 3600
	# Replica
	ssh -f -o ExitOnForwardFailure=yes -L 5434:localhost:5432 db-slave sleep 3600
    
# Cross compilation
build-linux:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 $(GOBUILD) -o $(BINARY_LINUX) -v github.com/hilli/$(BINARY_NAME)
# docker-build:
# 	docker run --rm -it -v "$(GOPATH)":/go -w /go golang:latest go build -o "$(BINARY_LINUX)" -v "github.com/hilli/pg-replication-lag/$(BINARY_LINUX)"
