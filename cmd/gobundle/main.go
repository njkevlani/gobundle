package main

import (
	"flag"
	"fmt"
	"os"
	"path"

	"github.com/njkevlani/gobundle/internal/gobundle"
)

var version = "unknown"

func usage() {
	fmt.Printf(`version: gobundle-%s
usage: gobundle [-h] file.go

file.go    path to input go file.
-h         show this message.
`, version)
}

func main() {
	flag.Usage = usage
	flag.Parse()

	posArgs := flag.Args()
	if len(posArgs) != 1 {
		fmt.Println("Path to go file is required.")
		usage()
		os.Exit(1)
	}

	fileName := path.Clean(posArgs[0])

	goimportedBytes, err := gobundle.GoBundle(fileName)

	if err != nil {
		panic(err)
	}

	goimportedBytes = append([]byte("// Auto generated using https://github.com/njkevlani/gobundle\n"), goimportedBytes...)

	err = os.MkdirAll("./build", os.ModePerm)
	if err != nil {
		panic(err)
	}
	err = os.WriteFile("./build/main.go", goimportedBytes, 0644)
	if err != nil {
		panic(err)
	}
}
