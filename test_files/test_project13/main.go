package main

import (
	"container/heap"
	"fmt"

	"github.com/njkevlani/gobundle/test_files/test_project13/algo/heapimpl"
)

func main() {
	h := &heapimpl.IntMinHeap{}
	heap.Init(h)
	heap.Push(h, 3)
	heap.Push(h, 1)
	heap.Push(h, 2)
	heap.Push(h, 5)
	for h.Len() > 0 {
		fmt.Println(heap.Pop(h))
	}
}
