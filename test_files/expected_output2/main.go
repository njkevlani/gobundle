package main

import "fmt"

func main() {
	tmp := lorem{}
	fmt.Println(tmp)
	bsTree := &BST{Value: 5}
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

type lorem struct{}
type BST struct {
	Value int
	Left  *BST
	Right *BST
}

func (bst *BST) Add(val int) {
	if bst.Value >= val {
		if bst.Left == nil {
			bst.Left = &BST{Value: val}
		} else {
			bst.Left.Add(val)
		}
	} else {
		if bst.Right == nil {
			bst.Right = &BST{Value: val}
		} else {
			bst.Right.Add(val)
		}
	}
}
func (bst BST) Has(val int) bool {
	if bst.Value == val {
		return true
	} else if bst.Left != nil && bst.Value >= val {
		return bst.Left.Has(val)
	} else if bst.Right != nil && bst.Value < val {
		return bst.Right.Has(val)
	}
	return false
}
