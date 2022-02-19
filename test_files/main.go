package main

import (
	"fmt"

	"github.com/njkevlani/go_bundle/test_files/algo"
)

func main() {
	input := []int{1, 2, 3, 4, 5, 2, 6}
	fmt.Println("Sum =", algo.Sum(input))

	s := algo.Sum(input)
	fmt.Println(s)

	empty()
}

func empty() {

}
