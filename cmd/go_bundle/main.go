package main

import (
	"log"
	"os"
	"path"

	"github.com/njkevlani/go_bundle/internal/go_bundle"
)

func main() {
	if len(os.Args) <= 1 {
		log.Fatalf("Arguments required")
	}

	fileName := path.Clean(os.Args[1])

	goimportedBytes, err := go_bundle.GoBundle(fileName)

	goimportedBytes = append([]byte("// Auto generated using https://github.com/njkevlani/go_bundle\n"), goimportedBytes...)

	err = os.MkdirAll("./build", os.ModePerm)
	if err != nil {
		log.Fatal(err)
	}
	err = os.WriteFile("./build/main.go", goimportedBytes, 0644)
	if err != nil {
		log.Fatal(err)
	}
}
