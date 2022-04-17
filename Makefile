.PHONY: build

# If BUILD_VERSION env variable is not defined, take it as dev
BUILD_VERSION ?= dev

LDFLAGS = "-s -w -X main.version=$(BUILD_VERSION)"

GOBUNDLE_PKG = "./cmd/gobundle"

clean:
	rm -rf ./bin ./build

build:
	go build -o bin/gobundle -ldflags $(LDFLAGS) $(GOBUNDLE_PKG)

test:
	go test $$(go list ./... | grep -v /test_files/)

install:
	go install -ldflags $(LDFLAGS) $(GOBUNDLE_PKG)
