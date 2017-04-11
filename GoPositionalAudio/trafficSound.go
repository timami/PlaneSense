package main

import (
	"container/heap"
	"fmt"
  "math"
)

// An IntHeap is a min-heap of ints.
//wip
type AnswerHeap []Answer


func (h AnswerHeap) Len() int           { return len(h) }
func (h AnswerHeap) Less(i, j int) bool { return h[i].distance < h[j].distance }
func (h AnswerHeap) Swap(i, j int)      { h[i].distance , h[j].distance = h[j].distance, h[i].distance }

func (h *AnswerHeap) Push(x interface{}) {
	// Push and Pop use pointer receivers because they modify the slice's length,
	// not just its contents.
	*h = append(*h, x.(Answer))
}

func (h *AnswerHeap) Pop() interface{} {
	old := *h
	n := len(old)
	x := old[n-1]
	*h = old[0 : n-1]
	return x
}

// This example inserts several ints into an IntHeap, checks the minimum,
// and removes them in order of priority.
func recieveTrafficUpdate(h *AnswerHeap) {
	heap.Init(h)
	heap.Push(h, 3)
	fmt.Printf("minimum: %d\n", (*h)[0])
	for h.Len() > 0 {
		fmt.Printf("%d ", heap.Pop(h))
	}
}
