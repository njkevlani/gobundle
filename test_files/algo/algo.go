package algo

func Sum(nums []int) int {
	var ans int
	for _, el := range nums {
		ans += el
	}

	return ans
}
