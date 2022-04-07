package main

import "fmt"

type node struct {
	val int
	in  []int
}

func main() {
    var n1, n2 node

	fmt.Printf("%#v\n", n1)
	fmt.Printf("%#v\n", n2)
}
