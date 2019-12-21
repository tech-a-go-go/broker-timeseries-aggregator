GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
MODULE_NAME=github.com/tech-a-go-go/timeseries-aggregator/cmd/aggregator
BINARY_NAME=bin/aggregator

all: build
build:
	$(GOBUILD) -o $(BINARY_NAME) -v $(MODULE_NAME)
test:
	$(GOTEST) -v ./...
clean:
	$(GOCLEAN) $(MODULE_NAME)
	rm -f $(BINARY_NAME)
run: build
	./$(BINARY_NAME)

