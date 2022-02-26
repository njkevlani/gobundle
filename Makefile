.PHONY: build

clean:
	rm -r ./bin ./build

build-release:
	go build -o bin/go_bundle -ldflags "-s -w" cmd/go_bundle/main.go

build:
	go build -o bin/go_bundle cmd/go_bundle/main.go

test: build
	./bin/go_bundle ./test_files/test_project1//main.go
	cmp ./build/main.go ./test_files/expected_output1.go || (echo -e "For checking diff, run\n\tnvim -d ./build/main.go ./test_files/expected_output1.go" && exit 1)
	go run ./build/main.go

install-in-gopath:
	go install -ldflags "-s -w" cmd/go_bundle/main.go
