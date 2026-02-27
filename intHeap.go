package main

import (
	"container/heap"
	//"fmt"
)

// IntHeap is a min-heap of ints
type IntHeap []int

func (h IntHeap) Len() int {
	return len(h)
}
func (h IntHeap) Less(i, j int) bool {
	return h[i] < h[j]
}
func (h IntHeap) Swap(i, j int) {
	h[i], h[j] = h[j], h[i]
}
func (h *IntHeap) Push(x any) {
	*h = append(*h, x.(int))
}
func (h *IntHeap) Pop() any {
	old := *h
	n := len(old)
	x := old[n-1]
	*h = old[:n-1]
	return x
}
func heapSort(values []int) []int {
	h := &IntHeap{}
	heap.Init(h)

	//push all values into the heap
	for _, v := range values {
		heap.Push(h, v)
	}

	// Pop them out in sorted order
	sorted := make([]int, 0, len(values))
	for h.Len() > 0 {
		sorted = append(sorted, heap.Pop(h).(int))
	}
	return sorted
}
