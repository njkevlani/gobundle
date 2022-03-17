package stdpkgdetector

import (
	"bytes"
	"encoding/gob"
	"log"
	"os"

	"golang.org/x/tools/go/packages"
)

var isStdPkg = make(map[string]bool)
var cacheFile = "/tmp/gobundle_stdpkgdetector.gob"

func init() {
	err := readFromCache()

	if err != nil {
		createNew()
		writeCache()
	}
}

func readFromCache() error {
	b, err := os.ReadFile(cacheFile)
	if err != nil {
		return err
	}

	d := gob.NewDecoder(bytes.NewBuffer(b))

	// Decoding the serialized data
	err = d.Decode(&isStdPkg)
	if err != nil {
		return err
	}

	return nil
}

func createNew() {
	pkgs, err := packages.Load(&packages.Config{Mode: packages.NeedSyntax}, "std")

	if err != nil {
		log.Fatal(err)
	}

	for _, pkg := range pkgs {
		if len(pkg.ID) != 0 {
			isStdPkg[pkg.ID] = true
		}
	}
}

func writeCache() {
	b := new(bytes.Buffer)

	e := gob.NewEncoder(b)

	// Encoding the map
	err := e.Encode(isStdPkg)
	if err != nil {
		log.Println("[WARN] [stdpkgdetector] failed to write cache")
	}

	err = os.WriteFile(cacheFile, b.Bytes(), 0644)

	if err != nil {
		log.Println("[WARN] [stdpkgdetector] failed to write cache")
	}
}

// IsStdPkg returns true if given string is a std package in golang, else false.
func IsStdPkg(funcName string) bool {
	return isStdPkg[funcName]
}
