package main

import (
	"fmt"

	"github.com/njkevlani/gobundle/test_files/test_project7/algo"
)

func main() {
	input := []int{1, 2, 3, 4, 5, 2, 6}

	for _, el := range input {
		fmt.Printf("isSeen[%v] = %v\n", el, algo.IsSeen(el))
	}
}
