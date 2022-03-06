.PHONY: build

clean:
	rm -rf ./bin ./build

build-release:
	go build -o bin/go_bundle -ldflags "-s -w" cmd/go_bundle/main.go

build:
	go build -o bin/go_bundle cmd/go_bundle/main.go

test:
	go test $$(go list ./... | grep -v /test_files/)

install: build-release
	mkdir -p ${GOPATH}/bin
	cp ./bin/go_bundle ${GOPATH}/bin/
