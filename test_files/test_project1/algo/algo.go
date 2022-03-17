package algo

import (
	"fmt"

	"github.com/njkevlani/gobundle/test_files/test_project1/algo2"
)

func Sum(nums []int) int {
	var ans int
	for _, el := range nums {
		ans += el
	}

	return ans
}

func Square(n int) int {
	fmt.Println("square =", algo2.Multiply(n, n))
	return algo2.Multiply(n, n)
}
