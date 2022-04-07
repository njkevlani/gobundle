package main

import "fmt"

func main() {
	n1 := node{}
	n2 := node{}
	fmt.Printf("%#v\n", n1)
	fmt.Printf("%#v\n", n2)
}

type node struct {
	val int
	in  []int
}
