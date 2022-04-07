package main

import "fmt"

func main() {
	var n1, n2 Node
	fmt.Printf("%#v\n", n1)
	fmt.Printf("%#v\n", n2)
}

type Node struct {
	val int
	in  []int
}
