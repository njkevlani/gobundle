package main

import (
	"fmt"

	"github.com/njkevlani/gobundle/test_files/test_project1/algo"
)

func main() {
	input := []int{1, 2, 3, 4, 5, 2, 6}
	input = append(input, 78)
	fmt.Println("Sum =", algo.Sum(input))

	s := algo.Sum(input)
	fmt.Println(s)

	empty()

	fmt.Println(algo.Square(2))
}

func empty() {

}

func unused() {}
