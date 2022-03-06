package compio

import (
	"bufio"
	"fmt"
	"strconv"
)

func NextToken(s *bufio.Scanner) string {
	if !s.Scan() {
		fmt.Println(s.Err().Error())
		panic("Could not read input")
	}
	return s.Text()
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
