package main

import "fmt"

type node struct {
	val int
	in  []int
}

func main() {
	nodeVals := []int{1, 2, 3, 4}
	nodeIns := [][]int{nil, {1}, {1}, {1}}

	var g []node
	for i := 0; i < len(nodeVals); i++ {
		g = append(g, node{nodeVals[i], nodeIns[i]})
	}

	for _, n := range g {
		fmt.Printf("%#v\n", n)
	}
}
