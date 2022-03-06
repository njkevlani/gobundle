package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
)

func main() {
	s := NewStdinScanner()
	T := NextInt(s)
	for tt := 1; tt <= T; tt++ {
		n := NextInt(s)
		solve(n)
	}
}
func NewStdinScanner() *bufio.Scanner {
	s := bufio.NewScanner(bufio.NewReader(os.Stdin))
	s.Split(bufio.ScanWords)
	return s
}
func NextInt(s *bufio.Scanner) int {
	return atoi(NextToken(s))
}
func atoi(str string) int {
	i, err := strconv.Atoi(str)
	if err != nil {
		panic(err)
	}
	return i
}
func NextToken(s *bufio.Scanner) string {
	if !s.Scan() {
		fmt.Println(s.Err().Error())
		panic("Could not read input")
	}
	return s.Text()
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
