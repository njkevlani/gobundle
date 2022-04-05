package helper

import "fmt"

func PrintSquares(nums []int) {
	for _, el := range nums {
		fmt.Printf("Square(%d) = %d\n", el, square(el))
	}
}
