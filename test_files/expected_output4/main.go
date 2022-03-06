package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
)

func main() {
	s := NewStdinScanner()
	n := NextInt(s)
	for i := 0; i < n; i++ {
		fmt.Println(i)
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
