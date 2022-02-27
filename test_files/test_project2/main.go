package main

import (
	"fmt"

	"github.com/njkevlani/go_bundle/test_files/test_project2/bst"
)

type lorem struct{}

func main() {
	tmp := lorem{}
	fmt.Println(tmp)
	bsTree := &bst.BST{Value: 5}
	bsTree.Add(4)
	bsTree.Add(6)
	bsTree.Add(3)
	bsTree.Add(7)
	bsTree.Add(2)
	bsTree.Add(8)
	bsTree.Add(1)
	bsTree.Add(9)

	fmt.Println("BST had 9 =", bsTree.Has(9))
	fmt.Println("BST had 99 =", bsTree.Has(99))
}
