// Auto generated using https://github.com/njkevlani/go_bundle
package main

import "fmt"

func main() {
	input := []int{1, 2, 3, 4, 5, 2, 6}
	fmt.Println("Sum =", algo_Sum(input))
	s := algo_Sum(input)
	fmt.Println(s)
	empty()
	fmt.Println(algo_Square(2))
}
func algo_Sum(nums []int) int {
	var ans int
	for _, el := range nums {
		ans += el
	}
	return ans
}
func empty() {
}
func algo_Square(n int) int {
	fmt.Println("square =", algo_Multiply(n, n))
	return algo_Multiply(n, n)
}
func algo_Multiply(x, y int) int {
	return x * y
}
