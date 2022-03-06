package bst

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
