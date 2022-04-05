package main

import (
	"fmt"
	"math"
)

func main() {
	nums := []int{1, 2, 3, 4, 5}
	PrintSquares(nums)
}
func PrintSquares(nums []int) {
	for _, el := range nums {
		fmt.Printf("Square(%d) = %d\n", el, square(el))
	}
}
func square(n int) int {
	return int(math.Sqrt(float64(n)))
}
