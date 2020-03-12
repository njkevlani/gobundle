package main

import (
	"fmt"

	"github.com/Workiva/go-datastructures/set"
)

func main() {
	s := set.New()
	input := []int{1, 2, 3, 4, 5, 2, 6}
	for _, i := range input {
		if s.Exists(i) {
			fmt.Println(i)
		} else {
			s.Add(i)
		}
	}
}
