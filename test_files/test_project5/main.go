package main

import (
	"errors"
	"fmt"

	"github.com/njkevlani/go_bundle/test_files/test_project5/compio"
	"github.com/njkevlani/go_bundle/test_files/test_project5/permutation"
)

func main() {
	s := compio.NewStdinScanner()

	T := compio.NextInt(s)
	for tt := 1; tt <= T; tt++ {
		n := compio.NextInt(s)
		solve(n)
	}
}

func solve(n int) {
	arr := make([]int, n)
	for i := 1; i <= n; i++ {
		arr[i-1] = i
	}

	var count int

	permutation.Permutation(arr, func(p []int) error {
		if isAntiFibonacci(p) {
			for _, el := range p {
				fmt.Print(el, " ")
			}
			fmt.Println()

			count++
			if count == n {
				return errors.New("count completed")
			}
		}
		return nil
	})
}

func isAntiFibonacci(arr []int) bool {
	for i := 2; i < len(arr); i++ {
		if arr[i-2]+arr[i-1] == arr[i] {
			return false
		}
	}

	return true
}
