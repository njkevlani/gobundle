test:
	go build -o bin/go_bundle cmd/go_bundle/main.go && ./bin/go_bundle ./test_files/main.go && go run ./build/main.go
