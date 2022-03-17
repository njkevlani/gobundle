.PHONY: build

clean:
	rm -rf ./bin ./build

build-release:
	go build -o bin/gobundle -ldflags "-s -w" cmd/gobundle/main.go

build:
	go build -o bin/gobundle cmd/gobundle/main.go

test:
	go test $$(go list ./... | grep -v /test_files/)

install:
	go install -ldflags "-s -w" ./cmd/gobundle
