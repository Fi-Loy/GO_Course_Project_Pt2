package main

import (
	"container/heap"
	"fmt"
)

type resRank struct {
	rank int
	res  *Resident
}

// residentHeap is a min-heap of ints
type residentHeap []resRank

func (h residentHeap) Len() int {
	return len(h)
}
func (h residentHeap) Less(i, j int) bool {
	return h[i].rank > h[j].rank
}
func (h residentHeap) Swap(i, j int) {
	h[i], h[j] = h[j], h[i]
}
func (h *residentHeap) Push(x any) {
	*h = append(*h, x.(resRank))
}
func (h *residentHeap) Pop() any {
	old := *h
	n := len(old)
	x := old[n-1]
	*h = old[:n-1]
	return x
}

// WARNING: probably doesnt work correctly
func (h residentHeap) Peek() any {
	return h[0]
}
func heapSort(values []int) []int {
	h := &residentHeap{}
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

func main() {
	rol := []string{"t", "g"}
	res1 := Resident{2, "David", "bungalo", rol, ""}
	res2 := Resident{2, "Jerry", "bungalo", rol, ""}
	res3 := Resident{3, "Todd", "bungalo", rol, ""}
	//res4 := Resident{4, "Bill", "bungalo", rol, ""}
	h := &residentHeap{{1, &res1}, {0, &res2}, {3, &res3}}
	heap.Init(h)
	//heap.Push(h, resRank{2, &res4})
	fmt.Printf("minimum: %d\n", h.Peek())
	for h.Len() > 0 {
		fmt.Printf("%d ", heap.Pop(h))
	}
}
