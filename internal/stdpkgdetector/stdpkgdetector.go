package stdpkgdetector

import (
	"log"

	"golang.org/x/tools/go/packages"
)

var isStdPkg = make(map[string]bool)

func init() {
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

// isStdPkg returns true if given string is a std package in golang, else false.
func IsStdPkg(funcName string) bool {
	return isStdPkg[funcName]
}
