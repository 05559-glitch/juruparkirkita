GOCMD=go
GOBUILD=$(GOCMD) build
GOINSTALL=$(GOCMD) install
GORUN=$(GOCMD) run

GOTEST=$(GOCMD) test
BINARY_NAME=my_app

all: build

build:
	$(GOBUILD) -o bin/$(BINARY_NAME) main.go
	@echo "Build complete! Binary created in bin/$(BINARY_NAME)"

run:
	$(GORUN) cmd/api/main.go

test:
	$(GOTEST) ./...

clean:
	rm -rf bin/
	@echo "Cleaned up bin/ directory"

.PHONY: all build run test clean
