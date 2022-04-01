package main

import "fmt"

var isSeen map[int]bool

func main() {
	input := []int{1, 2, 3, 4, 5, 2, 6}
	for _, el := range input {
		fmt.Printf("isSeen[%v] = %v\n", el, IsSeen(el))
	}
}
func init() {
	algo_init()
}
func algo_init() {
	isSeen = make(map[int]bool)
}
func IsSeen(num int) bool {
	ret := isSeen[num]
	isSeen[num] = true
	return ret
}
