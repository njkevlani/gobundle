package main

import (
	"fmt"

	"github.com/njkevlani/gobundle/test_files/test_project4/compio"
)

func main() {
	s := compio.NewStdinScanner()

	// Code
	n := compio.NextInt(s)
	for i := 0; i < n; i++ {
		fmt.Println(i)
	}
}
