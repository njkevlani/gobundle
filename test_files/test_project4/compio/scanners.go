package compio

import (
	"bufio"
	"os"
)

func NewStdinScanner() *bufio.Scanner {
	s := bufio.NewScanner(bufio.NewReader(os.Stdin))
	s.Split(bufio.ScanWords)

	return s
}
