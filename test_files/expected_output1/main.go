package main

import "fmt"

func main() {
	input := []int{1, 2, 3, 4, 5, 2, 6}
	input = append(input, 78)
	fmt.Println("Sum =", Sum(input))
	s := Sum(input)
	fmt.Println(s)
	empty()
	fmt.Println(Square(2))
}
func Sum(nums []int) int {
	var ans int
	for _, el := range nums {
		ans += el
	}
	return ans
}
func empty() {
}
func Square(n int) int {
	fmt.Println("square =", Multiply(n, n))
	return Multiply(n, n)
}
func Multiply(x, y int) int {
	return x * y
}
