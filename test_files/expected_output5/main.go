package main

import (
	"bufio"
	"errors"
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
	Permutation(arr, func(p []int) error {
		if isAntiFibonacci(p) {
			for _, el := range p {
				fmt.Print(el, " ")
			}
			fmt.Println()
			count++
			if count == n {
				return errors.New("count completed")
			}
		}
		return nil
	})
}
func Permutation(arr []int, processorFunc func([]int) error) {
	_ = permutationRecursive(arr, processorFunc, 0, len(arr)-1)
}
func permutationRecursive(arr []int, processorFunc func([]int) error, l, r int) error {
	if l == r {
		err := processorFunc(arr)
		if err != nil {
			return err
		}
	} else {
		for i := l; i <= r; i++ {
			arr[l], arr[i] = arr[i], arr[l]
			err := permutationRecursive(arr, processorFunc, l+1, r)
			if err != nil {
				return err
			}
			arr[l], arr[i] = arr[i], arr[l]
		}
	}
	return nil
}
func isAntiFibonacci(arr []int) bool {
	for i := 2; i < len(arr); i++ {
		if arr[i-2]+arr[i-1] == arr[i] {
			return false
		}
	}
	return true
}
