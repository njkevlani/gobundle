.PHONY: build

# If BUILD_VERSION env variable is not defined, take it as dev
BUILD_VERSION ?= dev

LDFLAGS = "-s -w -X main.version=$(BUILD_VERSION)"

GOBUNDLE_PKG = "./cmd/gobundle"

clean:
	rm -rf ./bin ./build ./packages

build:
	go build -o bin/gobundle -ldflags $(LDFLAGS) $(GOBUNDLE_PKG)

test:
	go test $$(go list ./... | grep -v /test_files/)

install:
	go install -ldflags $(LDFLAGS) $(GOBUNDLE_PKG)

package-linux-amd64:
	@echo Build Linux amd64
	mkdir packages -p
	env GOOS=linux GOARCH=amd64 go build -o bin/gobundle -ldflags $(LDFLAGS) $(GOBUNDLE_PKG)
	tar cf packages/linux_amd64.tar bin/gobundle
	md5sum < packages/linux_amd64.tar | cut -d ' ' -f 1 > packages/linux_amd64.tar.md5.txt

package-linux-arm64:
	@echo Build Linux arm64
	mkdir packages -p
	env GOOS=linux GOARCH=arm64 go build -o bin/gobundle -ldflags $(LDFLAGS) $(GOBUNDLE_PKG)
	tar cf packages/linux_arm64.tar bin/gobundle
	md5sum < packages/linux_arm64.tar | cut -d ' ' -f 1 > packages/linux_arm64.tar.md5.txt

package-darwin-amd64:
	@echo Build Darwin amd64
	mkdir packages -p
	env GOOS=darwin GOARCH=amd64 go build -o bin/gobundle -ldflags $(LDFLAGS) $(GOBUNDLE_PKG)
	tar cf packages/darwin_amd64.tar bin/gobundle
	md5sum < packages/darwin_amd64.tar | cut -d ' ' -f 1 > packages/darwin_amd64.tar.md5.txt

package-darwin-arm64:
	@echo Build Darwin arm64
	mkdir packages -p
	env GOOS=darwin GOARCH=arm64 go build -o bin/gobundle -ldflags $(LDFLAGS) $(GOBUNDLE_PKG)
	tar cf packages/darwin_arm64.tar bin/gobundle
	md5sum < packages/darwin_arm64.tar | cut -d ' ' -f 1 > packages/darwin_arm64.tar.md5.txt

package-windows-amd64:
	@echo Build Windows amd64
	mkdir packages -p
	env GOOS=windows GOARCH=amd64 go build -o bin/gobundle -ldflags $(LDFLAGS) $(GOBUNDLE_PKG)
	tar cf packages/windows_amd64.tar bin/gobundle
	md5sum < packages/windows_amd64.tar | cut -d ' ' -f 1 > packages/windows_amd64.tar.md5.txt

package-all: package-linux-amd64 package-linux-arm64 package-darwin-amd64 package-darwin-arm64 package-windows-amd64
