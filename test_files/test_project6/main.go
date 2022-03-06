package main

import (
	"fmt"

	"github.com/njkevlani/go_bundle/test_files/test_project6/compio"
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

	permutationRecursive(arr, &count, 0, len(arr)-1)
}

func permutationRecursive(arr []int, count *int, l, r int) {
	if l == r {
		if isAntiFibonacci(arr) {
			for _, el := range arr {
				fmt.Print(el, " ")
			}
			fmt.Println()

			*count++
			if *count == len(arr) {
				return
			}
		}
	} else {
		for i := l; i <= r; i++ {
			arr[l], arr[i] = arr[i], arr[l]

			if l < 2 || arr[l-2]+arr[l-1] != arr[l] {
				permutationRecursive(arr, count, l+1, r)
				if *count == len(arr) {
					return
				}
			}

			arr[l], arr[i] = arr[i], arr[l]
		}
	}
}

func isAntiFibonacci(arr []int) bool {
	for i := 2; i < len(arr); i++ {
		if arr[i-2]+arr[i-1] == arr[i] {
			return false
		}
	}

	return true
}
