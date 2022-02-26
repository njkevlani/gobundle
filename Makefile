.PHONY: build

clean:
	rm -r ./bin ./build

build-release:
	go build -o bin/go_bundle -ldflags "-s -w" cmd/go_bundle/main.go

build:
	go build -o bin/go_bundle cmd/go_bundle/main.go

test:
	./scripts/go-build-tests.sh
	./scripts/go-run-tests.sh

install: build-release
	mkdir -p ${GOPATH}/bin
	cp ./bin/go_bundle ${GOPATH}/bin/
