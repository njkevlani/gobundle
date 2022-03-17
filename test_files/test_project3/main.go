package main

import (
	"fmt"

	"github.com/njkevlani/gobundle/test_files/test_project2/bst"
)

type lorem struct{}

func main() {
	tmp := lorem{}
	fmt.Println(tmp)
	bsTree := &bst.BST{Value: 5}
	addInBst(bsTree, []int{4, 6, 3, 7, 2, 8, 1, 9})

	fmt.Println("BST had 9 =", bsTree.Has(9))
	fmt.Println("BST had 99 =", bsTree.Has(99))
	printBst(*bsTree)
}

func addInBst(bsTreeFuncParam *bst.BST, arr []int) {
	for _, el := range arr {
		bsTreeFuncParam.Add(el)
	}
}

func printBst(bsTreeFuncParam bst.BST) {
	fmt.Println(bsTreeFuncParam)
}
