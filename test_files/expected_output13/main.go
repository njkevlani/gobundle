package main

import (
	"container/heap"
	"fmt"
)

func main() {
	h := &IntMinHeap{}
	heap.Init(h)
	heap.Push(h, 3)
	heap.Push(h, 1)
	heap.Push(h, 2)
	heap.Push(h, 5)
	for h.Len() > 0 {
		fmt.Println(heap.Pop(h))
	}
}

type IntMinHeap []int

func (h IntMinHeap) Len() int {
	return len(h)
}

func (h IntMinHeap) Less(i, j int) bool {
	return h[i] < h[j]
}

func (h IntMinHeap) Swap(i, j int) {
	h[i], h[j] = h[j], h[i]
}

func (h *IntMinHeap) Push(x interface{}) {
	*h = append(*h, x.(int))
}

func (h *IntMinHeap) Pop() interface{} {
	old := *h
	n := len(old)
	x := old[n-1]
	*h = old[0 : n-1]
	return x
}
