package algo

import "fmt"

func Sum(nums []int) int {
	var ans int
	for _, el := range nums {
		ans += el
	}

	return ans
}

func Multiply(x, y int) int {
	return x * y
}

func Square(n int) int {
	fmt.Println("square =", Multiply(n, n))
	return Multiply(n, n)
}
