#!/usr/bin/env sh

# Consider removing this once https://github.com/golang/go/issues/15513 is resolved.

for i in $(go list ./...); do go test -v -c $i -o ./bin/$(basename $i).test; done
