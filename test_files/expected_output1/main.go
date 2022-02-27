package main

import "fmt"

func main() {
	input := []int{1, 2, 3, 4, 5, 2, 6}
	input = append(input, 78)
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
	fmt.Println("square =", algo2_Multiply(n, n))
	return algo2_Multiply(n, n)
}
func algo2_Multiply(x, y int) int {
	return x * y
}
